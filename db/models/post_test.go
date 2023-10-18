package models

import (
	ti "github.com/gage-technologies/gigo-lib/db"
	"reflect"
	"sort"
	"testing"
	"time"
)

func TestCreatePost(t *testing.T) {
	id := int64(5)
	attempt, err := CreatePost(
		69420,
		"test",
		"asdfasdf",
		"autor",
		42069,
		time.Now(),
		time.Now(),
		69,
		1,
		[]int64{1, 2, 3},
		&id,
		6969,
		20,
		40,
		24,
		27,
		[]ProgrammingLanguage{Go},
		PublicVisibility,
		[]int64{4, 5, 6},
		nil,
		nil,
		6942969,
		3,
		&DefaultWorkspaceSettings,
		false,
		false,
		nil,
	)
	if err != nil {
		t.Fatalf("%s failed\n    Error: %v", t.Name(), err)
	}

	if attempt == nil {
		t.Fatalf("%s failed\n    Error: returned nil", t.Name())
	}

	if attempt.ID != 69420 {
		t.Fatalf("%s failed\n    Error: wrong id %v", t.Name(), attempt.ID)
	}

	if attempt.AuthorID != 42069 {
		t.Fatalf("%s failed\n    Error: wrong author id %v", t.Name(), attempt.AuthorID)
	}

	if attempt.Published != false {
		t.Fatalf("%s failed\n    Error: wrong published %v", t.Name(), attempt.Published)
	}

	if attempt.Visibility != PublicVisibility {
		t.Fatalf("%s failed\n    Error: wrong visibility %v", t.Name(), attempt.Visibility)
	}

	if !reflect.DeepEqual(attempt.Languages, []ProgrammingLanguage{Go}) {
		t.Fatalf("%s failed\n    Error: wrong languages %v", t.Name(), attempt.Languages)
	}

	if !reflect.DeepEqual(attempt.Tags, []int64{4, 5, 6}) {
		t.Fatalf("%s failed\n    Error: wrong tags %v", t.Name(), attempt.Tags)
	}

	if attempt.WorkspaceConfig != 6942969 {
		t.Fatalf("%s failed\n    Error: wrong workspace config %v", t.Name(), attempt.WorkspaceConfig)
	}

	if attempt.WorkspaceConfigRevision != 3 {
		t.Fatalf("%s failed\n    Error: wrong workspace config revision %v", t.Name(), attempt.WorkspaceConfigRevision)
	}

	t.Logf("\n%s succeeded", t.Name())
}

func TestPost_ToSQLNative(t *testing.T) {
	id := int64(5)
	attempt, err := CreatePost(
		69420,
		"test",
		"description",
		"autor",
		42069,
		time.Now(),
		time.Now(),
		69,
		2,
		[]int64{1, 2, 3},
		&id,
		6969,
		20,
		40,
		24,
		27,
		[]ProgrammingLanguage{Go},
		PublicVisibility,
		[]int64{4, 5, 6},
		nil,
		nil,
		6942069,
		8,
		&DefaultWorkspaceSettings,
		false,
		false,
		nil,
	)
	if err != nil {
		t.Fatalf("%s failed\n    Error: %v", t.Name(), err)
	}

	statement, err := attempt.ToSQLNative()
	if err != nil {
		t.Fatalf("%s failed\n    Error: %v", t.Name(), err)
	}

	if statement[0].Statement != "insert ignore into post(_id, title, description, author, author_id, created_at, updated_at, repo_id, top_reply, tier, coffee, post_type, views, completions, attempts, published, visibility, stripe_price_id, challenge_cost, workspace_config, workspace_config_revision, workspace_settings, leads, embedded, exclusive_content, share_hash) values(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, uuid_to_bin(?));" {
		t.Fatalf("%s failed\n    Error: incorrect statement %v", t.Name(), statement[0].Statement)
	}

	if len(statement[0].Values) != 26 {
		t.Fatalf("%s failed\n    Error: incorrect number of values %v", t.Name(), len(statement[0].Values))
	}

	t.Logf("%v succeeded", t.Name())
}

func TestPostFromSQLNative(t *testing.T) {
	db, err := ti.CreateDatabase("gigo-dev-tidb", "4000", "mysql", "gigo-dev", "gigo-dev", "gigo_dev_test")
	if err != nil {
		t.Error("\nInitialize recommended post table sql failed\n    Error: ", err)
		return
	}

	defer db.DB.Exec("DROP TABLE post")
	defer db.DB.Exec("DROP TABLE post_awards")
	defer db.DB.Exec("DROP TABLE post_tags")

	id := int64(5)
	post, err := CreatePost(
		69420,
		"test",
		"asdfasdf",
		"autor",
		42069,
		time.Now(),
		time.Now(),
		69,
		3,
		[]int64{1, 2, 3},
		&id,
		6969,
		20,
		40,
		24,
		27,
		[]ProgrammingLanguage{Go, Java, Python},
		PublicVisibility,
		[]int64{4, 5, 6},
		nil,
		nil,
		42069420,
		12,
		&DefaultWorkspaceSettings,
		false,
		false,
		nil,
	)
	if err != nil {
		t.Error("\nCreate Attempt Post Failed")
		return
	}

	statements, err := post.ToSQLNative()
	if err != nil {
		t.Error("\nCreate Attempt Post ToSQLNative Failed")
		return
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

	rows, err := db.DB.Query("select * from post where _id = ?", 69420)
	if err != nil {
		t.Fatalf("%s failed\n    Error: %v", t.Name(), err)
	}

	if !rows.Next() {
		t.Fatalf("%s failed\n    Error: row not found", t.Name())
	}

	rec, err := PostFromSQLNative(db, rows)
	if err != nil {
		t.Fatalf("%s failed\n    Error: %v", t.Name(), err)
	}

	if err != nil {
		t.Fatalf("%s failed\n    Error: post returned as nil", t.Name())
	}

	if rec.ID != 69420 {
		t.Fatalf("%s failed\n    Error: incorrect id returned %v", t.Name(), rec.ID)
	}

	if rec.AuthorID != 42069 {
		t.Fatalf("%s failed\n    Error: incorrect author id returned %v", t.Name(), rec.AuthorID)
	}

	if rec.Published != false {
		t.Fatalf("%s failed\n    Error: incorrect published value returned %v", t.Name(), rec.Published)
	}

	if rec.Visibility != PublicVisibility {
		t.Fatalf("%s failed\n    Error: incorrect visibility value returned %v", t.Name(), rec.Visibility)
	}

	if rec.WorkspaceConfig != 42069420 {
		t.Fatalf("%s failed\n    Error: incorrect workspace value returned %v", t.Name(), rec.WorkspaceConfig)
	}

	if rec.WorkspaceConfigRevision != 12 {
		t.Fatalf("%s failed\n    Error: incorrect workspace revision value returned %v", t.Name(), rec.WorkspaceConfigRevision)
	}

	if !reflect.DeepEqual(rec.Awards, []int64{1, 2, 3}) {
		t.Fatalf("%s failed\n    Error: incorrect awards value returned %v", t.Name(), rec.Awards)
	}

	if !reflect.DeepEqual(rec.Tags, []int64{4, 5, 6}) {
		t.Fatalf("%s failed\n    Error: incorrect tags value returned %v", t.Name(), rec.Tags)
	}

	sort.Slice(rec.Languages, func(i, j int) bool {
		return rec.Languages[i] < rec.Languages[j]
	})

	if !reflect.DeepEqual(rec.Languages, []ProgrammingLanguage{Java, Python, Go}) {
		t.Fatalf("%s failed\n    Error: incorrect languages value returned %v", t.Name(), rec.Languages)
	}

	if !reflect.DeepEqual(*rec.WorkspaceSettings, DefaultWorkspaceSettings) {
		t.Fatalf("%s failed\n    Error: incorrect workspace settings returned %v", t.Name(), rec.WorkspaceSettings)
	}

	t.Logf("%s succeeded", t.Name())
}
