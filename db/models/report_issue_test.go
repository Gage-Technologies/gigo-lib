package models

import (
	"fmt"
	ti "github.com/gage-technologies/gigo-lib/db"
	"testing"
)

func TestCreateReportIssue(t *testing.T) {
	rec, err := CreateReportIssue(69420, "about", "page doesn't load", 42069)
	if err != nil {
		t.Error("\nReport Issue Failed")
		return
	}

	if rec == nil {
		t.Error("\nReport Issue Failed]\n    Error: creation returned nil")
		return
	}

	if rec.UserId != 69420 {
		t.Error("\nReport Issue Failed]\n    Error: wrong id")
		return
	}

	if rec.Id != 42069 {
		t.Error("\nReport Issue Failed]\n    Error: wrong author id")
		return
	}

	t.Log("\nReport Issue Succeeded")
}

func TestReportIssue_ToSQLNative(t *testing.T) {
	rec, err := CreateReportIssue(69420, "about", "page doesn't load", 42069)
	if err != nil {
		t.Error("\nReport Issue Failed")
		return
	}

	statement := rec.ToSQLNative()

	if statement.Statement != "insert into report_issue(_id, user_id, date, page, issue) values(?, ?, ?, ?, ?);" {
		t.Errorf("\nReport Issue to sql native failed\n    Error: incorrect statement returned")
		return
	}

	if len(statement.Values) != 5 {
		fmt.Println("number of values returned: ", len(statement.Values))
		t.Errorf("\nReport Issue to sql native failed\n    Error: incorrect values returned %v", statement.Values)
		return
	}

	t.Log("\nReport Issue to sql native succeeded")
}

func TestReportIssueFromSQLNative(t *testing.T) {
	db, err := ti.CreateDatabase("gigo-dev-tidb", "4000", "mysql", "gigo-dev", "gigo-dev", "gigo_dev_test")
	if err != nil {
		t.Error("\nInitialize Report Issue table sql failed\n    Error: ", err)
		return
	}

	defer db.DB.Exec("delete from report_issue")

	post, err := CreateReportIssue(69420, "about", "page doesn't load", 42069)
	if err != nil {
		t.Error("\nReport Issue Failed")
		return
	}

	statement := post.ToSQLNative()

	stmt, err := db.DB.Prepare(statement.Statement)
	if err != nil {
		t.Error("\nReport Issue from sql native failed\n    Error: ", err)
		return
	}

	_, err = stmt.Exec(statement.Values...)
	if err != nil {
		t.Error("\nReport Issue from sql native failed\n    Error: ", err)
		return
	}

	rows, err := db.DB.Query("select * from report_issue")
	if err != nil {
		t.Error("\nReport Issue from sql native failed\n    Error: ", err)
		return
	}

	if !rows.Next() {
		t.Error("\nReport Issue from sql native failed\n    Error: no rows found")
		return
	}

	rec, err := ReportIssueFromSQLNative(db, rows)
	if err != nil {
		t.Error("\nReport Issue from sql native failed\n    Error: ", err)
		return
	}

	if rec == nil {
		t.Error("\nReport Issue Failed]\n    Error: creation returned nil")
		return
	}

	if rec.UserId != 69420 {
		t.Error("\nReport Issue Failed]\n    Error: wrong id")
		return
	}

	if rec.Id != 42069 {
		t.Error("\nReport Issue Failed]\n    Error: wrong author id")
		return
	}

}
