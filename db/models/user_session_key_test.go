package models

import (
	"fmt"
	ti "github.com/gage-technologies/gigo-lib/db"
	"github.com/gage-technologies/gigo-lib/session"
	"math"
	"testing"
	"time"
)

func TestCreateUserSessionKey(t *testing.T) {
	e := time.Now()
	usk := CreateUserSessionKey(42069, "test", e)

	if usk.ID != 42069 {
		t.Fatalf("\n%v failed\n    Error: incorrect id returned", t.Name())
	}

	if usk.Key != "test" {
		t.Fatalf("\n%v failed\n    Error: incorrect key returned", t.Name())
	}

	if usk.Expiration.Unix() != e.Unix() {
		t.Fatalf("\n%v failed\n    Error: incorrect expiration returned", t.Name())
	}

	t.Logf("\n%s succeeded", t.Name())
}

func TestUserSessionKeyToFromSQLNative(t *testing.T) {
	db, err := ti.CreateDatabase("gigo-dev-tidb", "4000", "mysql", "gigo-dev", "gigo-dev", "gigo_dev_test")
	if err != nil {
		t.Fatalf("\n%v failed\n    Error: %v", t.Name(), err)
	}

	defer db.DB.Exec("DROP TABLE user_session_key")

	pass, err := session.GenerateServicePassword()
	if err != nil {
		t.Fatalf("\n%v failed\n    Error: %v", t.Name(), err)
	}

	e := time.Now()
	usk := CreateUserSessionKey(42069, pass, e)

	statement, err := usk.ToSQLNative()
	if err != nil {
		t.Fatalf("\n%v failed\n    Error: %v", t.Name(), err)
	}

	if statement[0].Statement != "insert ignore into user_session_key(_id, _key, expiration) values (?, ?, ?)" {
		t.Fatalf("\n%v failed\n    Error: incorrect statement", t.Name())
	}

	if len(statement[0].Values) != 3 {
		t.Fatalf("\n%v failed\n    Error: incorrect values %v", t.Name(), statement[0].Values)
	}

	stmt, err := db.DB.Prepare(statement[0].Statement)
	if err != nil {
		t.Fatalf("\n%v failed\n    Error: %v", t.Name(), err)
	}

	_, err = stmt.Exec(statement[0].Values...)
	if err != nil {
		t.Fatalf("\n%v failed\n    Error: %v", t.Name(), err)
	}

	rows, err := db.DB.Query("select * from user_session_key where _id = 42069")
	if err != nil {
		t.Fatalf("\n%v failed\n    Error: %v", t.Name(), err)
	}

	if !rows.Next() {
		t.Fatalf("\n%v failed\n    Error: no results", t.Name())
	}

	usk2, err := UserSessionKeyFromSQLNative(rows)
	if err != nil {
		t.Fatalf("\n%v failed\n    Error: %v", t.Name(), err)
	}

	if usk2.ID != 42069 {
		t.Fatalf("\n%v failed\n    Error: incorrect id returned", t.Name())
	}

	if usk2.Key != pass {
		fmt.Println(usk2.Key)
		fmt.Println(pass)
		t.Fatalf("\n%v failed\n    Error: incorrect key returned", t.Name())
	}

	if math.Abs(float64(usk2.Expiration.Unix()-e.Unix())) > 1 {
		t.Fatalf("\n%v failed\n    Error: incorrect expiration returned", t.Name())
	}

	t.Logf("\nTagFromSQLNative succeeded")
}
