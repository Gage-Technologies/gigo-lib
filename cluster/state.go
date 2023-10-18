package cluster

import (
	"context"
	"fmt"
)

const (
	ElectionPrefix  = "election"
	StateDataPrefix = "state-data"
	NodesPrefix     = "nodes"
)

// LeaderRoutine
//
//	Called by the leader every 50ms. The function should
//	exit as quickly as possible on a context cancel and
//	encompass an entire execution cycle for the leader.
//	Any logic container in this function that should not
//	be executed as frequently as 50ms should track the time
//	between executions independently and skip cycles that
//	are not appropriate for such logic.
type LeaderRoutine func(ctx context.Context) error

// FollowerRoutine
//
//	Called by the follower every 50ms. The function should
//	exit as quickly as possible on a context cancel and
//	encompass an entire execution cycle for the follower.
//	Any logic container in this function that should not
//	be executed as frequently as 50ms should track the time
//	between executions independently and skip cycles that
//	are not appropriate for such logic.
type FollowerRoutine func(ctx context.Context) error

// EventType
//
//	Type of event that occurred to the state
type EventType int

const (
	// EventTypeAdded value has been added to the state and previously did not exist
	EventTypeAdded EventType = iota
	// EventTypeModified value has been modified in the state and previously existed
	EventTypeModified
	// EventTypeDeleted value has been removed from the state and previously existed
	EventTypeDeleted
)

func (t EventType) String() string {
	switch t {
	case EventTypeAdded:
		return "Added"
	case EventTypeModified:
		return "Modified"
	case EventTypeDeleted:
		return "Deleted"
	default:
		return "Unknown"
	}
}

// StateChangeEvent
//
//	An event in which the state of state kv pair changes
type StateChangeEvent struct {
	Type     EventType
	NodeID   int64
	Key      string
	Value    string
	OldKey   string
	OldValue string
}

type NodeRole int

const (
	NodeRoleUnknown NodeRole = iota
	NodeRoleLeader
	NodeRoleFollower
)

func (s NodeRole) String() string {
	switch s {
	case NodeRoleUnknown:
		return "Unknown"
	case NodeRoleLeader:
		return "Leader"
	case NodeRoleFollower:
		return "Follower"
	default:
		return "Invalid"
	}
}

// KV
//
// Simple key-value pair
type KV struct {
	Key   string
	Value string
}

// formatPrefix
//
//	Helper function to format a prefix string
//	using the cluster name, prefix and node id
func formatPrefix(id int64, cluster string, prefix string) string {
	return fmt.Sprintf("/%s/%s/%d", cluster, prefix, id)
}

// formatPrefixCluster
//
//	Helper function to format a prefix string
//	using the cluster name and prefix. This function
//	excludes the node id so that keys for all nodes
//	can be retrieved.
func formatPrefixCluster(cluster string, prefix string) string {
	return fmt.Sprintf("/%s/%s/", cluster, prefix)
}
