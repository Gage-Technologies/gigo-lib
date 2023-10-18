package models

import (
	"encoding/gob"
	"github.com/gage-technologies/gigo-lib/db/models"
)

// we have to register the custom types with gob
// so that they can be marshaled and unmarshaled
func init() {
	gob.Register(&NewChatMsg{})
	gob.Register(&NewMessageMsg{})
	gob.Register(&ChatKickMsg{})
	gob.Register(&ChatUpdatedEventMsg{})
}

type ChatUpdateEvent int

const (
	ChatUpdateEventNameChange ChatUpdateEvent = iota
	ChatUpdateEventUserAdd
	ChatUpdateEventUserRemove
	ChatUpdateEventDeleted
)

type NewChatMsg struct {
	Chat models.Chat `json:"chat"`
}

type NewMessageMsg struct {
	Message models.ChatMessage `json:"message"`
}

type ChatKickMsg struct {
	ChatID int64 `json:"chat_id"`
}

type ChatUpdatedEventMsg struct {
	Chat                           *models.ChatFrontend `json:"chat"`
	UpdateEvents                   []ChatUpdateEvent    `json:"update_events"`
	OldName                        string               `json:"old_name"`
	AddedUsers                     []string             `json:"added_users"`
	RemovedUsers                   []string             `json:"removed_users"`
	Updater                        string               `json:"updater"`
	UpdaterIcon                    string               `json:"updater_icon"`
	UpdaterBackground              string               `json:"updater_background"`
	UpdaterBackgroundRenderInFront bool                 `json:"updater_background_render_in_front"`
	UpdaterBackgroundPalette       string               `json:"updater_background_palette"`
	UpdaterPro                     bool                 `json:"updater_pro"`
}
