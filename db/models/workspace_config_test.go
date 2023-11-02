package models

import (
	ti "github.com/gage-technologies/gigo-lib/db"
	"reflect"
	"testing"
)

func TestCreateWorkspaceConfig(t *testing.T) {
	wc := CreateWorkspaceConfig(
		69,
		"test",
		"test description",
		"fdsgsdfgsdfgsdfgdsfg",
		420,
		42069,
		[]int64{69420},
		[]ProgrammingLanguage{Go},
		int(0),
	)

	if wc.ID != 69 {
		t.Fatalf("%s failed\n    Error: id incorrect %d", t.Name(), wc.ID)
	}

	if wc.Title != "test" {
		t.Fatalf("%s failed\n    Error: title incorrect %s", t.Name(), wc.Title)
	}

	if wc.Description != "test description" {
		t.Fatalf("%s failed\n    Error: description incorrect %s", t.Name(), wc.Description)
	}

	if wc.Content != "fdsgsdfgsdfgsdfgdsfg" {
		t.Fatalf("%s failed\n    Error: content incorrect %s", t.Name(), wc.Content)
	}

	if wc.AuthorID != 420 {
		t.Fatalf("%s failed\n    Error: author id incorrect %d", t.Name(), wc.AuthorID)
	}

	if wc.Revision != 42069 {
		t.Fatalf("%s failed\n    Error: revision id incorrect %d", t.Name(), wc.Revision)
	}

	if !reflect.DeepEqual(wc.Tags, []int64{69420}) {
		t.Fatalf("%s failed\n    Error: tags incorrect %v", t.Name(), wc.Tags)
	}

	if !reflect.DeepEqual(wc.Languages, []ProgrammingLanguage{Go}) {
		t.Fatalf("%s failed\n    Error: languages incorrect %v", t.Name(), wc.Languages)
	}

	t.Logf("%s succeeded", t.Name())
}

func TestWorkspaceConfigFromSQLNative(t *testing.T) {
	db, err := ti.CreateDatabase("gigo-dev-tidb", "4000", "mysql", "gigo-dev", "gigo-dev", "gigo_dev_test")
	if err != nil {
		t.Fatalf("%s failed\n    Error: %v", t.Name(), err)
	}

	defer db.DB.Exec("delete from workspace_config")
	defer db.DB.Exec("delete from workspace_config_tags")
	defer db.DB.Exec("delete from workspace_config_langs")

	wc := CreateWorkspaceConfig(
		69,
		"test",
		"test description",
		"fdsgsdfgsdfgsdfgdsfg",
		420,
		42069,
		[]int64{69420},
		[]ProgrammingLanguage{Go},
		int(0),
	)

	statements, err := wc.ToSQLNative()
	if err != nil {
		t.Fatalf("%s failed\n    Error: %v", t.Name(), err)
	}

	for _, statement := range statements {
		stmt, err := db.DB.Prepare(statement.Statement)
		if err != nil {
			t.Fatalf("%s failed\n    Error: %v", t.Name(), err)
		}

		_, err = stmt.Exec(statement.Values...)
		if err != nil {
			t.Fatalf("%s failed\n    Error: %v", t.Name(), err)
		}
	}

	rows, err := db.DB.Query("select * from workspace_config where _id = ?", 69)
	if err != nil {
		t.Fatalf("%s failed\n    Error: %v", t.Name(), err)
	}

	if !rows.Next() {
		t.Fatalf("%s failed\n    Error: not found", t.Name())
	}

	rec, err := WorkspaceConfigFromSQLNative(db, rows)
	if err != nil {
		t.Fatalf("%s failed\n    Error: %v", t.Name(), err)
	}

	if rec == nil {
		t.Fatalf("%s failed\n    Error: config returned as nil", t.Name())
	}

	if rec.ID != 69 {
		t.Fatalf("%s failed\n    Error: id incorrect %d", t.Name(), rec.ID)
	}

	if rec.Title != "test" {
		t.Fatalf("%s failed\n    Error: title incorrect %s", t.Name(), rec.Title)
	}

	if rec.Description != "test description" {
		t.Fatalf("%s failed\n    Error: description incorrect %s", t.Name(), rec.Description)
	}

	if rec.Content != "fdsgsdfgsdfgsdfgdsfg" {
		t.Fatalf("%s failed\n    Error: content incorrect %s", t.Name(), rec.Content)
	}

	if rec.AuthorID != 420 {
		t.Fatalf("%s failed\n    Error: author id incorrect %d", t.Name(), rec.AuthorID)
	}

	if rec.Revision != 42069 {
		t.Fatalf("%s failed\n    Error: revision id incorrect %d", t.Name(), rec.Revision)
	}

	if !reflect.DeepEqual(rec.Tags, []int64{69420}) {
		t.Fatalf("%s failed\n    Error: tags incorrect %v", t.Name(), rec.Tags)
	}

	if !reflect.DeepEqual(rec.Languages, []ProgrammingLanguage{Go}) {
		t.Fatalf("%s failed\n    Error: languages incorrect %v", t.Name(), rec.Languages)
	}

	t.Logf("%s succeeded", t.Name())
}
