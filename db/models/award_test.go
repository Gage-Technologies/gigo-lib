package models

import (
	"fmt"
	ti "github.com/gage-technologies/gigo-lib/db"
	"testing"
)

func TestCreateAward(t *testing.T) {
	rec, err := CreateAward(69420, "test", ContentType(2))
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

	if rec.Award != "test" {
		t.Error("\nCreate Recommended Post Failed]\n    Error: wrong post id")
		return
	}

	t.Log("\nCreate Recommended Post Succeeded")
}

func TestAward_ToSQLNative(t *testing.T) {
	rec, err := CreateAward(69420, "test", ContentType(2))
	if err != nil {
		t.Error("\nCreate Repo Post Failed")
		return
	}
	statement := rec.ToSQLNative()

	if statement.Statement != "insert ignore into award(_id, award, types) values(?, ?, ?);" {
		t.Errorf("\nRec post to sql native failed\n    Error: incorrect statement returned")
		return
	}

	if len(statement.Values) != 3 {
		fmt.Println("number of values returned: ", len(statement.Values))
		t.Errorf("\nRec post to sql native failed\n    Error: incorrect values returned %v", statement.Values)
		return
	}

	t.Log("\nRec post to sql native succeeded")
}

func TestAwardFromSQLNative(t *testing.T) {
	db, err := ti.CreateDatabase("gigo-dev-tidb", "4000", "mysql", "gigo-dev", "gigo-dev", "gigo_dev_test")
	if err != nil {
		t.Error("\nInitialize recommended post table sql failed\n    Error: ", err)
		return
	}

	comment, err := CreateAward(69420, "test", ContentType(2))
	if err != nil {
		t.Error("\nCreate Repo Post Failed")
		return
	}

	statement := comment.ToSQLNative()

	stmt, err := db.DB.Prepare(statement.Statement)
	if err != nil {
		t.Error("\nRec post from sql native failed\n    Error: ", err)
		return
	}

	fmt.Println("statement is: ", statement)

	_, err = stmt.Exec(statement.Values...)
	if err != nil {
		t.Error("\nRec post from sql native failed\n    Error: ", err)
		return
	}

	rows, err := db.DB.Query("select * from award")
	if err != nil {
		t.Error("\nRec post from sql native failed\n    Error: ", err)
		return
	}

	if !rows.Next() {
		t.Error("\nRec post from sql native failed\n    Error: no rows found")
		return
	}

	rec, err := AwardFromSQLNative(rows)
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

	if rec.Award != "test" {
		t.Error("\nCreate Recommended Post Failed]\n    Error: wrong author id")
		return
	}

}
