package models

import (
	"encoding/gob"
)

// we have to register the custom types with gob
// so that they can be marshaled and unmarshaled
func init() {
	gob.Register(&NewWeekInactivityMsg{})
	gob.Register(&NewMonthInactivityMsg{})
}

type NewWeekInactivityMsg struct {
	Recipient string `json:"recipient"`
}

type NewMonthInactivityMsg struct {
	Recipient string `json:"recipient"`
}
