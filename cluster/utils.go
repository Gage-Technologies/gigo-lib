package cluster

import (
	"fmt"
	"strconv"
	"strings"
)

// extractNodeIDFromKey
//
//		Extracts the node ID from a key and returns the nodes id
//	 and the passed key without the node id
func extractNodeIDFromKey(key string) (int64, string, error) {
	// attempt to format the last part of the key to a node id
	parts := strings.Split(key, "/")
	if len(parts) < 2 {
		return -1, "", fmt.Errorf("failed to parse key: %s", key)
	}

	// separate the node id from the key
	nodeIdString := parts[len(parts)-1]
	keyBase := strings.Join(parts[:len(parts)-1], "/")

	// get node id
	id, err := strconv.ParseInt(nodeIdString, 10, 64)
	if err != nil {
		return -1, "", fmt.Errorf("failed to parse id from key: %s", key)
	}

	return id, keyBase, nil
}
