package models

import "encoding/gob"

// we have to register the custom types with gob
// so that they can be marshaled and unmarshaled
func init() {
	gob.Register(&ForgetConnMsg{})
}

type ForgetConnMsg struct {
	WorkspaceID int64 `json:"workspace_id"`
	AgentID     int64 `json:"conn_id"`
}
