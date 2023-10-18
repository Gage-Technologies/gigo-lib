package cluster

import (
	"context"
	"fmt"
	"math"
	"strconv"
	"sync"
	"time"

	"github.com/coder/retry"
	"github.com/gage-technologies/gigo-lib/logging"
	"github.com/sourcegraph/conc"
	etcd "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"
)

// ClusterNodeOptions
//
//	Options for a cluster node
type ClusterNodeOptions struct {
	Ctx             context.Context
	ID              int64
	Address         string
	Ttl             time.Duration
	ClusterName     string
	EtcdConfig      etcd.Config
	LeaderRoutine   LeaderRoutine
	FollowerRoutine FollowerRoutine
	RoutineTick     time.Duration
	Logger          logging.Logger
}

// ClusterNode
//
//	ClusterNode in a distributed cluster capable of
//	handling leader elections
type ClusterNode struct {
	ID              int64
	Address         string
	Role            NodeRole
	Ttl             time.Duration
	ClusterName     string
	StartTime       time.Time
	started         bool
	lease           etcd.LeaseID
	leaderRoutine   LeaderRoutine
	followerRoutine FollowerRoutine
	wg              *conc.WaitGroup
	campaignWg      *conc.WaitGroup
	lock            *sync.Mutex
	ctx             context.Context
	cancel          context.CancelFunc
	nodeStateCtx    context.Context
	nodeStateCancel context.CancelFunc
	client          *etcd.Client
	clientMu        *sync.RWMutex
	clientConfig    etcd.Config
	session         *concurrency.Session
	election        *concurrency.Election
	tick            time.Duration
	logger          logging.Logger
}

// NewClusterNode
//
//	Creates a new cluster node with a connection to the etcd cluster.
func NewClusterNode(opts ClusterNodeOptions) (*ClusterNode, error) {
	// create context from system context
	ctx, cancel := context.WithCancel(opts.Ctx)

	return &ClusterNode{
		ID:          opts.ID,
		Address:     opts.Address,
		Role:        NodeRoleUnknown,
		Ttl:         opts.Ttl,
		ClusterName: opts.ClusterName,
		// we start nodes with a 0 time and update when
		// they register with the cluster
		StartTime:       time.Unix(0, 0),
		leaderRoutine:   opts.LeaderRoutine,
		followerRoutine: opts.FollowerRoutine,
		wg:              conc.NewWaitGroup(),
		campaignWg:      conc.NewWaitGroup(),
		lock:            &sync.Mutex{},
		ctx:             ctx,
		cancel:          cancel,
		clientMu:        &sync.RWMutex{},
		clientConfig:    opts.EtcdConfig,
		tick:            opts.RoutineTick,
		logger:          opts.Logger,
	}, nil
}

// Start
//
//	Initiate a connection to the etcd cluster and
//	begin campaigning to become leader. Once a role
//	in the cluster has been established, the node
//	begins its work as a leader or follower in the
//	background.
func (n *ClusterNode) Start() {
	// acquire lock to check if the node has already started
	n.lock.Lock()
	defer n.lock.Unlock()

	// exit quietly if the node has already been started
	if n.started {
		return
	}

	// mark node as started
	n.started = true

	// launch state loop via wait group to
	// ensure that we track its closure
	n.wg.Go(n.stateLoop)
}

// Stop
//
//	Stops an active node within the cluster gracefully.
//	The node will cease all participation in the cluster
//	at its earliest convenience. If the node is a leader,
//	it will resign its position as leader before stopping.
//	Stopping removes the node from the cluster but preserves
//	the connection to the etcd cluster such that the node
//	can be started again.
func (n *ClusterNode) Stop() {
	// exit quietly if the node has already been stopped
	if !n.started {
		return
	}

	// cancel global context
	n.cancel()
	// wait for shutdown to complete
	n.wg.Wait()
	n.campaignWg.Wait()

	// lock node to mark node as stopped
	n.lock.Lock()
	defer n.lock.Unlock()

	// mark node as stopped
	n.started = false
	return
}

// Close
//
//	Closes the node and its connection to the etcd cluster.
//	A closed node cannot be restarted.
func (n *ClusterNode) Close() error {
	// ensure that the node is stopped
	n.Stop()
	// close the connection to etcd
	n.clientMu.RLock()
	defer n.clientMu.RUnlock()
	if n.client != nil {
		err := n.client.Close()
		if err != nil {
			return fmt.Errorf("failed to close etcd client: %v", err)
		}
	}
	return nil
}

// Put
//
//	Adds a new key-value pair to the cluster that is bound
//	to the node. If the node times out, the key-value pair
//	will be dropped from the cluster. If the key already
//	exists, it will be overwritten.
func (n *ClusterNode) Put(key string, value string) error {
	// retrieve node's etcd lease
	n.lock.Lock()
	lease := n.lease
	n.lock.Unlock()

	// fail if we don't have a valid lease yet
	if lease == 0 {
		return ErrNoLease
	}

	// get the etcd client
	client, err := n.getClient()
	if err != nil {
		return fmt.Errorf("failed to get etcd client: %v", err)
	}

	// write kv pair to etcd under node lease
	_, err = client.Put(
		context.TODO(),
		formatPrefix(n.ID, n.ClusterName, fmt.Sprintf("%s/%s", StateDataPrefix, key)),
		value,
		etcd.WithLease(lease),
	)
	if err != nil {
		return fmt.Errorf("failed to save key to etcd under node: %v", err)
	}
	return nil
}

// Get
//
//	Retrieves the value associated with a given key from
//	the cluster correlated to the node. If the key does
//	not exist, the method will return an empty string.
func (n *ClusterNode) Get(key string) (string, error) {
	// get the etcd client
	client, err := n.getClient()
	if err != nil {
		return "", fmt.Errorf("failed to get etcd client: %v", err)
	}

	// read kv pair from etcd under node lease
	pair, err := client.Get(
		context.TODO(),
		formatPrefix(n.ID, n.ClusterName, fmt.Sprintf("%s/%s", StateDataPrefix, key)),
		etcd.WithLease(n.lease),
	)
	if err != nil {
		return "", fmt.Errorf("failed to read key from etcd under node: %v", err)
	}
	if pair == nil || len(pair.Kvs) == 0 {
		return "", nil
	}
	return string(pair.Kvs[0].Value), nil
}

// GetAsNode
//
//	Retrieves the value associated with a given key from
//	the cluster correlated to the passed node id. If the
//	key does not exist, the method will return an empty
//	string.
func (n *ClusterNode) GetAsNode(nodeId int64, key string) (string, error) {
	// get the etcd client
	client, err := n.getClient()
	if err != nil {
		return "", fmt.Errorf("failed to get etcd client: %v", err)
	}

	// read kv pair from etcd under node lease
	pair, err := client.Get(
		context.TODO(),
		formatPrefix(nodeId, n.ClusterName, fmt.Sprintf("%s/%s", StateDataPrefix, key)),
		etcd.WithLease(n.lease),
	)
	if err != nil {
		return "", fmt.Errorf("failed to read key from etcd under node %d: %v", nodeId, err)
	}
	if pair == nil || len(pair.Kvs) == 0 {
		return "", nil
	}
	return string(pair.Kvs[0].Value), nil
}

// GetCluster
//
//	Retrieves the values associated with a given prefix for
//	all nodes in the cluster and returns the values in a
//	map where each key is the node's id and the value is
//	the value associated with that node.
//
//	Every node in the cluster will return a key in the map.
//	If a node does not have a value for the key the value
//	will be an empty string for that node's id.
func (n *ClusterNode) GetCluster(key string) (map[int64][]KV, error) {
	// retrieve all nodes in the cluster
	nodes, err := n.GetNodes()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve nodes in cluster: %v", err)
	}

	// create output map
	out := make(map[int64][]KV)
	for _, node := range nodes {
		out[node.ID] = make([]KV, 0)
	}

	// get the etcd client
	client, err := n.getClient()
	if err != nil {
		return nil, fmt.Errorf("failed to get etcd client: %v", err)
	}

	// read kv pair from etcd under node lease
	pairs, err := client.Get(
		context.TODO(),
		formatPrefixCluster(n.ClusterName, fmt.Sprintf("%s/%s", StateDataPrefix, key)),
		etcd.WithLease(n.lease),
		etcd.WithPrefix(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to read key from etcd under node: %v", err)
	}

	// handle nil return
	if pairs == nil {
		return nil, nil
	}

	// iterate over kv pairs updating the value
	// for each node in the output map
	for _, kv := range pairs.Kvs {
		// attempt to extract node if from key
		id, base, err := extractNodeIDFromKey(string(kv.Key))
		if err != nil {
			return nil, fmt.Errorf("failed to extract id from key: %s", kv.Key)
		}

		// save value to output map under node id
		out[id] = append(out[id], KV{
			Key:   base,
			Value: string(kv.Value),
		})
	}

	return out, nil
}

// Delete
//
//	Removes the value associated with a given key from
//	the cluster correlated to the node. If the key does
//	not exist, the method will return nil.
func (n *ClusterNode) Delete(key string) error {
	// get the etcd client
	client, err := n.getClient()
	if err != nil {
		return fmt.Errorf("failed to get etcd client: %v", err)
	}

	// delete kv pair from etcd under node
	_, err = client.Delete(
		context.TODO(),
		formatPrefix(n.ID, n.ClusterName, fmt.Sprintf("%s/%s", StateDataPrefix, key)),
	)
	if err != nil {
		return fmt.Errorf("failed to delete key from etcd under node: %v", err)
	}
	return nil
}

// WatchKeyCluster
//
//	Watches to changes for a particular key in any of the
//	nodes in the cluster.
func (n *ClusterNode) WatchKeyCluster(ctx context.Context, key string) (chan StateChangeEvent, error) {
	// get the etcd client
	client, err := n.getClient()
	if err != nil {
		return nil, fmt.Errorf("failed to get etcd client: %v", err)
	}

	// watch for changes in etcd under node
	watcher := client.Watch(
		ctx,
		formatPrefixCluster(n.ClusterName, fmt.Sprintf("%s/%s", StateDataPrefix, key)),
		etcd.WithPrefix(),
		etcd.WithPrevKV(),
	)

	// create channel to pipe state changes back to the caller
	ch := make(chan StateChangeEvent)

	// launch the watcher routine in the nodes wait group
	n.wg.Go(func() {
		// close channel when we're done
		defer close(ch)

		// loop indefinitely
		for {
			// wait for a state change or the context to cancel
			select {
			// exit if our context from the caller is done
			case <-ctx.Done():
				return
			// exit if our node context is done
			case <-n.ctx.Done():
				return
			// handle state changes
			case events, ok := <-watcher:
				// exit on channel close
				if !ok {
					return
				}

				// iterate events from the watcher
				for _, e := range events.Events {
					// determine event type
					// default to deleted because that is the simplest event type
					// then handle PUT events by checking if the state previously
					// existed to determine if this was an addition or modification
					eventType := EventTypeDeleted
					if e.Type == etcd.EventTypePut {
						// set event type to added if there is no previous state
						// otherwise set event type to modified
						if e.PrevKv == nil {
							eventType = EventTypeAdded
						} else {
							eventType = EventTypeModified
						}
					}

					// set key with the most recent key that we have
					k := ""
					if e.PrevKv != nil {
						k = string(e.PrevKv.Key)
					}
					if e.Kv != nil {
						k = string(e.Kv.Key)
					}

					// get node if from key
					nodeId, _, err := extractNodeIDFromKey(k)
					if err != nil {
						n.logger.Errorf("failed to extract node id from key %q: %v", k, err)
						continue
					}

					// create state change event
					event := StateChangeEvent{
						Type:   eventType,
						NodeID: nodeId,
					}

					// add remaining fields as they are available
					if eventType != EventTypeAdded && e.PrevKv != nil {
						event.OldKey = string(e.PrevKv.Key)
						event.OldValue = string(e.PrevKv.Value)
					}
					if eventType != EventTypeDeleted && e.Kv != nil {
						event.Key = string(e.Kv.Key)
						event.Value = string(e.Kv.Value)
					}

					// send state change to channel
					ch <- event
				}
			}
		}
	})

	return ch, nil
}

// GetNodes
//
//	Retrieves all the nodes in the cluster.
func (n *ClusterNode) GetNodes() ([]NodeMetadata, error) {
	// get the etcd client
	client, err := n.getClient()
	if err != nil {
		return nil, fmt.Errorf("failed to get etcd client: %v", err)
	}

	// retrieve all nodes in the cluster
	pairs, err := client.Get(
		context.TODO(),
		formatPrefixCluster(n.ClusterName, NodesPrefix),
		etcd.WithPrefix(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve nodes in cluster: %v", err)
	}

	// create output slice
	nodes := make([]NodeMetadata, 0)

	// iterate kv pairs
	for _, kv := range pairs.Kvs {
		// unmarshal node metadata
		meta, err := UnmarshalNodeMetadata(kv.Value)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal node metadata: %v", err)
		}

		// save value to output slice
		nodes = append(nodes, meta)
	}

	return nodes, nil
}

// GetLeader
//
//	Retrieves the current leader for the cluster. Returns
//	-1 if there is currently no leader for the cluster.
func (n *ClusterNode) GetLeader() (int64, error) {
	// return error if we are not a part of the cluster
	if n.election == nil {
		return -1, fmt.Errorf("node is not part of a cluster - call Start to join cluster")
	}

	// check for existing leader
	res, err := n.election.Leader(context.TODO())
	if err != nil {
		// return -1 with no error if we don't have a leader yet
		if err == concurrency.ErrElectionNoLeader {
			return -1, nil
		}
		return -1, fmt.Errorf("failed to check cluster leader: %v", err)
	}

	// ensure that we have a value to load
	if len(res.Kvs) == 0 {
		return -1, fmt.Errorf("empty response from leader election")
	}

	// attempt to load leader id and format to an int
	id, err := strconv.ParseInt(string(res.Kvs[0].Value), 10, 64)
	if err != nil {
		return -1, fmt.Errorf("failed to parse leader id: %v", err)
	}

	return id, nil
}

// GetSelfMetadata
//
//	Returns the called node's metadata
func (n *ClusterNode) GetSelfMetadata() NodeMetadata {
	return NodeMetadata{
		ID:      n.ID,
		Address: n.Address,
		Start:   n.StartTime,
		Role:    n.Role,
	}
}

// GetNodeMetadata
//
//	Returns the specified node's metadata
func (n *ClusterNode) GetNodeMetadata(id int64) (*NodeMetadata, error) {
	// get the etcd client
	client, err := n.getClient()
	if err != nil {
		return nil, fmt.Errorf("failed to get etcd client: %v", err)
	}

	// save latest node metadata to etcd storage
	metaRes, err := client.Get(
		context.TODO(),
		formatPrefix(id, n.ClusterName, NodesPrefix),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve node metadata: %v", err)
	}

	// return nil if the meta is not found
	if metaRes == nil || len(metaRes.Kvs) == 0 {
		return nil, nil
	}

	// unmarshall node metadata
	meta, err := UnmarshalNodeMetadata(metaRes.Kvs[0].Value)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal node metadata: %v", err)
	}

	return &meta, nil
}

// getClient
//
//	Returns the etcd client
func (n *ClusterNode) getClient() (*etcd.Client, error) {
	// acquire read on client mutex
	n.clientMu.RLock()

	// return client if it is not nil and the connection is healthy
	if n.client != nil {
		// check if the client is healthy
		ctx, cancel := context.WithTimeout(n.ctx, 3*time.Second)
		defer cancel()
		if _, err := n.client.Get(ctx, "/", etcd.WithLimit(1), etcd.WithKeysOnly()); err == nil {
			n.clientMu.RUnlock()
			return n.client, nil
		}
		n.logger.Debugf("(cluster: %d) etcd client is unhealthy - creating new client", n.ID)
	}

	// upgrade to write lock
	n.clientMu.RUnlock()
	n.clientMu.Lock()
	defer n.clientMu.Unlock()

	// create a new client
	client, err := etcd.New(n.clientConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create etcd client: %v", err)
	}

	// save client
	n.client = client

	return client, nil
}

// createClient
//
//	Creates a new etcd client
func (n *ClusterNode) createClient() (*etcd.Client, error) {
	// acquire write lock on client mutex
	n.clientMu.Lock()
	defer n.clientMu.Unlock()

	// create a new client
	client, err := etcd.New(n.clientConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create etcd client: %v", err)
	}

	// save client
	n.client = client

	return client, nil
}

// alive
//
//	Check if the node is alive
func (n *ClusterNode) alive() bool {
	select {
	case <-n.ctx.Done():
		return false
	default:
		return true
	}
}

// newClusterSession
//
//	Creates a new cluster session and leader election
func (n *ClusterNode) newClusterSession(sessionContext context.Context) (*concurrency.Session, *concurrency.Election, error) {
	// ensure that the ttl in at least 1 second
	ttl := int(math.Floor(n.Ttl.Seconds()))
	if ttl < 1 {
		ttl = 1
	}

	// get the etcd client
	client, err := n.getClient()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get etcd client: %v", err)
	}

	// create new session for leader election
	session, err := concurrency.NewSession(
		client,
		concurrency.WithTTL(ttl),
		concurrency.WithContext(sessionContext),
		concurrency.WithLease(n.lease),
	)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create session: %v", err)
	}

	// create new election
	election := concurrency.NewElection(session, fmt.Sprintf("/%s/%s", n.ClusterName, ElectionPrefix))

	return session, election, nil
}

// campaign
//
//	Campaigns the node to the cluster leader continuously
func (n *ClusterNode) campaign(ctx context.Context, election *concurrency.Election) {
	// create ticker to tick once per second
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		// exit if the context is now
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			// create timeout for leader election
			campaignCtx, campaignCancel := context.WithTimeout(ctx, time.Second)
			// campaign for leadership
			_ = election.Campaign(campaignCtx, fmt.Sprintf("%d", n.ID))
			campaignCancel()
		default:
		}
	}
}

// cleanupNode
//
//	Performs a cleanup for the node
func (n *ClusterNode) cleanupNode() {
	n.lock.Lock()
	defer n.lock.Unlock()

	// get the etcd client
	client, err := n.getClient()
	if err != nil {
		n.logger.Errorf("(cluster: %d) failed to get etcd client on cleanup: %v", n.ID, err)
		return
	}

	// close node state ctx
	if n.nodeStateCancel != nil {
		n.nodeStateCancel()
	}

	// wait for campaigner to close
	n.campaignWg.Wait()

	// we use only new contexts here since we want the cleanup
	// to occur regardless of the fact that the system may have
	// already cancelled

	// resign if we are the leader
	if n.Role == NodeRoleLeader {
		resignCtx, resignCancel := context.WithTimeout(context.TODO(), time.Second)
		err := n.election.Resign(resignCtx)
		if err != nil {
			n.logger.Errorf("(cluster: %d) failed to resign as leader on cleanup: %v", n.ID, err)
		}
		resignCancel()
	}

	// revoke lease
	revokeCtx, revokeCancel := context.WithTimeout(context.TODO(), time.Second)
	_, err = client.Revoke(revokeCtx, n.lease)
	if err != nil {
		n.logger.Errorf("(cluster: %d) failed to revoke lease on cleanup: %v", n.ID, err)
	}
	revokeCancel()

	// close session and election
	if n.session != nil {
		err = n.session.Close()
		if err != nil {
			n.logger.Warnf("(cluster: %d) failed to close session on cleanup: %v", n.ID, err)
		}
	}

	// clean node state
	n.Role = NodeRoleUnknown
	n.lease = 0
	n.session = nil
	n.election = nil
	n.nodeStateCtx = nil
	n.nodeStateCancel = nil

	// remove node from cluster state
	_, err = client.Delete(
		context.TODO(),
		formatPrefix(n.ID, n.ClusterName, NodesPrefix),
		etcd.WithLease(n.lease),
	)
	if err != nil {
		n.logger.Errorf("(cluster: %d) failed to remove node from cluster state: %v", n.ID, err)
	}
}

// updateNodeRole
//
//	Updates the node role
//	Handles context management, campaigning, and updating
//	the cluster state.
func (n *ClusterNode) updateNodeRole(leader int64, sessionCtx context.Context) error {
	// lock node for state update
	n.lock.Lock()
	defer n.lock.Unlock()

	// get the etcd client
	client, err := n.getClient()
	if err != nil {
		return fmt.Errorf("failed to get etcd client: %v", err)
	}

	// exit if there is no change to the node's role
	if (leader == n.ID && n.Role == NodeRoleLeader) ||
		(leader != n.ID && n.Role == NodeRoleFollower) {
		return nil
	}

	// define new role for the node
	newRole := NodeRoleFollower
	if leader == n.ID {
		newRole = NodeRoleLeader
	}

	n.logger.Infof("(cluster: %d) updating cluster role %v -> %v", n.ID, n.Role, newRole)

	// cancel current node state context and create a new one
	n.nodeStateCancel()
	n.nodeStateCtx, n.nodeStateCancel = context.WithCancel(sessionCtx)
	n.campaignWg.Go(func() {
		n.campaign(n.nodeStateCtx, n.election)
	})

	// validate our local cluster tick
	err = n.validateClusterTick()
	if err != nil {
		// don't wrap this error since the returned hour may carry
		// a formatted error
		return err
	}

	// mark the node's start time if this is the first time
	// we are registering this node
	if n.Role == NodeRoleUnknown {
		n.StartTime = time.Now()
	}

	// update node role
	n.Role = newRole

	// get clusters metadata
	meta := n.GetSelfMetadata()

	// marshall node metadata
	metaBytes, err := meta.Marshal()
	if err != nil {
		return fmt.Errorf("failed to marshal node metadata: %v", err)
	}

	// save latest node metadata to etcd storage
	_, err = client.Put(
		context.TODO(),
		formatPrefix(n.ID, n.ClusterName, NodesPrefix),
		string(metaBytes),
		etcd.WithLease(n.lease),
	)
	if err != nil {
		return fmt.Errorf("failed to save node to etcd: %v", err)
	}

	n.logger.Infof("(cluster: %d) cluster node is now %v", n.ID, n.Role)

	return nil
}

// validateClusterTick
//
//	Validates our cluster tick to ensure that
//	we are in agreement with the rest of the cluster
func (n *ClusterNode) validateClusterTick() error {
	// get the etcd client
	client, err := n.getClient()
	if err != nil {
		return fmt.Errorf("failed to get etcd client: %v", err)
	}

	// put our tick in the cluster storage
	_, err = client.Put(
		context.TODO(),
		formatPrefix(n.ID, n.ClusterName, fmt.Sprintf("%s/%s", StateDataPrefix, "tick")),
		fmt.Sprintf("%d", n.tick.Nanoseconds()),
		etcd.WithLease(n.lease),
	)
	if err != nil {
		return fmt.Errorf("failed to save tick to etcd: %v", err)
	}

	// get tick for all nodes
	clusterTicks, err := n.GetCluster("tick")
	if err != nil {
		return fmt.Errorf("failed to get tick from etcd: %v", err)
	}

	// iterate ticks checking them against out tick
	for _, tick := range clusterTicks {
		if len(tick) == 0 || tick[0].Value != fmt.Sprintf("%d", n.tick.Nanoseconds()) {
			return ErrClusterTickDisagreement
		}
	}

	// return now that we have confirmed there is
	// agreement on out tick rate
	return nil
}

// stateLoop
//
//	Manages node state in cluster
func (n *ClusterNode) stateLoop() {
	// retry forever
	for outerRetrier := retry.New(time.Millisecond*10, time.Second*3); outerRetrier.Wait(n.ctx); {
		n.logger.Debugf("(cluster: %d) starting state loop", n.ID)

		// exit if we are done
		if !n.alive() {
			return
		}

		n.logger.Debugf("(cluster: %d) retrieving client for state loop", n.ID)
		// get the etcd client
		client, err := n.getClient()
		if err != nil {
			n.logger.Errorf("(cluster: %d) failed to retrieve client: %v", n.ID, err)
			continue
		}

		n.logger.Debugf("(cluster: %d) retrieving lease for state loop", n.ID)
		// retrieve lease from etcd for up to Ttl so that we can campaign
		ttl := int64(math.Floor(n.Ttl.Seconds()))
		if ttl < 1 {
			ttl = 1
		}
		lease, err := client.Grant(n.ctx, ttl)
		if err != nil {
			n.logger.Errorf("(cluster: %d) failed to retrieve lease: %v", n.ID, err)
			continue
		}

		n.logger.Debugf("(cluster: %d) acquired lease %s", n.ID, lease.String())

		// update lease with node lock
		n.lock.Lock()
		n.lease = lease.ID
		n.lock.Unlock()

		n.logger.Debugf("(cluster: %d) creating session for state loop", n.ID)
		// create a new session for leader election
		session, election, err := n.newClusterSession(n.ctx)
		if err != nil {
			n.logger.Errorf("(cluster: %d) failed to create session: %v", n.ID, err)
			n.cleanupNode()
			continue
		}

		// create context for this session
		sessionCtx, sessionCancel := context.WithCancel(n.ctx)

		// update node with session and election
		n.session = session
		n.election = election

		n.logger.Debugf("(cluster: %d) checking cluster for leader in state loop", n.ID)
		// check for existing leader
		res, err := n.election.Leader(sessionCtx)
		if err != nil && err != concurrency.ErrElectionNoLeader {
			n.logger.Errorf("(cluster: %d) failed to check for existing leader: %v", n.ID, err)
			n.cleanupNode()
			sessionCancel()
			continue
		}

		// resign if we are the leader
		// the only case that this occurs is the same node restarts
		// at which point we should re-elect because it is too dangerous
		// to assume the leader role given that there could be a race condition
		if err == nil && string(res.Kvs[0].Value) == fmt.Sprintf("%d", n.ID) {
			// we use a new context here since we want the resignation
			// to occur regardless of if the system cancels
			resignCtx, resignCancel := context.WithTimeout(context.TODO(), time.Second)
			_ = n.election.Resign(resignCtx)
			resignCancel()
		}

		// begin state loop
		for innerRetrier := retry.New(time.Millisecond*10, time.Second); innerRetrier.Wait(sessionCtx); {
			// exit if we are done
			if !n.alive() {
				n.cleanupNode()
				sessionCancel()
				return
			}

			n.logger.Debugf("(cluster: %d) starting campaign for state loop", n.ID)
			// campaign for leadership using node state context
			// NOTE: this will kill the campaigner every time that the role
			// of the node changes which is not particularly desirable but it
			// is better than using the session context because we can terminate
			// the node context (and therefore the campaigner) before a leader
			// voluntarily resigns whereas the session context does not permit
			// us to cancel the context before a voluntary resignation
			n.nodeStateCtx, n.nodeStateCancel = context.WithCancel(sessionCtx)
			n.campaignWg.Go(func() {
				n.campaign(n.nodeStateCtx, n.election)
			})

			n.logger.Debugf("(cluster: %d) beginning campaign observation loop for state loop", n.ID)
			// create observation channel to watch changes
			// to the cluster leader
			observationChan := n.election.Observe(sessionCtx)

			// create default ticker for 50ms to execute
			tick := 50 * time.Millisecond
			ticker := time.NewTicker(tick)

			// create lease ticker to execute once every 500ms
			leaseTicker := time.NewTicker(500 * time.Millisecond)

			// create boolean to track if we can continue looping
			continueLoop := true

			// create index to track loop count
			loopCount := 0

			// track the amount of concurrent executions that we are unknown
			unknownCount := 0

			// loop until we have a role change executing
			// the leader or follower routine every 50ms
			for continueLoop {
				// increment loop count here because we always
				// execute at the start but may use a continue
				// statement and skip the end of the loop
				loopCount++

				select {
				case <-sessionCtx.Done():
					n.logger.Debugf("(cluster: %d) cluster node is exiting - beginning cleanup", n.ID)
					// perform node cleanup
					n.cleanupNode()

					// create a new client for the node
					client, err = n.createClient()
					if err != nil {
						n.logger.Errorf("(cluster: %d) failed to create new client during cleanup: %v", n.ID, err)
					}

					n.logger.Debugf("(cluster: %d) cluster node is exiting - cleanup complete", n.ID)
					// mark continue loop as false to exit the loop
					continueLoop = false
					continue
				case <-leaseTicker.C:
					// renew the node's lease
					_, err := client.KeepAliveOnce(sessionCtx, n.lease)
					if err != nil {
						n.logger.Errorf("(cluster: %d) failed to renew lease: %v", n.ID, err)
						// we re-elect here since it is the simplest way to resolve the
						// cluster state
						sessionCancel()
						continue
					}
					n.logger.Debugf("(cluster: %d) cluster node renewed lease", n.ID)
				case <-ticker.C:
					// TODO: replace with proper tracing
					start := time.Now()
					if n.Role == NodeRoleLeader {
						n.logger.Debugf("(cluster: %d) executing leader routine", n.ID)
						err := n.leaderRoutine(n.nodeStateCtx)
						if err != nil {
							n.logger.Errorf("(cluster: %d) leader routine failed: %s", n.ID, err)
						}
						n.logger.Debugf("(cluster: %d) leader routine took %v", n.ID, time.Since(start))
						unknownCount = 0
					} else if n.Role == NodeRoleFollower {
						n.logger.Debugf("(cluster: %d) executing follower routine", n.ID)
						err := n.followerRoutine(n.nodeStateCtx)
						if err != nil {
							n.logger.Errorf("(cluster: %d) follower routine failed: %s", n.ID, err)
						}
						n.logger.Debugf("(cluster: %d) follower routine took %v", n.ID, time.Since(start))
						unknownCount = 0
					} else {
						// abort if we are unknown for more than 100 executions
						if unknownCount > 100 {
							n.logger.Errorf("(cluster: %d) node role is unknown for too long - aborting", n.ID)
							sessionCancel()
							continue
						}

						// only log once per second
						if loopCount%20 == 0 {
							n.logger.Debugf("(cluster: %d) node role is currently unknown", n.ID)
						}
						unknownCount++
					}
				case clusterUpdate := <-observationChan:
					// handle lost leadership
					if len(clusterUpdate.Kvs) == 0 {
						n.logger.Errorf("(cluster: %d) cluster without leader for too long - aborting", n.ID)
						sessionCancel()
						continue
					}

					// parse node id from key
					leaderId, err := strconv.ParseInt(string(clusterUpdate.Kvs[0].Value), 10, 64)
					if err != nil {
						n.logger.Errorf("(cluster: %d) failed to parse leader id: %v", n.ID, err)
						// we have to bail here since we lost track of the leader
						sessionCancel()
						continue
					}

					// update node role
					err = n.updateNodeRole(leaderId, sessionCtx)
					if err != nil {
						n.logger.Errorf("(cluster: %d) failed to update node role: %v", n.ID, err)
						// we have to bail here since we lost track of the leader
						sessionCancel()
						continue
					}

					if tick != n.tick {
						n.logger.Infof("(cluster: %d) updating cluster tick from %v -> %v", n.ID, tick, n.tick)
						// update the ticker
						tick = n.tick
						ticker.Reset(tick)
					}
				}
			}

			// break loop since the only case that we are here is if it's time
			// to exit or re-establish a new connection
			ticker.Stop()
			break
		}

		// cancel election context
		sessionCancel()
	}
}
