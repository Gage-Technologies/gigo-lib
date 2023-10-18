package models

import (
	"context"
	"database/sql"
	"fmt"
	ti "github.com/gage-technologies/gigo-lib/db"
	"github.com/kisielk/sqlstruct"
	"go.opentelemetry.io/otel"
	"time"
)

type ChatType int

const (
	ChatTypeGlobal ChatType = iota
	ChatTypeRegional
	ChatTypePublicGroup
	ChatTypePrivateGroup
	ChatTypeDirectMessage
	ChatTypeChallenge
)

type ChatIconInfo struct {
	Icon                    string `json:"icon" sql:"icon"`
	Background              string `json:"background" sql:"background"`
	BackgroundRenderInFront bool   `json:"background_render_in_front" sql:"background_render_in_front"`
	BackgroundPalette       string `json:"background_palette" sql:"background_palette"`
	Pro                     bool   `json:"pro" sql:"pro"`
}

type Chat struct {
	ID              int64         `json:"_id" sql:"_id"`
	Name            string        `json:"name" sql:"name"`
	Type            ChatType      `json:"type" sql:"type"`
	Users           []int64       `json:"users" sql:"users"`
	Usernames       []string      `json:"user_names" sql:"user_names"`
	LastMessage     *int64        `json:"last_message" sql:"last_message"`
	LastMessageTime *time.Time    `json:"last_message_time" sql:"last_message_time"`
	ChatIconInfo    *ChatIconInfo `json:"icon" sql:"icon"`
	LastReadMessage *int64        `json:"last_read_message" sql:"last_read_message"`
	Muted           bool          `json:"muted" sql:"muted"`
}

type ChatSQL struct {
	ID              int64      `json:"_id" sql:"_id"`
	Name            string     `json:"name" sql:"name"`
	Type            ChatType   `json:"type" sql:"type"`
	Users           []int64    `json:"users" sql:"users"`
	Usernames       []string   `json:"user_names" sql:"user_names"`
	LastMessage     *int64     `json:"last_message" sql:"last_message"`
	LastMessageTime *time.Time `json:"last_message_time" sql:"last_message_time"`
}

type ChatFrontend struct {
	ID              string        `json:"_id" sql:"_id"`
	Name            string        `json:"name" sql:"name"`
	Type            ChatType      `json:"type" sql:"type"`
	Users           []string      `json:"users" sql:"users"`
	UserNames       []string      `json:"user_names" sql:"user_names"`
	LastMessage     *string       `json:"last_message" sql:"last_message"`
	LastMessageTime *time.Time    `json:"last_message_time" sql:"last_message_time"`
	ChatIconInfo    *ChatIconInfo `json:"icon" sql:"icon"`
	LastReadMessage *string       `json:"last_read_message" sql:"last_read_message"`
	Muted           bool          `json:"muted" sql:"muted"`
}

func CreateChat(id int64, name string, chatType ChatType, users []int64) *Chat {
	now := time.Now()
	return &Chat{
		ID:              id,
		Name:            name,
		Type:            chatType,
		Users:           users,
		LastMessageTime: &now,
	}
}

func ChatFromSQLNative(callerId int64, db *ti.Database, rows *sql.Rows) (*Chat, error) {
	// create new Chat object to load into
	chatSql := new(ChatSQL)

	// scan row into Chat object
	err := sqlstruct.Scan(chatSql, rows)
	if err != nil {
		return nil, err
	}

	// query for the users
	var users []int64
	var userNames []string

	// create span for telemetry
	ctx, span := otel.Tracer("gigo-core").Start(context.Background(), "gigo-lib")
	defer span.End()
	callerName := "ChatFromSQLNative"

	var lastReadMessage *int64
	var muted bool

	// load the users if this isn't a global, regional, or challenge chat
	if chatSql.Type != ChatTypeGlobal && chatSql.Type != ChatTypeRegional && chatSql.Type != ChatTypeChallenge {
		// query for all the users in this chat
		rows, err = db.Query(ctx, &span, &callerName, "select cu.user_id, u.user_name, cu.last_read_message, cu.muted from chat_users cu join users u on u._id = cu.user_id where cu.chat_id = ?", chatSql.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to query for users for chat %d: %w", chatSql.ID, err)
		}

		defer rows.Close()

		// iterate over the rows loading the user ids and names
		for rows.Next() {
			var userId int64
			var userName string
			var localLastReadId *int64
			var localMuted bool

			err = rows.Scan(&userId, &userName, &localLastReadId, &localMuted)
			if err != nil {
				return nil, fmt.Errorf("failed to scan user for chat %d: %w", chatSql.ID, err)
			}

			users = append(users, userId)
			userNames = append(userNames, userName)

			// conditionally set the last read message and muted if this is the caller
			if userId == callerId {
				lastReadMessage = localLastReadId
				muted = localMuted
			}
		}
	}

	// if this is a dm then we need to load the chat icon info
	var chatIconInfo *ChatIconInfo
	if chatSql.Type == ChatTypeDirectMessage && callerId != 0 && len(users) == 2 {
		// select the user that is not the caller
		var otherUser int64
		if users[0] == callerId {
			otherUser = users[1]
		} else {
			otherUser = users[0]
		}

		// create a new chat icon info
		chatIconInfo = &ChatIconInfo{
			Icon: fmt.Sprintf("/static/user/pfp/%v", otherUser),
		}

		// query for the user's icon info
		err = db.QueryRow(
			ctx, &span, &callerName,
			"select r.name as background, r.color_palette as background_palette, r.render_in_front as render_in_front, u.user_status = 1 as pro from users u left join rewards r on r._id = u.avatar_reward  where u._id = ?",
			otherUser,
		).Scan(
			&chatIconInfo.Background,
			&chatIconInfo.BackgroundPalette,
			&chatIconInfo.BackgroundRenderInFront,
			&chatIconInfo.Pro,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to query for user %d: %w", otherUser, err)
		}
	}

	return &Chat{
		ID:              chatSql.ID,
		Name:            chatSql.Name,
		Type:            chatSql.Type,
		Users:           users,
		Usernames:       userNames,
		LastMessage:     chatSql.LastMessage,
		LastMessageTime: chatSql.LastMessageTime,
		ChatIconInfo:    chatIconInfo,
		LastReadMessage: lastReadMessage,
		Muted:           muted,
	}, nil
}

func (c *Chat) ToFrontend() *ChatFrontend {
	// conditionally format the last message and last read message
	var lastMessage *string
	if c.LastMessage != nil {
		lastMessageStr := fmt.Sprintf("%d", *c.LastMessage)
		lastMessage = &lastMessageStr
	}

	var lastReadMessage *string
	if c.LastReadMessage != nil {
		lastReadMessageStr := fmt.Sprintf("%d", *c.LastReadMessage)
		lastReadMessage = &lastReadMessageStr
	}

	// create new Chat frontend
	mf := &ChatFrontend{
		ID:              fmt.Sprintf("%d", c.ID),
		Name:            c.Name,
		Type:            c.Type,
		Users:           make([]string, 0),
		UserNames:       c.Usernames,
		LastMessage:     lastMessage,
		LastMessageTime: c.LastMessageTime,
		ChatIconInfo:    c.ChatIconInfo,
		LastReadMessage: lastReadMessage,
		Muted:           c.Muted,
	}

	// iterate over the users loading them into the frontend
	for _, user := range c.Users {
		mf.Users = append(mf.Users, fmt.Sprintf("%d", user))
	}

	return mf
}

func (c *Chat) ToSQLNative() []SQLInsertStatement {
	// create the insert statement
	inserts := []SQLInsertStatement{
		{
			Statement: "insert ignore into chat(_id, name, type, last_message_time) values(?, ?, ?, ?);",
			Values: []interface{}{
				c.ID,
				c.Name,
				c.Type,
				c.LastMessageTime,
			},
		},
	}

	// create the insert statements for the users
	for _, user := range c.Users {
		inserts = append(inserts, SQLInsertStatement{
			Statement: "insert ignore into chat_users(chat_id, user_id) values(?, ?);",
			Values:    []interface{}{c.ID, user},
		})
	}

	return inserts
}
