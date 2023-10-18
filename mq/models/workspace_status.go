package models

import (
	"encoding/gob"
	"github.com/gage-technologies/gigo-lib/db/models"
)

// we have to register the custom types with gob
// so that they can be marshaled and unmarshaled
func init() {
	gob.Register(&WorkspaceStatusUpdateMsg{})
}

type WorkspaceStatusUpdateMsg struct {
	Workspace *models.WorkspaceFrontend `json:"workspace"`
}
