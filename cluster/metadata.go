package cluster

import (
	"encoding/json"
	"time"
)

// NodeMetadata
//
//	Metadata for a node that is persisted in the cluster state
type NodeMetadata struct {
	ID      int64
	Address string
	Start   time.Time
	Role    NodeRole
}

func UnmarshalNodeMetadata(data []byte) (NodeMetadata, error) {
	metadata := NodeMetadata{}
	err := json.Unmarshal(data, &metadata)
	if err != nil {
		return NodeMetadata{}, err
	}
	return metadata, nil
}

func (m *NodeMetadata) Marshal() ([]byte, error) {
	return json.Marshal(m)
}
