package models

import "encoding/gob"

// we have to register the custom types with gob
// so that they can be marshaled and unmarshaled
func init() {
	gob.Register(&AddStreakXPMsg{})
}

type AddStreakXPMsg struct {
	ID      int64
	OwnerID int64
}
