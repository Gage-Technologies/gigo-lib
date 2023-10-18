package models

import (
	"fmt"
	ti "github.com/gage-technologies/gigo-lib/db"
	"testing"
	"time"
)

func TestCreateBroadcastEvent(t *testing.T) {
	ev, err := CreateBroadcastEvent(69420, 6969, "test", "test", 0, time.Now())
	if err != nil {
		t.Error("\nCreateBroadcastEvent Failed")
		return
	}

	if ev == nil {
		t.Error("\nCreateBroadcastEvent Failed]\n    Error: creation returned nil")
		return
	}

	if ev.ID != 69420 {
		t.Error("\nCreateBroadcastEvent Failed]\n    Error: wrong id")
		return
	}

	if ev.UserID != 6969 {
		t.Error("\nCreateBroadcastEvent Failed]\n    Error: wrong author id")
		return
	}

	if ev.Message != "test" {
		t.Error("\nCreateBroadcastEvent Failed]\n    Error: wrong post id")
		return
	}

	t.Log("\nCreateBroadcastEvent Succeeded")
}

func TestBroadcastEvent_ToSQLNative(t *testing.T) {
	rec, err := CreateBroadcastEvent(69420, 6969, "test", "test", 0, time.Now())
	if err != nil {
		t.Error("\nBroadcastEvent_ToSQLNative Failed")
		return
	}

	statement := rec.ToSQLNative()

	if statement.Statement != "insert ignore into broadcast_event(_id, user_id, user_name, message, broadcast_type, time_posted) values(?, ?, ?, ?, ?, ?);" {
		t.Errorf("\nBroadcastEvent_ToSQLNative failed\n    Error: incorrect statement returned")
		return
	}

	if len(statement.Values) != 6 {
		fmt.Println("number of values returned: ", len(statement.Values))
		t.Errorf("\nBroadcastEvent_ToSQLNative failed\n    Error: incorrect values returned %v", statement.Values)
		return
	}

	t.Log("\nBroadcastEvent_ToSQLNative succeeded")
}

func TestBroadcastEventFromSQLNative(t *testing.T) {
	db, err := ti.CreateDatabase("gigo-dev-tidb", "4000", "mysql", "gigo-dev", "gigo-dev", "gigo_dev_test")
	if err != nil {
		t.Error("\nBroadcastEventFromSQLNative failed\n    Error: ", err)
		return
	}

	defer db.DB.Exec("DROP TABLE broadcast_event")

	message := "test"
	ev, err := CreateBroadcastEvent(69420, 6969, "test", "test", 0, time.Now())
	if err != nil {
		t.Error("\nBroadcastEventFromSQLNative Failed")
		return
	}

	statement := ev.ToSQLNative()

	stmt, err := db.DB.Prepare(statement.Statement)
	if err != nil {
		t.Error("\nBroadcastEventFromSQLNative failed\n    Error: ", err)
		return
	}

	_, err = stmt.Exec(statement.Values...)
	if err != nil {
		t.Error("\nBroadcastEventFromSQLNative failed\n    Error: ", err)
		return
	}

	rows, err := db.DB.Query("select * from broadcast_event")
	if err != nil {
		t.Error("\nBroadcastEventFromSQLNative failed\n    Error: ", err)
		return
	}

	if !rows.Next() {
		t.Error("\nBroadcastEventFromSQLNative failed\n    Error: no rows found")
		return
	}

	rec, err := BroadcastEventFromSQLNative(rows)
	if err != nil {
		t.Error("\nBroadcastEventFromSQLNative failed\n    Error: ", err)
		return
	}

	if rec == nil {
		t.Error("\nBroadcastEventFromSQLNative Failed]\n    Error: creation returned nil")
		return
	}

	if rec.ID != 69420 {
		t.Error("\nBroadcastEventFromSQLNative Failed]\n    Error: wrong id")
		return
	}

	if rec.UserID != 6969 {
		t.Error("\nBroadcastEventFromSQLNative Failed]\n    Error: wrong user id")
		return
	}

	if rec.Message != message {
		t.Error("\nBroadcastEventFromSQLNative Failed]\n    Error: wrong message")
		return
	}

}
