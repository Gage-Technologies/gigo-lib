package models

import (
	"database/sql"
	"fmt"
	"github.com/kisielk/sqlstruct"
	"time"
)

type BroadcastType int

const (
	BroadcastMessage BroadcastType = iota
	BroadcastNotification
)

func (b BroadcastType) String() string {
	switch b {
	case BroadcastMessage:
		return "broadcast_message"
	case BroadcastNotification:
		return "broadcast_notification"
	}
	return "Unknown"
}

type BroadcastEvent struct {
	ID            int64         `json:"_id" sql:"_id"`
	UserID        int64         `json:"user_id" sql:"user_id"`
	UserName      string        `json:"user_name" sql:"user_name"`
	Message       string        `json:"message" sql:"message"`
	BroadcastType BroadcastType `json:"broadcast_type" sql:"broadcast_type"`
	TimePosted    time.Time     `json:"time_posted" sql:"time_posted"`
}

type BroadcastEventSQL struct {
	ID            int64         `json:"_id" sql:"_id"`
	UserID        int64         `json:"user_id" sql:"user_id"`
	UserName      string        `json:"user_name" sql:"user_name"`
	Message       string        `json:"message" sql:"message"`
	BroadcastType BroadcastType `json:"broadcast_type" sql:"broadcast_type"`
	TimePosted    time.Time     `json:"time_posted" sql:"time_posted"`
}

type BroadcastEventFrontend struct {
	ID            string        `json:"_id" sql:"_id"`
	UserID        string        `json:"user_id" sql:"user_id"`
	UserName      string        `json:"user_name" sql:"user_name"`
	Message       string        `json:"message" sql:"message"`
	BroadcastType BroadcastType `json:"broadcast_type" sql:"broadcast_type"`
	TimePosted    time.Time     `json:"time_posted" sql:"time_posted"`
}

func CreateBroadcastEvent(id int64, userID int64, userName string, message string, broadcastType BroadcastType, date time.Time) (*BroadcastEvent, error) {

	return &BroadcastEvent{
		ID:            id,
		UserID:        userID,
		UserName:      userName,
		Message:       message,
		BroadcastType: broadcastType,
		TimePosted:    date,
	}, nil
}

func BroadcastEventFromSQLNative(rows *sql.Rows) (*BroadcastEvent, error) {
	// create new coffee object to load into
	broadcastSql := new(BroadcastEventSQL)

	// scan row into coffee object
	err := sqlstruct.Scan(broadcastSql, rows)
	if err != nil {
		return nil, err
	}

	// create new coffee for the output
	event := &BroadcastEvent{
		ID:            broadcastSql.ID,
		UserID:        broadcastSql.UserID,
		UserName:      broadcastSql.UserName,
		Message:       broadcastSql.Message,
		BroadcastType: broadcastSql.BroadcastType,
		TimePosted:    broadcastSql.TimePosted,
	}

	return event, nil
}

func (i *BroadcastEvent) ToFrontend() *BroadcastEventFrontend {

	// create new BroadcastEvent frontend
	mf := &BroadcastEventFrontend{
		ID:            fmt.Sprintf("%d", i.ID),
		UserID:        fmt.Sprintf("%d", i.UserID),
		UserName:      i.UserName,
		Message:       i.Message,
		BroadcastType: i.BroadcastType,
		TimePosted:    i.TimePosted,
	}

	return mf
}

func (i *BroadcastEvent) ToSQLNative() *SQLInsertStatement {

	// create insertion statement and return
	return &SQLInsertStatement{
		Statement: "insert ignore into broadcast_event(_id, user_id, user_name, message, broadcast_type, time_posted) values(?, ?, ?, ?, ?, ?);",
		Values:    []interface{}{i.ID, i.UserID, i.UserName, i.Message, i.BroadcastType, i.TimePosted},
	}
}
