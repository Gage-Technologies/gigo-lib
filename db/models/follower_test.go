package models

import (
	"fmt"
	ti "github.com/gage-technologies/gigo-lib/db"
	"testing"
)

func TestCreateFollower(t *testing.T) {
	rec, err := CreateFollower(69420, 6969)
	if err != nil {
		t.Error("\nCreate Recommended Post Failed")
		return
	}

	if rec == nil {
		t.Error("\nCreate Recommended Post Failed]\n    Error: creation returned nil")
		return
	}

	if rec.Follower != 69420 {
		t.Error("\nCreate Recommended Post Failed]\n    Error: wrong id")
		return
	}

	if rec.Following != 6969 {
		t.Error("\nCreate Recommended Post Failed]\n    Error: wrong author id")
		return
	}

	t.Log("\nCreate Recommended Post Succeeded")
}

func TestFollower_ToSQLNative(t *testing.T) {
	rec, err := CreateFollower(69420, 6969)
	if err != nil {
		t.Error("\nCreate Recommended Post Failed")
		return
	}

	statement := rec.ToSQLNative()

	if statement.Statement != "insert ignore into follower(follower, following) values(?, ?);" {
		t.Errorf("\nRec post to sql native failed\n    Error: incorrect statement returned")
		return
	}

	if len(statement.Values) != 2 {
		fmt.Println("number of values returned: ", len(statement.Values))
		t.Errorf("\nRec post to sql native failed\n    Error: incorrect values returned %v", statement.Values)
		return
	}

	t.Log("\nRec post to sql native succeeded")
}

func TestFollowerFromSQLNative(t *testing.T) {
	db, err := ti.CreateDatabase("gigo-dev-tidb", "4000", "mysql", "gigo-dev", "gigo-dev", "gigo_dev_test")
	if err != nil {
		t.Error("\nInitialize recommended post table sql failed\n    Error: ", err)
		return
	}

	defer db.DB.Exec("DROP TABLE follower")

	post, err := CreateFollower(69420, 6969)
	if err != nil {
		t.Error("\nCreate Recommended Post Failed")
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

	rows, err := db.DB.Query("select * from follower")
	if err != nil {
		t.Error("\nRec post from sql native failed\n    Error: ", err)
		return
	}

	if !rows.Next() {
		t.Error("\nRec post from sql native failed\n    Error: no rows found")
		return
	}

	rec, err := FollowerFromSQLNative(db, rows)
	if err != nil {
		t.Error("\nRec post from sql native failed\n    Error: ", err)
		return
	}

	if rec == nil {
		t.Error("\nCreate Recommended Post Failed]\n    Error: creation returned nil")
		return
	}

	if rec.Follower != 69420 {
		t.Error("\nCreate Recommended Post Failed]\n    Error: wrong id")
		return
	}

	if rec.Following != 6969 {
		t.Error("\nCreate Recommended Post Failed]\n    Error: wrong author id")
		return
	}

}
