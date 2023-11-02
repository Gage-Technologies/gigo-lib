package models

import (
	"fmt"
	ti "github.com/gage-technologies/gigo-lib/db"
	"testing"
	"time"
)

func TestCreateRecommendedPost(t *testing.T) {
	n := time.Now()
	e := time.Now().Add(time.Minute)
	rec, err := CreateRecommendedPost(69420, 6969, 420, RecommendationTypeSematic, 69, 0.9, n, e, Tier1)
	if err != nil {
		t.Error("\nCreate Recommended Post Failed")
		return
	}

	if rec == nil {
		t.Error("\nCreate Recommended Post Failed\n    Error: creation returned nil")
		return
	}

	if rec.ID != 69420 {
		t.Error("\nCreate Recommended Post Failed\n    Error: wrong id")
		return
	}

	if rec.UserID != 6969 {
		t.Error("\nCreate Recommended Post Failed\n    Error: wrong author id")
		return
	}

	if rec.PostID != 420 {
		t.Error("\nCreate Recommended Post Failed\n    Error: wrong post id")
		return
	}

	if rec.Type != RecommendationTypeSematic {
		t.Error("\nCreate Recommended Post Failed\n    Error: wrong recommendation type")
		return
	}

	if rec.Score != .9 {
		t.Error("\nCreate Recommended Post Failed\n    Error: wrong score")
		return
	}

	if rec.ReferenceID != 69 {
		t.Error("\nCreate Recommended Post Failed\n    Error: wrong reference id")
		return
	}

	if rec.CreatedAt.Unix() != n.Unix() {
		t.Error("\nCreate Recommended Post Failed\n    Error: wrong created at")
		return
	}

	if rec.ExpiresAt.Unix() != e.Unix() {
		t.Error("\nCreate Recommended Post Failed\n    Error: wrong expires at")
		return
	}

	t.Log("\nCreate Recommended Post Succeeded")
}

func TestRecommendedPost_ToSQLNative(t *testing.T) {
	n := time.Now()
	e := time.Now().Add(time.Minute)
	rec, err := CreateRecommendedPost(69420, 6969, 420, RecommendationTypeSematic, 69, 0.9, n, e, Tier1)
	if err != nil {
		t.Error("\nCreate Recommended Post Failed")
		return
	}

	statement := rec.ToSQLNative()

	if statement.Statement != "insert ignore into recommended_post(_id, user_id, post_id, type, reference_id, score, created_at, expires_at, reference_tier) values(?, ?, ?, ?, ?, ?, ?, ?, ?);" {
		t.Errorf("\nRec post to sql native failed\n    Error: incorrect statement returned")
		return
	}

	if len(statement.Values) != 9 {
		fmt.Println("number of values returned: ", len(statement.Values))
		t.Errorf("\nRec post to sql native failed\n    Error: incorrect values returned %v", statement.Values)
		return
	}

	t.Log("\nRec post to sql native succeeded")
}

func TestRecommendedPostFromSQLNative(t *testing.T) {
	db, err := ti.CreateDatabase("gigo-dev-tidb", "4000", "mysql", "gigo-dev", "gigo-dev", "gigo_dev_test")
	if err != nil {
		t.Error("\nInitialize recommended post table sql failed\n    Error: ", err)
		return
	}

	defer db.DB.Exec("delete from recommended_post")

	n := time.Now()
	e := time.Now().Add(time.Minute)
	rec, err := CreateRecommendedPost(69420, 6969, 420, RecommendationTypeSematic, 69, 0.9, n, e, Tier1)
	if err != nil {
		t.Error("\nCreate Recommended Post Failed")
		return
	}

	statement := rec.ToSQLNative()

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

	rows, err := db.DB.Query("select * from recommended_post")
	if err != nil {
		t.Error("\nRec post from sql native failed\n    Error: ", err)
		return
	}

	if !rows.Next() {
		t.Error("\nRec post from sql native failed\n    Error: no rows found")
		return
	}

	r, err := RecommendedPostFromSQLNative(rows)
	if err != nil {
		t.Error("\nRec post from sql native failed\n    Error: ", err)
		return
	}

	if r == nil {
		t.Error("\nRec post from sql native failed\n    Error: creation returned nil")
		return
	}

	if r.ID != 69420 {
		t.Error("\nRec post from sql native failed\n    Error: wrong id")
		return
	}

	if r.UserID != 6969 {
		t.Error("\nRec post from sql native failed\n    Error: wrong author id")
		return
	}

	if r.PostID != 420 {
		t.Error("\nRec post from sql native failed\n    Error: wrong post id")
		return
	}

	if r.Type != RecommendationTypeSematic {
		t.Error("\nRec post from sql native failed\n    Error: wrong recommendation type")
		return
	}

	if r.Score != .9 {
		t.Error("\nRec post from sql native failed\n    Error: wrong score")
		return
	}

	if r.ReferenceID != 69 {
		t.Error("\nRec post from sql native failed\n    Error: wrong reference id")
		return
	}

	t.Log("\nRec post from sql native succeeded")
}
