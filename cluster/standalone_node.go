package cluster

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/gage-technologies/gigo-lib/logging"
	"github.com/sourcegraph/conc"
)

// StandaloneNode
//
//	StandaloneNode mocks the Node interface but executes the roles of
//	leader and follower itself
type StandaloneNode struct {
	ID              int64
	Address         string
	StartTime       time.Time
	Role            NodeRole
	clusterName     string
	ctx             context.Context
	cancel          context.CancelFunc
	leaderRoutine   LeaderRoutine
	followerRoutine FollowerRoutine
	wg              *conc.WaitGroup
	lock            *sync.Mutex
	kv              *sync.Map
	started         bool
	tick            time.Duration
	changeChan      chan StateChangeEvent
	logger          logging.Logger
}

// NewStandaloneNode
//
//	Creates a new StandaloneNode
func NewStandaloneNode(ctx context.Context, id int64, address string, leaderRoutine LeaderRoutine,
	followerRoutine FollowerRoutine, tick time.Duration, logger logging.Logger) *StandaloneNode {
	// create context from system context
	ctx, cancel := context.WithCancel(ctx)

	return &StandaloneNode{
		ID:              id,
		Address:         address,
		StartTime:       time.Now(),
		Role:            NodeRoleLeader,
		clusterName:     "standalone",
		ctx:             ctx,
		cancel:          cancel,
		leaderRoutine:   leaderRoutine,
		followerRoutine: followerRoutine,
		wg:              conc.NewWaitGroup(),
		lock:            &sync.Mutex{},
		kv:              &sync.Map{},
		tick:            tick,
		changeChan:      make(chan StateChangeEvent, 100),
		logger:          logger,
	}
}

// Start
//
//	Starts the StandaloneNode
func (n *StandaloneNode) Start() {
	// acquire lock to check if the node has already started
	n.lock.Lock()
	defer n.lock.Unlock()

	// exit quietly if the node has already been started
	if n.started {
		return
	}

	// mark node as started
	n.started = true

	// launch main loop via wait group
	n.wg.Go(n.loop)
}

// Stop
//
//	Stops the StandaloneNode
func (n *StandaloneNode) Stop() {
	// acquire lock to check if the node has already stopped
	n.lock.Lock()
	defer n.lock.Unlock()

	// exit quietly if the node has already been stopped
	if !n.started {
		return
	}

	// mark node as stopped
	n.started = false

	// cancel the global context
	n.cancel()

	// wait for the main loop to exit
	n.wg.Wait()
}

// Close
//
//	Closes the StandaloneNode
//	This function is a no-op but is needed to satisfy the Node interface
func (n *StandaloneNode) Close() error {
	n.wg.Wait()
	return nil
}

// Put
//
//	Adds a new key-value pair to the cluster that is bound
//	to the node. If the node times out, the key-value pair
//	will be dropped from the cluster. If the key already
//	exists, it will be overwritten.
func (n *StandaloneNode) Put(key string, value string) error {
	key = formatPrefix(n.ID, n.clusterName, fmt.Sprintf("%s/%s", StateDataPrefix, key))
	old, exists := n.kv.Load(key)
	n.kv.Store(key, value)
	event := StateChangeEvent{
		NodeID: n.ID,
		Key:    key,
		Value:  value,
		Type:   EventTypeAdded,
	}
	if exists {
		event.Type = EventTypeModified
		event.OldValue = old.(string)
		event.OldKey = key
	}
	if len(n.changeChan) < 100 {
		n.changeChan <- StateChangeEvent{
			NodeID: n.ID,
			Key:    key,
			Value:  value,
			Type:   EventTypeAdded,
		}
	}
	return nil
}

// Get
//
//	Retrieves the value associated with a given key from
//	the cluster correlated to the node. If the key does
//	not exist, the method will return an empty string.
func (n *StandaloneNode) Get(key string) (string, error) {
	key = formatPrefix(n.ID, n.clusterName, fmt.Sprintf("%s/%s", StateDataPrefix, key))
	value, ok := n.kv.Load(key)
	if ok {
		return value.(string), nil
	}
	return "", nil
}

// GetAsNode
//
//		This is a compliance function for the Node interface.
//	 The function will operate the same as Get when provided
//	 with the self-id of the node instance and will return an
//	 empty string for any other id.
func (n *StandaloneNode) GetAsNode(nodeId int64, key string) (string, error) {
	if nodeId != n.ID {
		return "", nil
	}
	key = formatPrefix(n.ID, n.clusterName, fmt.Sprintf("%s/%s", StateDataPrefix, key))
	value, ok := n.kv.Load(key)
	if ok {
		return value.(string), nil
	}
	return "", nil
}

// GetCluster
//
//	Retrieves the values associated with a given prefix for
//	all nodes in the cluster and returns the values in a
//	map where each key is the node's id and the value is
//	a slice of key value pairs found with passed prefix
//	for the node with the node id stripped from the key.
//
//	Every node in the cluster will return a key in the map.
//	If a node does not have a value for the key the value
//	will be an empty slice for that node's id.
func (n *StandaloneNode) GetCluster(key string) (map[int64][]KV, error) {
	key = formatPrefixCluster(n.clusterName, fmt.Sprintf("%s/%s", StateDataPrefix, key))
	out := make([]KV, 0)
	n.kv.Range(func(k, v any) bool {
		if strings.HasPrefix(k.(string), key) {
			out = append(out, KV{k.(string), v.(string)})
		}
		return true
	})
	return map[int64][]KV{
		n.ID: out,
	}, nil
}

// Delete
//
//	Removes the value associated with a given key from
//	the cluster correlated to the node. If the key does
//	not exist, the method will return nil.
func (n *StandaloneNode) Delete(key string) error {
	key = formatPrefix(n.ID, n.clusterName, fmt.Sprintf("%s/%s", StateDataPrefix, key))
	n.kv.Delete(key)
	return nil
}

// WatchKeyCluster
//
//	Watches to changes for a particular key in any of the
//	nodes in the cluster.
func (n *StandaloneNode) WatchKeyCluster(ctx context.Context, key string) (chan StateChangeEvent, error) {
	key = formatPrefixCluster(n.clusterName, fmt.Sprintf("%s/%s", StateDataPrefix, key))
	out := make(chan StateChangeEvent)
	n.wg.Go(func() {
		for {
			select {
			case <-ctx.Done():
				close(out)
				return
			case event := <-n.changeChan:
				if event.Key != "" && !strings.HasPrefix(event.Key, key) {
					continue
				}
				if event.OldKey != "" && !strings.HasPrefix(event.Key, key) {
					continue
				}
				out <- event
			}
		}
	})
	return out, nil
}

// GetNodes
//
//	Retrieves all the nodes in the cluster.
func (n *StandaloneNode) GetNodes() ([]NodeMetadata, error) {
	return []NodeMetadata{n.GetSelfMetadata()}, nil
}

// GetLeader
//
//	Retrieves the current leader for the cluster. Returns
//	-1 if there is currently no leader for the cluster.
func (n *StandaloneNode) GetLeader() (int64, error) {
	return n.ID, nil
}

// GetSelfMetadata
//
//	Returns the called node's metadata
func (n *StandaloneNode) GetSelfMetadata() NodeMetadata {
	return NodeMetadata{
		ID:      n.ID,
		Address: n.Address,
		Start:   n.StartTime,
		Role:    n.Role,
	}
}

// GetNodeMetadata
//
//	Returns the specified node's metadata and nil if the node is not found
func (n *StandaloneNode) GetNodeMetadata(id int64) (*NodeMetadata, error) {
	if id != n.ID {
		return nil, nil
	}
	meta := n.GetSelfMetadata()
	return &meta, nil
}

// loop
//
//	Loops every 50ms executing the leaderRoutine and followerRoutine
func (n *StandaloneNode) loop() {
	// create ticker to execute every 50ms
	ticker := time.NewTicker(n.tick)
	defer ticker.Stop()

	// loop until context is cancelled
	for {
		select {
		case <-n.ctx.Done():
			return
		case <-ticker.C:
			err := n.leaderRoutine(n.ctx)
			if err != nil {
				n.logger.Errorf("leaderRoutine failed: %v", err)
			}
			err = n.followerRoutine(n.ctx)
			if err != nil {
				n.logger.Errorf("followerRoutine failed: %v", err)
			}
		}
	}
}
