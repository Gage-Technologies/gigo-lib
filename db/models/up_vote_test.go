package models

import (
	"fmt"
	ti "github.com/gage-technologies/gigo-lib/db"
	"testing"
)

func TestCreateUpVote(t *testing.T) {
	upVote := CreateUpVote(42069, 0, 69, 420)

	if upVote.ID != 42069 {
		t.Errorf("\nTestCreateUpVote failed\n    Error: incorrect id returned")
		return
	}

	if upVote.DiscussionType != 0 {
		t.Errorf("\nTestCreateUpVote failed\n    Error: incorrect DiscussionType returned")
		return
	}

	if upVote.DiscussionId != 69 {
		t.Errorf("\nTestCreateUpVote failed\n    Error: incorrect DiscussionId returned")
		return
	}

	if upVote.UserId != 420 {
		t.Errorf("\nTestCreateUpVote failed\n    Error: incorrect UserId returned")
		return
	}

	t.Logf("\nTestCreateUpVote succeeded")
}

func TestUpVote_ToSQLNative(t *testing.T) {
	upVote := CreateUpVote(42069, 0, 69, 420)

	statement := upVote.ToSQLNative()

	if statement[0].Statement != "insert ignore into up_vote(_id, discussion_type, discussion_id, user_id) values (?, ?, ?, ?)" {
		t.Errorf("\nTestUpVote_ToSQLNative failed\n    Error: incorrect statement returned")
		return
	}

	if len(statement[0].Values) != 4 {
		fmt.Println("number of values returned: ", len(statement[0].Values))
		t.Errorf("\nTestUpVote_ToSQLNative failed\n    Error: incorrect values returned %v", statement[0].Values)
		return
	}

	t.Log("\nTestUpVote_ToSQLNative succeeded")
}

func TestUpVoteFromSQLNative(t *testing.T) {
	db, err := ti.CreateDatabase("gigo-dev-tidb", "4000", "mysql", "gigo-dev", "gigo-dev", "gigo_dev_test")
	if err != nil {
		t.Error("\nTestUpVoteFromSQLNative failed\n    Error: ", err)
		return
	}

	defer db.DB.Exec("DROP TABLE post")

	upVote := CreateUpVote(42069, 0, 69, 420)

	statement := upVote.ToSQLNative()

	fmt.Println("statement: ", statement[0].Statement)

	stmt, err := db.DB.Prepare(statement[0].Statement)
	if err != nil {
		t.Error("\nTestUpVoteFromSQLNative failed\n    Error: ", err)
		return
	}

	_, err = stmt.Exec(statement[0].Values...)
	if err != nil {
		t.Error("\nTestUpVoteFromSQLNative failed\n    Error: ", err)
		return
	}

	rows, err := db.DB.Query("select * from up_vote where _id = 42069")
	if err != nil {
		t.Error("\nTestUpVoteFromSQLNative failed\n    Error: ", err)
		return
	}

	if !rows.Next() {
		t.Error("\nTestUpVoteFromSQLNative failed\n    Error: no rows found")
		return
	}

	upVote2, err := UpVoteFromSQLNative(rows)
	if err != nil {
		t.Error("\nTestUpVoteFromSQLNative failed\n    Error: ", err)
		return
	}

	if upVote.ID != upVote.ID {
		t.Errorf("\nTestUpVoteFromSQLNative failed\n    Error: incorrect id returned")
		return
	}

	if upVote.DiscussionType != upVote2.DiscussionType {
		t.Errorf("\nTestUpVoteFromSQLNative failed\n    Error: incorrect value returned")
		return
	}

	if upVote.DiscussionId != upVote2.DiscussionId {
		t.Errorf("\nTestUpVoteFromSQLNative failed\n    Error: incorrect official returned")
		return
	}

	if upVote.UserId != upVote2.UserId {
		t.Errorf("\nTestUpVoteFromSQLNative failed\n    Error: incorrect usage count returned")
		return
	}

	t.Logf("\nTestUpVoteFromSQLNative succeeded")
}
