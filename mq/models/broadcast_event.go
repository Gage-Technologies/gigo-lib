package models

import (
	"encoding/gob"
)

// we have to register the custom types with gob
// so that they can be marshaled and unmarshaled
func init() {
	gob.Register(&BroadcastMessage{})
	gob.Register(&BroadcastNotification{})
}

type BroadcastMessage struct {
	InitMessage string `json:"init_message"`
}

type BroadcastNotification struct {
	Notification string `json:"notification"`
}
