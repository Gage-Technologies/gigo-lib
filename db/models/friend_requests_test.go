package models

import (
	"fmt"
	ti "github.com/gage-technologies/gigo-lib/db"
	"testing"
	"time"
)

func TestCreateFriendRequests(t *testing.T) {
	date := time.Date(1, 0, 0, 0, 0, 0, 0, time.UTC)
	friend, err := CreateFriendRequests(1, 69420, "userName", 6969, "name", date, 5)
	if err != nil {
		t.Error("\nCreate FriendRequests Failed")
		return
	}

	if friend == nil {
		t.Error("\nCreate FriendRequests Failed]\n    Error: creation returned nil")
		return
	}

	if friend.UserID != 69420 {
		t.Error("\nCreate FriendRequests Failed]\n    Error: wrong id")
		return
	}

	if friend.Friend != 6969 {
		t.Error("\nCreate FriendRequests Failed]\n    Error: wrong author id")
		return
	}

	if friend.Response != nil {
		t.Error("\nCreate FriendRequests Failed]\n    Error: wrong response")
		return
	}

	if friend.Date != date {
		t.Error("\nCreate FriendRequests Failed]\n    Error: wrong date")
		return
	}

	t.Log("\nCreate FriendRequests Succeeded")
}

func TestFriendRequests_ToSQLNative(t *testing.T) {
	date := time.Date(1, 0, 0, 0, 0, 0, 0, time.UTC)
	fr, err := CreateFriendRequests(1, 69420, "userName", 6969, "name", date, 5)
	if err != nil {
		t.Error("\nCreate FriendRequests Failed")
		return
	}
	statement := fr.ToSQLNative()

	if statement.Statement != "insert ignore into friend_requests(_id, user_id, user_name, friend, friend_name, response, date, notification_id) values(?, ?, ?, ?, ?, ?, ?, ?);" {
		t.Errorf("\nRec post to sql native failed\n    Error: incorrect statement returned")
		return
	}

	if len(statement.Values) != 8 {
		fmt.Println("number of values returned: ", len(statement.Values))
		t.Errorf("\nRec post to sql native failed\n    Error: incorrect values returned %v", statement.Values)
		return
	}

	t.Log("\nRec post to sql native succeeded")
}

func TestFriendRequestsFromSQLNative(t *testing.T) {
	db, err := ti.CreateDatabase("gigo-dev-tidb", "4000", "mysql", "gigo-dev", "gigo-dev", "gigo_dev_test")
	if err != nil {
		t.Error("\nInitialize FriendRequests table sql failed\n    Error: ", err)
		return
	}

	defer db.DB.Exec("delete from friend_requests")

	date := time.Date(1999, 5, 19, 0, 0, 0, 0, time.UTC)
	friend, err := CreateFriendRequests(1, 69420, "userName", 6969, "name", date, 5)
	if err != nil {
		t.Error("\nCreate FriendRequests Failed")
		return
	}

	statement := friend.ToSQLNative()

	stmt, err := db.DB.Prepare(statement.Statement)
	if err != nil {
		t.Error("\nRec post from sql native failed\n    Error: ", err)
		return
	}

	_, err = stmt.Exec(statement.Values...)
	if err != nil {
		t.Error("\nRec post from sql native failed\n    Error: ", err)
		return
	}

	rows, err := db.DB.Query("select * from friend_requests")
	if err != nil {
		t.Error("\nRec post from sql native failed\n    Error: ", err)
		return
	}

	if !rows.Next() {
		t.Error("\nRec post from sql native failed\n    Error: no rows found")
		return
	}

	rec, err := FriendRequestsFromSQLNative(db, rows)
	if err != nil {
		t.Error("\nRec post from sql native failed\n    Error: ", err)
		return
	}

	if rec == nil {
		t.Error("\nCreate FriendRequests Failed]\n    Error: creation returned nil")
		return
	}

	if rec.UserID != 69420 {
		t.Error("\nCreate FriendRequests Failed]\n    Error: wrong id")
		return
	}

	if rec.Friend != 6969 {
		t.Error("\nCreate FriendRequests Failed]\n    Error: wrong author id")
		return
	}

	if rec.Response != nil {
		t.Error("\nCreate FriendRequests Failed]\n    Error: wrong response")
		return
	}

	if rec.Date != date {
		t.Error("\nCreate FriendRequests Failed]\n    Error: wrong date")
		return
	}

}
