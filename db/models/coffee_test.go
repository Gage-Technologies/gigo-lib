package models

import (
	"fmt"
	ti "github.com/gage-technologies/gigo-lib/db"
	"testing"
	"time"
)

func TestCreateCoffee(t *testing.T) {
	attemptId := int64(323)
	postId := int64(455)
	discussionId := int64(9595)
	rec, err := CreateCoffee(69420, 6969, time.Now(), time.Now(), &attemptId, &postId, &discussionId)
	if err != nil {
		t.Error("\nCreate Repo Post Failed")
		return
	}

	if rec == nil {
		t.Error("\nCreate Recommended Post Failed]\n    Error: creation returned nil")
		return
	}

	if rec.ID != 69420 {
		t.Error("\nCreate Recommended Post Failed]\n    Error: wrong id")
		return
	}

	if rec.UserID != 6969 {
		t.Error("\nCreate Recommended Post Failed]\n    Error: wrong author id")
		return
	}

	if *rec.PostID != postId {
		t.Error("\nCreate Recommended Post Failed]\n    Error: wrong post id")
		return
	}

	if *rec.DiscussionID != discussionId {
		t.Error("\nCreate Recommended Post Failed]\n    Error: wrong similarity")
		return
	}

	t.Log("\nCreate Recommended Post Succeeded")
}

func TestCoffee_ToSQLNative(t *testing.T) {
	attemptId := int64(323)
	postId := int64(455)
	discussionId := int64(9595)
	rec, err := CreateCoffee(69420, 6969, time.Now(), time.Now(), &attemptId, &postId, &discussionId)
	if err != nil {
		t.Error("\nCreate Repo Post Failed")
		return
	}

	statement := rec.ToSQLNative()

	if statement.Statement != "insert ignore into coffee(_id, created_at, updated_at, post_id, attempt_id, user_id, discussion_id) values(?, ?, ?, ?, ?, ?, ?);" {
		t.Errorf("\nRec post to sql native failed\n    Error: incorrect statement returned")
		return
	}

	if len(statement.Values) != 7 {
		fmt.Println("number of values returned: ", len(statement.Values))
		t.Errorf("\nRec post to sql native failed\n    Error: incorrect values returned %v", statement.Values)
		return
	}

	t.Log("\nRec post to sql native succeeded")
}

func TestCoffeeFromSQLNative(t *testing.T) {
	db, err := ti.CreateDatabase("gigo-dev-tidb", "4000", "mysql", "gigo-dev", "gigo-dev", "gigo_dev_test")
	if err != nil {
		t.Error("\nInitialize recommended post table sql failed\n    Error: ", err)
		return
	}

	defer db.DB.Exec("delete from coffee")

	attemptId := int64(323)
	postId := int64(455)
	discussionId := int64(9595)
	post, err := CreateCoffee(69420, 6969, time.Now(), time.Now(), &attemptId, &postId, &discussionId)
	if err != nil {
		t.Error("\nCreate Repo Post Failed")
		return
	}

	statement := post.ToSQLNative()

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

	rows, err := db.DB.Query("select * from coffee")
	if err != nil {
		t.Error("\nRec post from sql native failed\n    Error: ", err)
		return
	}

	if !rows.Next() {
		t.Error("\nRec post from sql native failed\n    Error: no rows found")
		return
	}

	rec, err := CoffeeFromSQLNative(rows)
	if err != nil {
		t.Error("\nRec post from sql native failed\n    Error: ", err)
		return
	}

	if rec == nil {
		t.Error("\nCreate Recommended Post Failed]\n    Error: creation returned nil")
		return
	}

	if rec.ID != 69420 {
		t.Error("\nCreate Recommended Post Failed]\n    Error: wrong id")
		return
	}

	if rec.UserID != 6969 {
		t.Error("\nCreate Recommended Post Failed]\n    Error: wrong author id")
		return
	}

	if *rec.PostID != postId {
		t.Error("\nCreate Recommended Post Failed]\n    Error: wrong post id")
		return
	}

	if *rec.DiscussionID != discussionId {
		t.Error("\nCreate Recommended Post Failed]\n    Error: wrong similarity")
		return
	}

}
