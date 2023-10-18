package models

import (
	"database/sql"
	"fmt"
	ti "github.com/gage-technologies/gigo-lib/db"
	"github.com/kisielk/sqlstruct"
	"time"
)

type Friends struct {
	ID         int64     `json:"_id"`
	UserID     int64     `json:"user_id" sql:"user_id"`
	UserName   string    `json:"user_name" sql:"user_name"`
	Friend     int64     `json:"friend" sql:"friend"`
	FriendName string    `json:"friend_name" sql:"friend_name"`
	Date       time.Time `json:"date" sql:"date"`
}

type FriendsSQL struct {
	ID         int64     `json:"_id"`
	UserID     int64     `json:"user_id" sql:"user_id"`
	UserName   string    `json:"user_name" sql:"user_name"`
	Friend     int64     `json:"friend" sql:"friend"`
	FriendName string    `json:"friend_name" sql:"friend_name"`
	Date       time.Time `json:"date" sql:"date"`
}

type FriendsFrontend struct {
	ID         string    `json:"_id"`
	UserID     string    `json:"user_id" sql:"user_id"`
	UserName   string    `json:"user_name" sql:"user_name"`
	Friend     string    `json:"friend" sql:"friend"`
	FriendName string    `json:"friend_name" sql:"friend_name"`
	Date       time.Time `json:"date" sql:"date"`
}

func CreateFriends(id int64, userID int64, username string, friend int64, friendName string, date time.Time) (*Friends, error) {

	return &Friends{
		ID:         id,
		UserID:     userID,
		UserName:   username,
		Friend:     friend,
		FriendName: friendName,
		Date:       date,
	}, nil
}

func FriendsFromSQLNative(db *ti.Database, rows *sql.Rows) (*Friends, error) {
	// create new friends object to load into
	friendsSQL := new(FriendsSQL)

	// scan row into attempt object
	err := sqlstruct.Scan(friendsSQL, rows)
	if err != nil {
		return nil, err
	}

	// create new attempt for the output
	friend := &Friends{
		ID:         friendsSQL.ID,
		UserID:     friendsSQL.UserID,
		UserName:   friendsSQL.UserName,
		Friend:     friendsSQL.Friend,
		FriendName: friendsSQL.FriendName,
		Date:       friendsSQL.Date,
	}

	return friend, nil
}

func (i *Friends) ToFrontend() *FriendsFrontend {

	// create new attempt frontend
	f := &FriendsFrontend{
		ID:         fmt.Sprintf("%d", i.ID),
		UserID:     fmt.Sprintf("%d", i.UserID),
		UserName:   i.UserName,
		Friend:     fmt.Sprintf("%d", i.Friend),
		FriendName: i.FriendName,
		Date:       i.Date,
	}

	return f
}

func (i *Friends) ToSQLNative() *SQLInsertStatement {

	sqlStatements := &SQLInsertStatement{
		Statement: "insert ignore into friends(_id, user_id, user_name, friend, friend_name, date) values(?, ?, ?, ?, ?, ?);",
		Values:    []interface{}{i.ID, i.UserID, i.UserName, i.Friend, i.FriendName, i.Date},
	}

	// create insertion statement and return
	return sqlStatements
}
