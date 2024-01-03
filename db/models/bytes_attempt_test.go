package models

import (
	"fmt"
	ti "github.com/gage-technologies/gigo-lib/db"
	"testing"
)

func TestCreateByteAttempts(t *testing.T) {
	byteAttempts, err := CreateByteAttempts(42069, 69, 420, "content for the attempt")
	if err != nil {
		t.Error("\nCreate ByteAttempts Post Failed")
		return
	}

	if byteAttempts == nil {
		t.Error("\nCreate ByteAttempts Post Failed\n    Error: creation returned nil")
		return
	}

	if byteAttempts.ID != 42069 {
		t.Error("\nCreate ByteAttempts Post Failed\n    Error: wrong id")
		return
	}

	if byteAttempts.ByteID != 69 {
		t.Error("\nCreate ByteAttempts Post Failed\n    Error: wrong byte id")
		return
	}

	if byteAttempts.AuthorID != 420 {
		t.Error("\nCreate ByteAttempts Post Failed\n    Error: wrong author id")
		return
	}

	if byteAttempts.Content != "content for the attempt" {
		t.Error("\nCreate ByteAttempts Post Failed\n    Error: wrong content")
		return
	}

	t.Log("\nCreate ByteAttempts Post Succeeded")
}

func TestByteAttempts_ToSQLNative(t *testing.T) {
	byteAttempts, err := CreateByteAttempts(42069, 69, 420, "content for the attempt")
	if err != nil {
		t.Error("\nCreate ByteAttempts Post Failed")
	}

	statement, err := byteAttempts.ToSQLNative()
	if err != nil {
		t.Error("\nCreate ByteAttempts ToSQLNative Failed")
	}

	if statement[0].Statement != "insert ignore into byte_attempts(_id, byte_id, author_id, content) values(?,?,?,?);" {
		t.Errorf("\nbyte attempts to sql native failed\n    Error: incorrect statement returned")
		return
	}

	if len(statement[0].Values) != 4 {
		fmt.Println("number of values returned: ", len(statement[0].Values))
		t.Errorf("\nbyte attempts to sql native failed\n    Error: incorrect values returned %v", statement[0].Values)
	}

	t.Log("\nbyte attempts to sql native succeeded")

}

func TestByteAttemptsFromSQLNative(t *testing.T) {
	db, err := ti.CreateDatabase("gigo-dev-tidb", "4000", "mysql", "gigo-dev", "gigo-dev", "gigo_dev_test")
	if err != nil {
		t.Fatalf("\n%s failed\n    Error: %v", t.Name(), err)
	}

	byteAttempt, err := CreateByteAttempts(42069, 69, 420, "content for the attempt")
	if err != nil {
		t.Fatalf("\n%s failed\n    Error: %v", t.Name(), err)
	}

	statement, err := byteAttempt.ToSQLNative()
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

	rows, err := db.DB.Query("select * from byte_attempts")
	if err != nil {
		t.Fatalf("\n%s failed\n    Error: %v", t.Name(), err)
	}

	if !rows.Next() {
		t.Fatalf("\n%s failed\n    Error: no rows found", t.Name())
		return
	}

	b, err := ByteAttemptsFromSQLNative(rows)
	if err != nil {
		t.Fatalf("\n%s failed\n    Error: %v", t.Name(), err)
	}

	if b == nil {
		t.Fatalf("\n%s failed\n    Error: nil return", t.Name())
	}

	if b.ID != 42069 {
		t.Fatalf("\n%s failed\n    Error: wrong id", t.Name())
	}

	if b.ByteID != 69 {
		t.Fatalf("\n%s failed\n    Error: wrong byte id", t.Name())
	}

	if b.AuthorID != 420 {
		t.Fatalf("\n%s failed\n    Error: wrong author id", t.Name())
	}

	if b.Content != "content for the attempt" {
		t.Fatalf("\n%s failed\n    Error: wrong content", t.Name())
	}

	t.Logf("\n%s succeeded", t.Name())
}