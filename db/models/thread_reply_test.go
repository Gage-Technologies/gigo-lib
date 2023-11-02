package models

import (
	"fmt"
	ti "github.com/gage-technologies/gigo-lib/db"
	"testing"
	"time"
)

func TestCreateThreadReply(t *testing.T) {
	rec, err := CreateThreadReply(69420, "test", "author", 42069, time.Now(), TierType(3), 6969, 20, 0, 3)
	if err != nil {
		t.Error("\nCreateThreadReply Failed")
		return
	}

	if rec == nil {
		t.Error("\nCreateThreadReply Failed]\n    Error: creation returned nil")
		return
	}

	if rec.ID != 69420 {
		t.Error("\nCreateThreadReply Failed]\n    Error: wrong id")
		return
	}

	if rec.AuthorID != 42069 {
		t.Error("\nCreateThreadReply Failed]\n    Error: wrong author id")
		return
	}

	if rec.AuthorTier != 3 {
		t.Error("\nCreateThreadReply Failed]\n    Error: wrong post id")
		return
	}

	t.Log("\nCreateThreadReply Succeeded")
}

func TestInitializeThreadReplyTableSQL(t *testing.T) {
	db, err := ti.CreateDatabase("gigo-dev-tidb", "4000", "mysql", "gigo-dev", "gigo-dev", "gigo_dev_test")
	if err != nil {
		t.Error("\nInitializeThreadReplyTable sql failed\n    Error: ", err)
		return
	}

	defer db.DB.Exec("delete from thread_reply")

	res, err := db.DB.Query("SELECT * FROM information_schema.tables WHERE table_schema = 'gigo_dev_test' AND table_name = 'thread_reply' LIMIT 1;")
	if err != nil {
		t.Error("\nInitializeThreadReplyTable sql failed\n    Error: ", err)
		return
	}

	if !res.Next() {
		t.Error("\nInitializeThreadReplyTable sql failed\n    Error: table was not created")
		return
	}

	t.Log("\nInitializeThreadReplyTable sql succeeded")
}

func TestThreadReply_ToSQLNative(t *testing.T) {
	rec, err := CreateThreadReply(69420, "test", "author", 42069, time.Now(), TierType(3), 6969, 20, 0, 3)
	if err != nil {
		t.Error("\nThreadReply_ToSQLNative Failed")
		return
	}
	statement := rec.ToSQLNative()

	if statement[0].Statement != "insert ignore into thread_reply(_id, body, author, author_id, created_at, author_tier, coffee, thread_comment_id, revision, discussion_level) values(?, ?, ?, ?, ?, ?, ?, ?, ?, ?);" {
		t.Errorf("\nThreadReply_ToSQLNative failed\n    Error: incorrect statement returned")
		return
	}

	if len(statement[0].Values) != 10 {
		fmt.Println("number of values returned: ", len(statement[0].Values))
		t.Errorf("\nThreadReply_ToSQLNative failed\n    Error: incorrect values returned %v", statement[0].Values)
		return
	}

	t.Log("\nThreadReply_ToSQLNative succeeded")
}

func TestThreadReplyFromSQLNative(t *testing.T) {
	db, err := ti.CreateDatabase("gigo-dev-tidb", "4000", "mysql", "gigo-dev", "gigo-dev", "gigo_dev_test")
	if err != nil {
		t.Error("\nThreadReplyFromSQLNative failed\n    Error: ", err)
		return
	}

	defer db.DB.Exec("delete from thread_reply")

	comment, err := CreateThreadReply(69420, "test", "author", 42069, time.Now(), TierType(3), 6969, 20, 0, 3)
	if err != nil {
		t.Error("\nThreadReplyFromSQLNative Failed")
		return
	}

	statement := comment.ToSQLNative()

	stmt, err := db.DB.Prepare(statement[0].Statement)
	if err != nil {
		t.Error("\nThreadReplyFromSQLNative failed\n    Error: ", err)
		return
	}

	fmt.Println("statement is: ", statement[0])

	_, err = stmt.Exec(statement[0].Values...)
	if err != nil {
		t.Error("\nThreadReplyFromSQLNative failed\n    Error: ", err)
		return
	}

	rows, err := db.DB.Query("select * from thread_reply")
	if err != nil {
		t.Error("\nThreadReplyFromSQLNative failed\n    Error: ", err)
		return
	}

	if !rows.Next() {
		t.Error("\nThreadReplyFromSQLNative failed\n    Error: no rows found")
		return
	}

	rec, err := ThreadReplyFromSQLNative(rows)
	if err != nil {
		t.Error("\nThreadReplyFromSQLNative failed\n    Error: ", err)
		return
	}

	if rec == nil {
		t.Error("\nThreadReplyFromSQLNative Failed]\n    Error: creation returned nil")
		return
	}

	if rec.ID != 69420 {
		t.Error("\nThreadReplyFromSQLNative Failed]\n    Error: wrong id")
		return
	}

	if rec.AuthorID != 42069 {
		t.Error("\nThreadReplyFromSQLNative Failed]\n    Error: wrong author id")
		return
	}

	if rec.Coffee != 6969 {
		t.Error("\nThreadReplyFromSQLNative Failed]\n    Error: wrong post id")
		return
	}
}
