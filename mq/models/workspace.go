package models

import "encoding/gob"

// we have to register the custom types with gob
// so that they can be marshaled and unmarshaled
func init() {
	gob.Register(&CreateWorkspaceMsg{})
	gob.Register(&StopWorkspaceMsg{})
	gob.Register(&DestroyWorkspaceMsg{})
}

type CreateWorkspaceMsg struct {
	WorkspaceID int64
	OwnerID     int64
	OwnerEmail  string
	OwnerName   string
	Disk        int
	CPU         int
	Memory      int
	Container   string
	AccessUrl   string
}

type StartWorkspaceMsg struct {
	ID int64
}

type StopWorkspaceMsg struct {
	ID              int64
	OwnerID         int64
	WorkspaceFailed bool
}

type DestroyWorkspaceMsg struct {
	ID      int64
	OwnerID int64
}
