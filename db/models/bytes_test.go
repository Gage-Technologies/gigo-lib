package models

import (
	"fmt"
	ti "github.com/gage-technologies/gigo-lib/db"
	"testing"
)

func TestCreateBytes(t *testing.T) {
	bytes, err := CreateBytes(69420, "test", "this is a test description", "test outline for bytes", "steps 1 2 3 4", 5, "#fff")
	if err != nil {
		t.Error("\nCreate Bytes Post Failed")
		return
	}

	if bytes == nil {
		t.Error("\nCreate Bytes Post Failed\n    Error: creation returned nil")
		return
	}

	if bytes.ID != 69420 {
		t.Error("\nCreate Bytes Post Failed\n    Error: wrong id")
		return
	}

	if bytes.Name != "test" {
		t.Error("\nCreate Bytes Post Failed\n    Error: wrong name")
		return
	}

	if bytes.DescriptionEasy != "this is a test description" {
		t.Error("\nCreate Bytes Post Failed\n    Error: wrong description")
		return
	}

	if bytes.OutlineContentEasy != "test outline for bytes" {
		t.Error("\nCreate Bytes Post Failed\n    Error: wrong outline content")
		return
	}

	if bytes.DevStepsEasy != "steps 1 2 3 4" {
		t.Error("\nCreate Bytes Post Failed\n    Error: wrong dev steps")
		return
	}

	t.Log("\nCreate Bytes Post Succeeded")
}

func TestBytes_ToSQLNative(t *testing.T) {
	bytes, err := CreateBytes(69420, "test", "this is a test description", "test outline for bytes", "steps 1 2 3 4", 5, "#fff")
	if err != nil {
		t.Error("\nCreate Bytes Post Failed")
		return
	}

	statement, err := bytes.ToSQLNative()
	if err != nil {
		t.Error("\nCreate Bytes ToSQLNative Failed")
		return
	}

	if statement[0].Statement != "insert ignore into bytes(_id, name, description, outline_content, dev_steps, lang) values(?,?,?,?,?,?);" {
		t.Errorf("\nbytes to sql native failed\n    Error: incorrect statement returned")
		return
	}

	if len(statement[0].Values) != 5 {
		fmt.Println("number of values returned: ", len(statement[0].Values))
		t.Errorf("\nbytes to sql native failed\n    Error: incorrect values returned %v", statement[0].Values)
	}

	t.Log("\nbytes to sql native succeeded")
}

func TestBytesFromSQLNative(t *testing.T) {
	db, err := ti.CreateDatabase("gigo-dev-tidb", "4000", "mysql", "gigo-dev", "gigo-dev", "gigo_dev_test")
	if err != nil {
		t.Fatalf("\n%s failed\n    Error: %v", t.Name(), err)
	}

	bytes, err := CreateBytes(69420, "test", "this is a test description", "test outline for bytes", "steps 1 2 3 4", 5, "#fff")
	if err != nil {
		t.Fatalf("\n%s failed\n    Error: %v", t.Name(), err)
	}

	statement, err := bytes.ToSQLNative()
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

	rows, err := db.DB.Query("select * from bytes")
	if err != nil {
		t.Fatalf("\n%s failed\n    Error: %v", t.Name(), err)
	}

	if !rows.Next() {
		t.Fatalf("\n%s failed\n    Error: no rows found", t.Name())
		return
	}

	b, err := BytesFromSQLNative(rows)
	if err != nil {
		t.Fatalf("\n%s failed\n    Error: %v", t.Name(), err)
	}

	if b == nil {
		t.Fatalf("\n%s failed\n    Error: nil return", t.Name())
	}

	if b.ID != 69420 {
		t.Fatalf("\n%s failed\n    Error: wrong id", t.Name())
	}

	if b.Name != "test" {
		t.Fatalf("\n%s failed\n    Error: wrong name", t.Name())
	}

	if b.DescriptionEasy != "this is a test description" {
		t.Fatalf("\n%s failed\n    Error: wrong description", t.Name())
	}

	if b.OutlineContentEasy != "test outline for bytes" {
		t.Fatalf("\n%s failed\n    Error: wrong outline content", t.Name())
	}

	if b.DevStepsEasy != "steps 1 2 3 4" {
		t.Fatalf("\n%s failed\n    Error: wrong dev steps", t.Name())
	}

	t.Logf("\n%s succeeded", t.Name())
}
