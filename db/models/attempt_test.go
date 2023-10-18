package models

import (
	"fmt"
	ti "github.com/gage-technologies/gigo-lib/db"
	"reflect"
	"testing"
	"time"
)

func TestCreateAttempt(t *testing.T) {
	parent := int64(12345)
	attempt, err := CreateAttempt(69420, "test", "test", "author", 42069, time.Now(), time.Now(), 69, 2, []int64{}, uint64(3), 6969, 20, &parent, 0)
	if err != nil {
		t.Error("\nCreate Attempt Post Failed")
		return
	}

	if attempt == nil {
		t.Error("\nCreate Attempt Post Failed\n    Error: creation returned nil")
		return
	}

	if attempt.ID != 69420 {
		t.Error("\nCreate Attempt Post Failed\n    Error: wrong id")
		return
	}

	if attempt.AuthorID != 42069 {
		t.Error("\nCreate Attempt Post Failed\n    Error: wrong author id")
		return
	}

	if *attempt.ParentAttempt != 12345 {
		t.Error("\nCreate Attempt Post Failed\n    Error: wrong parent attempt")
		return
	}

	t.Log("\nCreate Attempt Post Succeeded")
}

func TestAttempt_ToSQLNative(t *testing.T) {
	parent := int64(12345)
	attempt, err := CreateAttempt(69420, "test", "test", "author", 42069, time.Now(), time.Now(), 69, 4, []int64{}, uint64(3), 6969, 20, &parent, 0)
	if err != nil {
		t.Error("\nCreate Attempt Post Failed")
		return
	}

	statement, err := attempt.ToSQLNative()
	if err != nil {
		t.Error("\nCreate Attempt ToSQLNative Failed")
		return
	}

	if statement[0].Statement != "insert ignore into attempt(_id, post_title, description, author, author_id, created_at, updated_at, repo_id, "+
		"author_tier, coffee, post_id, closed, success, closed_date, tier, parent_attempt, workspace_settings, post_type) values(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);" {
		t.Errorf("\nattempt to sql native failed\n    Error: incorrect statement returned")
		return
	}

	if len(statement[0].Values) != 18 {
		fmt.Println("number of values returned: ", len(statement[0].Values))
		t.Errorf("\nattempt to sql native failed\n    Error: incorrect values returned %v", statement[0].Values)
		return
	}

	t.Log("\nattempt to sql native succeeded")
}

func TestAttemptFromSQLNative(t *testing.T) {
	db, err := ti.CreateDatabase("gigo-dev-tidb", "4000", "mysql", "gigo-dev", "gigo-dev", "gigo_dev_test")
	if err != nil {
		t.Fatalf("\n%s failed\n    Error: %v", t.Name(), err)
	}

	parent := int64(12345)
	post, err := CreateAttempt(69420, "test", "test", "author", 42069, time.Now(), time.Now(), 69, 5, []int64{}, uint64(3), 6969, 3, &parent, 0)
	if err != nil {
		t.Fatalf("\n%s failed\n    Error: %v", t.Name(), err)
	}

	post.WorkspaceSettings = &DefaultWorkspaceSettings

	statement, err := post.ToSQLNative()
	if err != nil {
		t.Fatalf("\n%s failed\n    Error: %v", t.Name(), err)
	}

	stmt, err := db.DB.Prepare(statement[0].Statement)
	if err != nil {
		t.Fatalf("\n%s failed\n    Error: %v", t.Name(), err)
	}

	_, err = stmt.Exec(statement[0].Values...)
	if err != nil {
		t.Fatalf("\n%s failed\n    Error: %v", t.Name(), err)
	}

	rows, err := db.DB.Query("select * from attempt")
	if err != nil {
		t.Fatalf("\n%s failed\n    Error: %v", t.Name(), err)
	}

	if !rows.Next() {
		t.Fatalf("\n%s failed\n    Error: no rows found", t.Name())
		return
	}

	rec, err := AttemptFromSQLNative(db, rows)
	if err != nil {
		t.Fatalf("\n%s failed\n    Error: %v", t.Name(), err)
	}

	if rec == nil {
		t.Fatalf("\n%s failed\n    Error: nil return", t.Name())
	}

	if rec.ID != 69420 {
		t.Fatalf("\n%s failed\n    Error: wrong id", t.Name())
	}

	if rec.AuthorID != 42069 {
		t.Fatalf("\n%s failed\n    Error: wrong author id", t.Name())
	}

	if rec.PostID != 6969 {
		t.Fatalf("\n%s failed\n    Error: wrong post id", t.Name())
	}

	if rec.Tier != 3 {
		t.Fatalf("\n%s failed\n    Error: wrong tier", t.Name())
	}

	if *rec.ParentAttempt != 12345 {
		t.Fatalf("\n%s failed\n    Error: wrong parent attempt id", t.Name())
	}

	if !reflect.DeepEqual(*rec.WorkspaceSettings, DefaultWorkspaceSettings) {
		t.Fatalf("\n%s failed\n    Error: wrong workspace settings", t.Name())
	}

	t.Logf("\n%s succeeded", t.Name())
}
