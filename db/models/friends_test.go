package models

import (
	"fmt"
	ti "github.com/gage-technologies/gigo-lib/db"
	"testing"
	"time"
)

func TestCreateFriends(t *testing.T) {
	date := time.Date(1, 0, 0, 0, 0, 0, 0, time.UTC)
	friend, err := CreateFriends(1, 69420, "username", 6969, "friendname", date)
	if err != nil {
		t.Error("\nCreate Friends Failed")
		return
	}

	if friend == nil {
		t.Error("\nCreate Friends Failed]\n    Error: creation returned nil")
		return
	}

	if friend.UserID != 69420 {
		t.Error("\nCreate Friends Failed]\n    Error: wrong id")
		return
	}

	if friend.Friend != 6969 {
		t.Error("\nCreate Friends Failed]\n    Error: wrong author id")
		return
	}

	if friend.Date != date {
		t.Error("\nCreate Friends Failed]\n    Error: wrong date")
		return
	}

	if friend.UserName != "username" {
		t.Error("\nCreate Friends Failed]\n    Error: wrong author id")
		return
	}

	if friend.FriendName != "friendname" {
		t.Error("\nCreate Friends Failed]\n    Error: wrong author id")
		return
	}

	t.Log("\nCreate Friends Succeeded")
}

func TestFriends_ToSQLNative(t *testing.T) {
	date := time.Date(1, 0, 0, 0, 0, 0, 0, time.UTC)
	fr, err := CreateFriends(1, 69420, "username", 6969, "friendname", date)
	if err != nil {
		t.Error("\nCreate Friends Failed")
		return
	}
	statement := fr.ToSQLNative()

	if statement.Statement != "insert ignore into friends(_id, user_id, user_name, friend, friend_name, date) values(?, ?, ?, ?, ?, ?);" {
		t.Errorf("\nRec post to sql native failed\n    Error: incorrect statement returned")
		return
	}

	if len(statement.Values) != 6 {
		fmt.Println("number of values returned: ", len(statement.Values))
		t.Errorf("\nRec post to sql native failed\n    Error: incorrect values returned %v", statement.Values)
		return
	}

	t.Log("\nRec post to sql native succeeded")
}

func TestFriendsFromSQLNative(t *testing.T) {
	db, err := ti.CreateDatabase("gigo-dev-tidb", "4000", "mysql", "gigo-dev", "gigo-dev", "gigo_dev_test")
	if err != nil {
		t.Error("\nInitialize friends table sql failed\n    Error: ", err)
		return
	}

	defer db.DB.Exec("DROP TABLE friends")

	date := time.Date(1999, 5, 19, 0, 0, 0, 0, time.UTC)
	friend, err := CreateFriends(1, 69420, "username", 6969, "friendname", date)
	if err != nil {
		t.Error("\nCreate Friends Failed")
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

	rows, err := db.DB.Query("select * from friends")
	if err != nil {
		t.Error("\nRec post from sql native failed\n    Error: ", err)
		return
	}

	if !rows.Next() {
		t.Error("\nRec post from sql native failed\n    Error: no rows found")
		return
	}

	rec, err := FriendsFromSQLNative(db, rows)
	if err != nil {
		t.Error("\nRec post from sql native failed\n    Error: ", err)
		return
	}

	if rec == nil {
		t.Error("\nCreate Friends Failed]\n    Error: creation returned nil")
		return
	}

	if rec.UserID != 69420 {
		t.Error("\nCreate Friends Failed]\n    Error: wrong id")
		return
	}

	if rec.Friend != 6969 {
		t.Error("\nCreate Friends Failed]\n    Error: wrong author id")
		return
	}

	if rec.Date != date {
		t.Error("\nCreate Friends Failed]\n    Error: wrong date")
		return
	}

	if friend.UserName != "username" {
		t.Error("\nCreate Friends Failed]\n    Error: wrong author id")
		return
	}

	if friend.FriendName != "friendname" {
		t.Error("\nCreate Friends Failed]\n    Error: wrong author id")
		return
	}

}
