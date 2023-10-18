package models

import (
	"fmt"
	ti "github.com/gage-technologies/gigo-lib/db"
	"testing"
	"time"
)

func TestCreateComment(t *testing.T) {
	rec, err := CreateComment(69420, "test", "author", 42069, time.Now(), TierType(3), []int64{}, 6969, 20, true, 0, 1)
	if err != nil {
		t.Error("\nTestCreateComment Failed")
		return
	}

	if rec == nil {
		t.Error("\nTestCreateComment Failed]\n    Error: creation returned nil")
		return
	}

	if rec.ID != 69420 {
		t.Error("\nTestCreateComment Failed]\n    Error: wrong id")
		return
	}

	if rec.AuthorID != 42069 {
		t.Error("\nTestCreateComment Failed]\n    Error: wrong author id")
		return
	}

	if rec.AuthorTier != 3 {
		t.Error("\nTestCreateComment Failed]\n    Error: wrong post id")
		return
	}

	t.Log("\nTestCreateComment Succeeded")
}

func TestComment_ToSQLNative(t *testing.T) {
	rec, err := CreateComment(69420, "test", "author", 42069, time.Now(), TierType(3), []int64{}, 6969, 20, true, 0, 1)
	if err != nil {
		t.Error("\nTestComment_ToSQLNative Failed")
		return
	}
	statement := rec.ToSQLNative()

	if statement[0].Statement != "insert ignore into comment(_id, body, author, author_id, created_at, author_tier, coffee, discussion_id, leads, revision, discussion_level) values(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);" {
		t.Errorf("\nTestComment_ToSQLNative failed\n    Error: incorrect statement returned")
		return
	}

	if len(statement[0].Values) != 11 {
		fmt.Println("number of values returned: ", len(statement[0].Values))
		t.Errorf("\nTestComment_ToSQLNative\n    Error: incorrect values returned %v", statement[0].Values)
		return
	}

	t.Log("\nTestComment_ToSQLNative succeeded")
}

func TestCommentFromSQLNative(t *testing.T) {
	db, err := ti.CreateDatabase("gigo-dev-tidb", "4000", "mysql", "gigo-dev", "gigo-dev", "gigo_dev_test")
	if err != nil {
		t.Error("\nTestCommentFromSQLNative failed\n    Error: ", err)
		return
	}

	defer db.DB.Exec("DROP TABLE comment")
	defer db.DB.Exec("DROP TABLE comment_awards")

	comment, err := CreateComment(69420, "test", "author", 42069, time.Now(), TierType(3), []int64{5}, 6969, 20, true, 0, 1)
	if err != nil {
		t.Error("\nTestCommentFromSQLNative Failed")
		return
	}

	statements := comment.ToSQLNative()

	for _, statement := range statements {
		stmt, err := db.DB.Prepare(statement.Statement)
		if err != nil {
			t.Fatalf("%s failed\n    Error: %v", t.Name(), err)
		}

		_, err = stmt.Exec(statement.Values...)
		if err != nil {
			t.Fatalf("%s failed\n    Error: %v", t.Name(), err)
		}
	}

	rows, err := db.DB.Query("select * from comment")
	if err != nil {
		t.Error("\nTestCommentFromSQLNative failed\n    Error: ", err)
		return
	}

	if !rows.Next() {
		t.Error("\nTestCommentFromSQLNative failed\n    Error: no rows found")
		return
	}

	rec, err := CommentFromSQLNative(db, rows)
	if err != nil {
		t.Error("\nTestCommentFromSQLNative failed\n    Error: ", err)
		return
	}

	if rec == nil {
		t.Error("\nTestCommentFromSQLNative Failed]\n    Error: creation returned nil")
		return
	}

	if rec.ID != 69420 {
		t.Error("\nTestCommentFromSQLNative Failed]\n    Error: wrong id")
		return
	}

	if rec.AuthorID != 42069 {
		t.Error("\nTestCommentFromSQLNative Failed]\n    Error: wrong author id")
		return
	}

	if rec.Coffee != 6969 {
		t.Error("\nTestCommentFromSQLNative Failed]\n    Error: wrong post id")
		return
	}

}
