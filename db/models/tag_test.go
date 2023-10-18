package models

import (
	"fmt"
	ti "github.com/gage-technologies/gigo-lib/db"
	"testing"
)

func TestCreateTag(t *testing.T) {
	tag := CreateTag(42069, "test")

	if tag.ID != 42069 {
		t.Errorf("\nCreateTag failed\n    Error: incorrect id returned")
		return
	}

	if tag.Value != "test" {
		t.Errorf("\nCreateTag failed\n    Error: incorrect value returned")
		return
	}

	if tag.Official != false {
		t.Errorf("\nCreateTag failed\n    Error: incorrect official returned")
		return
	}

	if tag.UsageCount != 1 {
		t.Errorf("\nCreateTag failed\n    Error: incorrect usage count returned")
		return
	}

	t.Logf("\nCreateTag succeeded")
}

func TestTag_ToSQLNative(t *testing.T) {
	tag := CreateTag(42069, "test")

	statement := tag.ToSQLNative()

	if statement[0].Statement != "insert ignore into tag(_id, value, official, usage_count) values (?, ?, ?, ?)" {
		t.Errorf("\nTag ToSQLNative failed\n    Error: incorrect statement returned")
		return
	}

	if len(statement[0].Values) != 4 {
		fmt.Println("number of values returned: ", len(statement[0].Values))
		t.Errorf("\nTag ToSQLNative failed\n    Error: incorrect values returned %v", statement[0].Values)
		return
	}

	t.Log("\nTag ToSQLNative succeeded")
}

func TestTagFromSQLNative(t *testing.T) {
	db, err := ti.CreateDatabase("gigo-dev-tidb", "4000", "mysql", "gigo-dev", "gigo-dev", "gigo_dev_test")
	if err != nil {
		t.Error("\nTagFromSQLNative failed\n    Error: ", err)
		return
	}

	defer db.DB.Exec("DROP TABLE post")

	tag := CreateTag(42069, "test")
	tag.Official = true
	tag.UsageCount = 420

	statement := tag.ToSQLNative()

	fmt.Println("statement: ", statement[0].Statement)

	stmt, err := db.DB.Prepare(statement[0].Statement)
	if err != nil {
		t.Error("\nTagFromSQLNative failed\n    Error: ", err)
		return
	}

	_, err = stmt.Exec(statement[0].Values...)
	if err != nil {
		t.Error("\nTagFromSQLNative failed\n    Error: ", err)
		return
	}

	rows, err := db.DB.Query("select * from tag where _id = 42069")
	if err != nil {
		t.Error("\nTagFromSQLNative failed\n    Error: ", err)
		return
	}

	if !rows.Next() {
		t.Error("\nTagFromSQLNative failed\n    Error: no rows found")
		return
	}

	tag2, err := TagFromSQLNative(rows)
	if err != nil {
		t.Error("\nTagFromSQLNative failed\n    Error: ", err)
		return
	}

	if tag.ID != tag2.ID {
		t.Errorf("\nTagFromSQLNative failed\n    Error: incorrect id returned")
		return
	}

	if tag.Value != tag2.Value {
		t.Errorf("\nTagFromSQLNative failed\n    Error: incorrect value returned")
		return
	}

	if tag.Official != tag2.Official {
		t.Errorf("\nTagFromSQLNative failed\n    Error: incorrect official returned")
		return
	}

	if tag.UsageCount != tag.UsageCount {
		t.Errorf("\nTagFromSQLNative failed\n    Error: incorrect usage count returned")
		return
	}

	t.Logf("\nTagFromSQLNative succeeded")
}
