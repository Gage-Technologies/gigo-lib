package models

import (
	"encoding/gob"

	"github.com/gage-technologies/gigo-lib/db/models"
)

// we have to register the custom types with gob
// so that they can be marshaled and unmarshaled
func init() {
	gob.Register(&WorkspaceStatusUpdateMsg{})
	gob.Register(&WorkspaceResourceUtil{})
}

type WorkspaceResourceUtil struct {
	CPU         float64 `json:"cpu"`
	Memory      float64 `json:"memory"`
	CPULimit    int64   `json:"cpu_limit"`
	MemoryLimit int64   `json:"memory_limit"`
	CPUUsage    int64   `json:"cpu_usage"`
	MemoryUsage int64   `json:"memory_usage"`
}

type WorkspaceStatusUpdateMsg struct {
	Workspace *models.WorkspaceFrontend `json:"workspace"`
	Resources *WorkspaceResourceUtil    `json:"resources"`
}
