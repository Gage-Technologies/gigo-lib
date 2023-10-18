package models

import (
	"fmt"
	ti "github.com/gage-technologies/gigo-lib/db"
	"testing"
	"time"
)

func TestCreateThreadComment(t *testing.T) {
	rec, err := CreateThreadComment(69420, "test", "author", 42069, time.Now(), TierType(3), 6969, 20, false, 0, 2)
	if err != nil {
		t.Error("\nTestCreateThreadComment Failed")
		return
	}

	if rec == nil {
		t.Error("\nTestCreateThreadComment Failed]\n    Error: creation returned nil")
		return
	}

	if rec.ID != 69420 {
		t.Error("\nTestCreateThreadComment Failed]\n    Error: wrong id")
		return
	}

	if rec.AuthorID != 42069 {
		t.Error("\nTestCreateThreadComment Failed]\n    Error: wrong author id")
		return
	}

	if rec.AuthorTier != 3 {
		t.Error("\nTestCreateThreadComment Failed]\n    Error: wrong post id")
		return
	}

	t.Log("\nTestCreateThreadComment Succeeded")
}

func TestThreadComment_ToSQLNative(t *testing.T) {
	rec, err := CreateThreadComment(69420, "test", "author", 42069, time.Now(), TierType(3), 6969, 20, false, 0, 2)
	if err != nil {
		t.Error("\nTestThreadComment_ToSQLNative Failed")
		return
	}
	statement := rec.ToSQLNative()

	if statement[0].Statement != "insert ignore into thread_comment(_id, body, author, author_id, created_at, author_tier, coffee, comment_id, leads, revision, discussion_level) values(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);" {
		t.Errorf("\nTestThreadComment_ToSQLNative failed\n    Error: incorrect statement returned")
		return
	}

	if len(statement[0].Values) != 11 {
		fmt.Println("number of values returned: ", len(statement[0].Values))
		t.Errorf("\nTestThreadComment_ToSQLNative failed\n    Error: incorrect values returned %v", statement[0].Values)
		return
	}

	t.Log("\nTestThreadComment_ToSQLNative succeeded")
}

func TestThreadCommentFromSQLNative(t *testing.T) {
	db, err := ti.CreateDatabase("gigo-dev-tidb", "4000", "mysql", "gigo-dev", "gigo-dev", "gigo_dev_test")
	if err != nil {
		t.Error("\nTestThreadCommentFromSQLNative failed\n    Error: ", err)
		return
	}

	defer db.DB.Exec("DROP TABLE thread_comment")

	comment, err := CreateThreadComment(69420, "test", "author", 42069, time.Now(), TierType(3), 6969, 20, false, 0, 2)
	if err != nil {
		t.Error("\nTestThreadCommentFromSQLNative Failed")
		return
	}

	statement := comment.ToSQLNative()

	stmt, err := db.DB.Prepare(statement[0].Statement)
	if err != nil {
		t.Error("\nTestThreadCommentFromSQLNative failed\n    Error: ", err)
		return
	}

	fmt.Println("statement is: ", statement[0])

	_, err = stmt.Exec(statement[0].Values...)
	if err != nil {
		t.Error("\nTestThreadCommentFromSQLNative failed\n    Error: ", err)
		return
	}

	rows, err := db.DB.Query("select * from thread_comment")
	if err != nil {
		t.Error("\nTestThreadCommentFromSQLNative failed\n    Error: ", err)
		return
	}

	if !rows.Next() {
		t.Error("\nTestThreadCommentFromSQLNative failed\n    Error: no rows found")
		return
	}

	rec, err := ThreadCommentFromSQLNative(rows)
	if err != nil {
		t.Error("\nTestThreadCommentFromSQLNative failed\n    Error: ", err)
		return
	}

	if rec == nil {
		t.Error("\nTestThreadCommentFromSQLNative Failed]\n    Error: creation returned nil")
		return
	}

	if rec.ID != 69420 {
		t.Error("\nTestThreadCommentFromSQLNative Failed]\n    Error: wrong id")
		return
	}

	if rec.AuthorID != 42069 {
		t.Error("\nTestThreadCommentFromSQLNative Failed]\n    Error: wrong author id")
		return
	}

	if rec.Coffee != 6969 {
		t.Error("\nTestThreadCommentFromSQLNative Failed]\n    Error: wrong post id")
		return
	}
}
