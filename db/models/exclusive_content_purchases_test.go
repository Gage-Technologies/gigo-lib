package models

import (
	"fmt"
	ti "github.com/gage-technologies/gigo-lib/db"
	"testing"
)

func TestCreateExclusiveContentPurchases(t *testing.T) {
	rec, err := CreateExclusiveContentPurchases(69420, 6969)
	if err != nil {
		t.Error("\nCreate exclusive content Failed")
		return
	}

	if rec == nil {
		t.Error("\nCreate exclusive content Failed]\n    Error: creation returned nil")
		return
	}

	if rec.UserId != 69420 {
		t.Error("\nCreate exclusive content Failed]\n    Error: wrong id")
		return
	}

	if rec.Post != 6969 {
		t.Error("\nCreate exclusive content Failed]\n    Error: wrong author id")
		return
	}

	t.Log("\nCreate exclusive content Succeeded")
}

func TestExclusiveContentPurchases_ToSQLNative(t *testing.T) {
	rec, err := CreateExclusiveContentPurchases(69420, 6969)
	if err != nil {
		t.Error("\nCreate exclusive content Failed")
		return
	}

	statement := rec.ToSQLNative()

	if statement.Statement != "insert ignore into exclusive_content_purchases(user_id, post, date) values(?, ?, ?);" {
		t.Errorf("\nexclusive content to sql native failed\n    Error: incorrect statement returned")
		return
	}

	if len(statement.Values) != 3 {
		fmt.Println("number of values returned: ", len(statement.Values))
		t.Errorf("\nexclusive content to sql native failed\n    Error: incorrect values returned %v", statement.Values)
		return
	}

	t.Log("\nexclusive content to sql native succeeded")
}

func TestExclusiveContentPurchasesFromSQLNative(t *testing.T) {
	db, err := ti.CreateDatabase("gigo-dev-tidb", "4000", "mysql", "gigo-dev", "gigo-dev", "gigo_dev_test")
	if err != nil {
		t.Error("\nInitialize exclusive content table sql failed\n    Error: ", err)
		return
	}

	defer db.DB.Exec("delete from exclusive_content_purchases")

	post, err := CreateExclusiveContentPurchases(69420, 6969)
	if err != nil {
		t.Error("\nCreate exclusive content Failed")
		return
	}

	statement := post.ToSQLNative()

	stmt, err := db.DB.Prepare(statement.Statement)
	if err != nil {
		t.Error("\nexclusive content from sql native failed\n    Error: ", err)
		return
	}

	_, err = stmt.Exec(statement.Values...)
	if err != nil {
		t.Error("\nexclusive content from sql native failed\n    Error: ", err)
		return
	}

	rows, err := db.DB.Query("select * from exclusive_content_purchases")
	if err != nil {
		t.Error("\nexclusive content from sql native failed\n    Error: ", err)
		return
	}

	if !rows.Next() {
		t.Error("\nexclusive content from sql native failed\n    Error: no rows found")
		return
	}

	rec, err := ExclusiveContentPurchasesFromSQLNative(db, rows)
	if err != nil {
		t.Error("\nexclusive content from sql native failed\n    Error: ", err)
		return
	}

	if rec == nil {
		t.Error("\nCreate exclusive content Failed]\n    Error: creation returned nil")
		return
	}

	if rec.UserId != 69420 {
		t.Error("\nCreate exclusive content Failed]\n    Error: wrong id")
		return
	}

	if rec.Post != 6969 {
		t.Error("\nCreate exclusive content Failed]\n    Error: wrong author id")
		return
	}

}
