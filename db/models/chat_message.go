package models

import (
	"database/sql"
	"fmt"
	"github.com/kisielk/sqlstruct"
	"time"
)

type ChatMessageType int

const (
	ChatMessageTypeInsecure ChatMessageType = iota
	ChatMessageTypeSecure
)

type ChatMessage struct {
	ID        int64           `json:"_id" sql:"_id"`
	ChatID    int64           `json:"chat_id" sql:"chat_id"`
	AuthorID  int64           `json:"author_id" sql:"author_id"`
	Author    string          `json:"author" sql:"author"`
	Message   string          `json:"message" sql:"message"`
	CreatedAt time.Time       `json:"created_at" sql:"created_at"`
	Revision  int64           `json:"revision" sql:"revision"`
	Type      ChatMessageType `json:"type" sql:"type"`

	// Join fields - not included in the model but common for queries
	AuthorRenown TierType `json:"author_renown" sql:"author_renown"`
}

type ChatMessageSQL struct {
	ID        int64           `json:"_id" sql:"_id"`
	ChatID    int64           `json:"chat_id" sql:"chat_id"`
	AuthorID  int64           `json:"author_id" sql:"author_id"`
	Author    string          `json:"author" sql:"author"`
	Message   string          `json:"message" sql:"message"`
	CreatedAt time.Time       `json:"created_at" sql:"created_at"`
	Revision  int64           `json:"revision" sql:"revision"`
	Type      ChatMessageType `json:"type" sql:"type"`

	// Join fields - not included in the model but common for queries
	AuthorRenown TierType `json:"author_renown" sql:"author_renown"`
}

type ChatMessageFrontend struct {
	ID        string          `json:"_id" sql:"_id"`
	ChatID    string          `json:"chat_id" sql:"chat_id"`
	AuthorID  string          `json:"author_id" sql:"author_id"`
	Author    string          `json:"author" sql:"author"`
	Message   string          `json:"message" sql:"message"`
	CreatedAt time.Time       `json:"created_at" sql:"created_at"`
	Revision  int64           `json:"revision" sql:"revision"`
	Type      ChatMessageType `json:"type" sql:"type"`

	// Join fields - not included in the model but common for queries
	AuthorRenown TierType `json:"author_renown" sql:"author_renown"`
}

func CreateChatMessage(id int64, chatID int64, authorID int64, author string, message string, createdAt time.Time, revision int64, messageType ChatMessageType) *ChatMessage {
	return &ChatMessage{
		ID:        id,
		ChatID:    chatID,
		AuthorID:  authorID,
		Author:    author,
		Message:   message,
		CreatedAt: createdAt,
		Revision:  revision,
		Type:      messageType,
	}
}

func ChatMessageFromSQLNative(rows *sql.Rows) (*ChatMessage, error) {
	// create new ChatMessage object to load into
	chatMessageSql := new(ChatMessageSQL)

	// scan row into ChatMessage object
	err := sqlstruct.Scan(chatMessageSql, rows)
	if err != nil {
		return nil, fmt.Errorf("failed to scan chat message: %w", err)
	}

	// convert to ChatMessage
	chatMessage := &ChatMessage{
		ID:           chatMessageSql.ID,
		ChatID:       chatMessageSql.ChatID,
		AuthorID:     chatMessageSql.AuthorID,
		Author:       chatMessageSql.Author,
		Message:      chatMessageSql.Message,
		CreatedAt:    chatMessageSql.CreatedAt,
		Revision:     chatMessageSql.Revision,
		Type:         chatMessageSql.Type,
		AuthorRenown: chatMessageSql.AuthorRenown,
	}

	return chatMessage, nil
}

func (i *ChatMessage) ToFrontend() *ChatMessageFrontend {
	return &ChatMessageFrontend{
		ID:           fmt.Sprintf("%d", i.ID),
		ChatID:       fmt.Sprintf("%d", i.ChatID),
		AuthorID:     fmt.Sprintf("%d", i.AuthorID),
		Author:       i.Author,
		Message:      i.Message,
		CreatedAt:    i.CreatedAt,
		Revision:     i.Revision,
		Type:         i.Type,
		AuthorRenown: i.AuthorRenown,
	}
}

func (i *ChatMessage) ToSQLNative() *SQLInsertStatement {
	return &SQLInsertStatement{
		Statement: "insert into chat_messages (_id, chat_id, author_id, author, message, created_at, revision, type) values (?, ?, ?, ?, ?, ?, ?, ?)",
		Values: []interface{}{
			i.ID,
			i.ChatID,
			i.AuthorID,
			i.Author,
			i.Message,
			i.CreatedAt,
			i.Revision,
			i.Type,
		},
	}
}
