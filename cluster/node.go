package cluster

import "context"

// Node
//
//	Interface to interact with a cluster node
type Node interface {
	// Start
	//
	//  Initiate a connection to the etcd cluster and
	//  begin campaigning to become leader. Once a role
	//  in the cluster has been established, the node
	//  begins its work as a leader or follower in the
	//  background.
	Start()

	// Stop
	//
	//  Stops an active node within the cluster gracefully.
	//  The node will cease all participation in the cluster
	//  at its earliest convenience. If the node is a leader,
	//  it will resign its position as leader before stopping.
	//  Stopping removes the node from the cluster but preserves
	//  the connection to the etcd cluster such that the node
	//  can be started again.
	Stop()

	// Close
	//
	//  Closes the node and its connection to the etcd cluster.
	//  A closed node cannot be restarted.
	Close() error

	// Put
	//
	//  Adds a new key-value pair to the cluster that is bound
	//  to the node. If the node times out, the key-value pair
	//  will be dropped from the cluster. If the key already
	//  exists, it will be overwritten.
	Put(key string, value string) error

	// Get
	//
	//  Retrieves the value associated with a given key from
	//  the cluster correlated to the node. If the key does
	//  not exist, the method will return an empty string.
	Get(key string) (string, error)

	// GetAsNode
	//
	//  Retrieves the value associated with a given key from
	//  the cluster correlated to the passed node id. If the
	//  key does not exist, the method will return an empty
	//  string.
	GetAsNode(nodeId int64, key string) (string, error)

	// GetCluster
	//
	//  Retrieves the values associated with a given prefix for
	//  all nodes in the cluster and returns the values in a
	//  map where each key is the node's id and the value is
	//  a slice of key value pairs found with passed prefix
	//  for the node with the node id stripped from the key.
	//
	//  Every node in the cluster will return a key in the map.
	//  If a node does not have a value for the key the value
	//  will be an empty slice for that node's id.
	GetCluster(key string) (map[int64][]KV, error)

	// Delete
	//
	//  Removes the value associated with a given key from
	//  the cluster correlated to the node. If the key does
	//  not exist, the method will return nil.
	Delete(key string) error

	// WatchKeyCluster
	//
	//  Watches to changes for a particular key in any of the
	//  nodes in the cluster.
	WatchKeyCluster(ctx context.Context, key string) (chan StateChangeEvent, error)

	// GetNodes
	//
	//  Retrieves all the nodes in the cluster.
	GetNodes() ([]NodeMetadata, error)

	// GetLeader
	//
	//  Retrieves the current leader for the cluster. Returns
	//  -1 if there is currently no leader for the cluster.
	GetLeader() (int64, error)

	// GetSelfMetadata
	//
	//  Returns the called node's metadata
	GetSelfMetadata() NodeMetadata

	// GetNodeMetadata
	//
	//  Returns the specified node's metadata and nil if the node is not found
	GetNodeMetadata(id int64) (*NodeMetadata, error)
}
