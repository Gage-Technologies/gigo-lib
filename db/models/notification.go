package models

import (
	"database/sql"
	"fmt"
	"github.com/kisielk/sqlstruct"
	"time"
)

type NotificationType int

const (
	FriendRequest NotificationType = iota
	NemesisRequest
	NemesisAlert
	StreakInfo
)

func (b NotificationType) String() string {
	switch b {
	case FriendRequest:
		return "friend_request"
	case NemesisRequest:
		return "nemesis_request"
	case NemesisAlert:
		return "nemesis_alert"
	case StreakInfo:
		return "streak_info"
	}
	return "Unknown"
}

type Notification struct {
	ID                int64            `json:"_id" sql:"_id"`
	UserID            int64            `json:"user_id" sql:"user_id"`
	Message           string           `json:"message" sql:"message"`
	NotificationType  NotificationType `json:"notification_type" sql:"notification_type"`
	CreatedAt         time.Time        `json:"created_at" sql:"created_at"`
	Acknowledged      bool             `json:"acknowledged" sql:"acknowledged"`
	InteractingUserID *int64           `json:"interacting_user_id" sql:"interacting_user_id"`
}

type NotificationSQL struct {
	ID                int64            `json:"_id" sql:"_id"`
	UserID            int64            `json:"user_id" sql:"user_id"`
	Message           string           `json:"message" sql:"message"`
	NotificationType  NotificationType `json:"notification_type" sql:"notification_type"`
	CreatedAt         time.Time        `json:"created_at" sql:"created_at"`
	Acknowledged      bool             `json:"acknowledged" sql:"acknowledged"`
	InteractingUserID *int64           `json:"interacting_user_id" sql:"interacting_user_id"`
}

type NotificationFrontend struct {
	ID                string           `json:"_id" sql:"_id"`
	UserID            string           `json:"user_id" sql:"user_id"`
	Message           string           `json:"message" sql:"message"`
	NotificationType  NotificationType `json:"notification_type" sql:"notification_type"`
	CreatedAt         time.Time        `json:"created_at" sql:"created_at"`
	Acknowledged      bool             `json:"acknowledged" sql:"acknowledged"`
	InteractingUserID *string          `json:"interacting_user_id" sql:"interacting_user_id"`
}

func CreateNotification(id int64, userID int64, message string, notificationType NotificationType, createdAt time.Time, acknowledged bool, interactingUser *int64) (*Notification, error) {

	return &Notification{
		ID:                id,
		UserID:            userID,
		Message:           message,
		NotificationType:  notificationType,
		CreatedAt:         createdAt,
		Acknowledged:      acknowledged,
		InteractingUserID: interactingUser,
	}, nil
}

func NotificationFromSQLNative(rows *sql.Rows) (*Notification, error) {
	// create new Notification object to load into
	notificationSql := new(NotificationSQL)

	// scan row into Notification object
	err := sqlstruct.Scan(notificationSql, rows)
	if err != nil {
		return nil, err
	}

	// create new Notification for the output
	notification := &Notification{
		ID:                notificationSql.ID,
		UserID:            notificationSql.UserID,
		Message:           notificationSql.Message,
		NotificationType:  notificationSql.NotificationType,
		CreatedAt:         notificationSql.CreatedAt,
		Acknowledged:      notificationSql.Acknowledged,
		InteractingUserID: notificationSql.InteractingUserID,
	}

	return notification, nil
}

func (i *Notification) ToFrontend() *NotificationFrontend {
	// conditionally format top reply id into string
	var interactingUser string
	if i.InteractingUserID != nil {
		interactingUser = fmt.Sprintf("%d", *i.InteractingUserID)
	}

	// create new BroadcastEvent frontend
	mf := &NotificationFrontend{
		ID:                fmt.Sprintf("%d", i.ID),
		UserID:            fmt.Sprintf("%d", i.UserID),
		Message:           i.Message,
		NotificationType:  i.NotificationType,
		CreatedAt:         i.CreatedAt,
		Acknowledged:      i.Acknowledged,
		InteractingUserID: &interactingUser,
	}

	return mf
}

func (i *Notification) ToSQLNative() *SQLInsertStatement {

	// create insertion statement and return
	return &SQLInsertStatement{
		Statement: "insert ignore into notification(_id, user_id, message, notification_type, created_at, acknowledged, interacting_user_id) values(?, ?, ?, ?, ?, ?, ?);",
		Values:    []interface{}{i.ID, i.UserID, i.Message, i.NotificationType, i.CreatedAt, i.Acknowledged, i.InteractingUserID},
	}
}
