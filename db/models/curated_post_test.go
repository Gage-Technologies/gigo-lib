package models

import (
	ti "github.com/gage-technologies/gigo-lib/db"
	"testing"
)

func TestCreateCuratedPost(t *testing.T) {
	rec, err := CreateCuratedPost(1, 100, []ProficiencyType{Beginner, Intermediate}, Go)
	if err != nil {
		t.Error("Create Curated Post Failed:", err)
		return
	}

	if rec == nil {
		t.Error("Create Curated Post Failed: creation returned nil")
		return
	}

	if rec.ID != 1 {
		t.Error("Create Curated Post Failed: wrong ID")
		return
	}

	if rec.PostID != 100 {
		t.Error("Create Curated Post Failed: wrong PostID")
		return
	}
	t.Log("Create Curated Post Succeeded")
}

func TestCuratedPost_ToSQLNative(t *testing.T) {
	rec, err := CreateCuratedPost(1, 100, []ProficiencyType{Beginner, Intermediate}, Go)
	if err != nil {
		t.Error("Create Curated Post Failed:", err)
		return
	}

	statement := rec.ToSQLNative()

	if statement.Statement != "insert ignore into curated_post(_id, post_id, proficiency_type, post_language) values(?, ?, ?, ?);" {
		t.Errorf("Curated post to sql native failed: incorrect statement returned")
		return
	}

	if len(statement.Values) != 4 {
		t.Errorf("Curated post to sql native failed: incorrect values returned %v", statement.Values)
		return
	}

	t.Log("Curated post to sql native succeeded")
}

func TestCuratedPostFromSQLNative(t *testing.T) {
	db, err := ti.CreateDatabase("gigo-dev-tidb", "4000", "mysql", "gigo-dev", "gigo-dev", "gigo_dev_test")
	if err != nil {
		t.Error("\nInitialize recommended post table sql failed\n    Error: ", err)
		return
	}

	defer db.DB.Exec("DROP TABLE curated_post")

	// Prepare and insert a test record
	post, err := CreateCuratedPost(1, 100, []ProficiencyType{Beginner, Intermediate}, Go)
	if err != nil {
		t.Error("Create Curated Post Failed:", err)
		return
	}

	statement := post.ToSQLNative()
	stmt, err := db.DB.Prepare(statement.Statement)
	if err != nil {
		t.Error("Curated post from sql native failed:", err)
		return
	}

	_, err = stmt.Exec(statement.Values...)
	if err != nil {
		t.Error("Curated post from sql native failed:", err)
		return
	}

	// Query and check
	rows, err := db.DB.Query("select * from curated_post")
	if err != nil {
		t.Error("Curated post from sql native failed:", err)
		return
	}

	if !rows.Next() {
		t.Error("Curated post from sql native failed: no rows found")
		return
	}

	rec, err := CuratedPostFromSQLNative(rows)
	if err != nil {
		t.Error("Curated post from sql native failed:", err)
		return
	}

	if rec == nil {
		t.Error("Curated post from sql native failed: creation returned nil")
		return
	}

	if rec.ID != 1 {
		t.Error("Curated post from sql native failed: wrong ID")
		return
	}

	if rec.PostID != 100 {
		t.Error("Curated post from sql native failed: wrong PostID")
		return
	}
}
