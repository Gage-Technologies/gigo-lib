package models

import (
	"fmt"
	ti "github.com/gage-technologies/gigo-lib/db"
	"testing"
	"time"
)

func TestCreateDiscussion(t *testing.T) {
	rec, err := CreateDiscussion(69420, "test", "author", 42069, time.Now(), time.Now(), 3, []int64{420}, 6969, 20, "test", []int64{69}, false, 0, 0)
	if err != nil {
		t.Error("\nTestCreateDiscussion Failed")
		return
	}

	if rec == nil {
		t.Error("\nTestCreateDiscussion Failed]\n    Error: creation returned nil")
		return
	}

	if rec.ID != 69420 {
		t.Error("\nTestCreateDiscussion Failed]\n    Error: wrong id")
		return
	}

	if rec.AuthorID != 42069 {
		t.Error("\nTestCreateDiscussion Failed]\n    Error: wrong author id")
		return
	}

	if rec.AuthorTier != 3 {
		t.Error("\nTestCreateDiscussion Failed]\n    Error: wrong post id")
		return
	}

	t.Log("\nTestCreateDiscussion Succeeded")
}

func TestDiscussion_ToSQLNative(t *testing.T) {
	rec, err := CreateDiscussion(69420, "test", "author", 42069, time.Now(), time.Now(), 4, []int64{}, 6969, 20, "test", []int64{}, false, 0, 0)
	if err != nil {
		t.Error("\nDiscussion_ToSQLNative Failed")
		return
	}
	statement := rec.ToSQLNative()

	if statement[0].Statement != "insert ignore into discussion(_id, body, author, author_id, created_at, updated_at, author_tier, coffee, post_id, title, leads, revision, discussion_level) values(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);" {
		t.Errorf("\nDiscussion_ToSQLNative failed\n    Error: incorrect statement returned")
		return
	}

	if len(statement[0].Values) != 13 {
		fmt.Println("number of values returned: ", len(statement[0].Values))
		t.Errorf("\nDiscussion_ToSQLNative failed\n    Error: incorrect values returned %v", statement[0].Values)
		return
	}

	t.Log("\nDiscussion_ToSQLNative succeeded")
}

func TestDiscussionFromSQLNative(t *testing.T) {
	db, err := ti.CreateDatabase("gigo-dev-tidb", "4000", "mysql", "gigo-dev", "gigo-dev", "gigo_dev_test")
	if err != nil {
		t.Error("\nDiscussionFromSQLNative failed\n    Error: ", err)
		return
	}

	defer db.DB.Exec("DROP TABLE discussion")
	defer db.DB.Exec("DROP TABLE discussion_awards")
	defer db.DB.Exec("DROP TABLE discussion_tags")

	d, err := CreateDiscussion(69420, "test", "author", 42069, time.Now(), time.Now(), 5, []int64{69}, 6969, 20, "test", []int64{420}, false, 0, 0)
	if err != nil {
		t.Error("\nDiscussionFromSQLNative Failed")
		return
	}

	statements := d.ToSQLNative()

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

	rows, err := db.DB.Query("select * from discussion where _id = ?", d.ID)
	if err != nil {
		t.Error("\nDiscussionFromSQLNative failed\n    Error: ", err)
		return
	}

	if !rows.Next() {
		t.Error("\nDiscussionFromSQLNative failed\n    Error: no rows found")
		return
	}

	rec, err := DiscussionFromSQLNative(db, rows)
	if err != nil {
		t.Error("\nDiscussionFromSQLNative failed\n    Error: ", err)
		return
	}

	if rec == nil {
		t.Error("\nDiscussionFromSQLNative Failed]\n    Error: creation returned nil")
		return
	}

	if rec.ID != 69420 {
		t.Error("\nDiscussionFromSQLNative Failed]\n    Error: wrong id")
		return
	}

	if rec.AuthorID != 42069 {
		t.Error("\nDiscussionFromSQLNative Failed]\n    Error: wrong author id")
		return
	}

	if rec.Coffee != 6969 {
		t.Error("\nDiscussionFromSQLNative Failed]\n    Error: wrong post id")
		return
	}

	t.Log("\nDiscussionFromSQLNative succeeded")
}
