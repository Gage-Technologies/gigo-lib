package models

import (
	"database/sql"
	"fmt"
	ti "github.com/gage-technologies/gigo-lib/db"
	"github.com/kisielk/sqlstruct"
	"time"
)

type FriendRequests struct {
	ID             int64     `json:"_id"`
	UserID         int64     `json:"user_id" sql:"user_id"`
	UserName       string    `json:"user_name" sql:"user_name"`
	Friend         int64     `json:"friend" sql:"friend"`
	FriendName     string    `json:"friend_name" sql:"friend_name"`
	Response       *bool     `json:"response" sql:"response"`
	Date           time.Time `json:"date" sql:"date"`
	NotificationID int64     `json:"notification_id" sql:"notification_id"`
}

type FriendRequestsSQL struct {
	ID             int64     `json:"_id"`
	UserID         int64     `json:"user_id" sql:"user_id"`
	UserName       string    `json:"user_name" sql:"user_name"`
	Friend         int64     `json:"friend" sql:"friend"`
	FriendName     string    `json:"friend_name" sql:"friend_name"`
	Response       *bool     `json:"response" sql:"response"`
	Date           time.Time `json:"date" sql:"date"`
	NotificationID int64     `json:"notification_id" sql:"notification_id"`
}

type FriendRequestsFrontend struct {
	ID             string    `json:"_id"`
	UserID         string    `json:"user_id" sql:"user_id"`
	UserName       string    `json:"user_name" sql:"user_name"`
	Friend         string    `json:"friend" sql:"friend"`
	FriendName     string    `json:"friend_name" sql:"friend_name"`
	Response       *bool     `json:"response" sql:"response"`
	Date           time.Time `json:"date" sql:"date"`
	NotificationID string    `json:"notification_id" sql:"notification_id"`
}

func CreateFriendRequests(id int64, userID int64, userName string, friend int64, friendName string, date time.Time, notificationId int64) (*FriendRequests, error) {

	return &FriendRequests{
		ID:             id,
		UserID:         userID,
		UserName:       userName,
		Friend:         friend,
		FriendName:     friendName,
		Response:       nil,
		Date:           date,
		NotificationID: notificationId,
	}, nil
}

func FriendRequestsFromSQLNative(db *ti.Database, rows *sql.Rows) (*FriendRequests, error) {
	// create new friend request object to load into
	requestSQL := new(FriendRequestsSQL)

	// scan row into friend requests object
	err := sqlstruct.Scan(requestSQL, rows)
	if err != nil {
		return nil, err
	}

	// create new friend for the output
	request := &FriendRequests{
		ID:             requestSQL.ID,
		UserID:         requestSQL.UserID,
		UserName:       requestSQL.UserName,
		Friend:         requestSQL.Friend,
		FriendName:     requestSQL.FriendName,
		Response:       requestSQL.Response,
		Date:           requestSQL.Date,
		NotificationID: requestSQL.NotificationID,
	}

	return request, nil
}

func (i *FriendRequests) ToFrontend() *FriendRequestsFrontend {

	// create new friend requests frontend
	r := &FriendRequestsFrontend{
		ID:             fmt.Sprintf("%d", i.ID),
		UserID:         fmt.Sprintf("%d", i.UserID),
		UserName:       i.UserName,
		Friend:         fmt.Sprintf("%d", i.Friend),
		FriendName:     i.FriendName,
		Response:       i.Response,
		Date:           i.Date,
		NotificationID: fmt.Sprintf("%d", i.NotificationID),
	}

	return r
}

func (i *FriendRequests) ToSQLNative() *SQLInsertStatement {

	sqlStatements := &SQLInsertStatement{
		Statement: "insert ignore into friend_requests(_id, user_id, user_name, friend, friend_name, response, date, notification_id) values(?, ?, ?, ?, ?, ?, ?, ?);",
		Values:    []interface{}{i.ID, i.UserID, i.UserName, i.Friend, i.FriendName, i.Response, i.Date, i.NotificationID},
	}

	// create insertion statement and return
	return sqlStatements
}
