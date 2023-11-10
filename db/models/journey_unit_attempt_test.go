package models

import (
	"context"
	ti "github.com/gage-technologies/gigo-lib/db"
	"go.opentelemetry.io/otel"
	"reflect"
	"testing"
	"time"
)

func TestCreateJourneyUnitAttempt(t *testing.T) {
	journey, err := CreateJourneyUnitAttempt(
		1,
		69420,
		2,
		"test-title",
		UnitFocusFromString("Fullstack"),
		[]ProgrammingLanguage{Go, Python, JavaScript},
		"this a test description",
		69420,
		time.Now(),
		time.Now(),
		nil,
		Tier1,
		&DefaultWorkspaceSettings,
		69420,
		0,
		nil,
	)
	if err != nil {
		t.Error("\nCreate Journey Unit Failed\n    Error: ", err)
		return
	}

	if journey == nil {
		t.Error("\nCreate Journey Unit Failed\n    Error: creation returned nil")
		return
	}

	if journey.ID != 1 {
		t.Error("\nCreate Journey Unit Failed\n    Error: wrong id")
		return
	}

	if journey.UserID != 69420 {
		t.Error("\nCreate Journey Unit Failed\n    Error: wrong user id")
		return
	}

	if journey.ParentUnit != 2 {
		t.Error("\nCreate Journey Unit Failed\n    Error: wrong parent id")
		return
	}

	if journey.Title != "test-title" {
		t.Error("\nCreate Journey Unit Failed\n    Error: wrong name")
		return
	}

	if journey.UnitFocus.String() != "fullstack" {
		t.Error("\nCreate Journey Unit Failed\n    Error: wrong unit focus")
		return
	}

	if !reflect.DeepEqual(journey.LanguageList, []*JourneyUnitLanguages{
		{
			1,
			"Go",
			true,
		},
		{
			1,
			"Python",
			true,
		},
		{
			1,
			"JavaScript",
			true,
		},
	}) {
		t.Errorf("\nCreate Journey Unit Failed\n    Error: wrong language list: \n%v\n%v\n%v", journey.LanguageList[0], journey.LanguageList[1], journey.LanguageList[2])
		return
	}

	if journey.Description != "this a test description" {
		t.Error("\nCreate Journey Unit Failed\n    Error: wrong description")
		return
	}

	if journey.WorkspaceConfig != 69420 {
		t.Error("\nCreate Journey Unit Failed\n    Error: wrong workspace config")
		return
	}

	t.Log("\nCreate Journey Unit Succeeded")
}

func TestJourneyUnitAttempt_ToSQLNative(t *testing.T) {
	journey, err := CreateJourneyUnitAttempt(
		1,
		69420,
		2,
		"test-title",
		UnitFocusFromString("Fullstack"),
		[]ProgrammingLanguage{Go, Python, JavaScript},
		"this a test description",
		69420,
		time.Now(),
		time.Now(),
		nil,
		Tier1,
		&DefaultWorkspaceSettings,
		69420,
		0,
		nil,
	)
	if err != nil {
		t.Error("\nCreate Journey Unit Failed\n    Error: ", err)
		return
	}

	statement, err := journey.ToSQLNative()
	if err != nil {
		t.Error("\nCreate Journey Unit Failed\n    Error: ", err)
		return
	}

	if statement[len(statement)-1].Statement != "insert ignore into journey_unit_attempts (_id, title, user_id, parent_unit, unit_focus, description, repo_id, created_at, updated_at, challenge_cost, completions, attempts, tier, embedded, workspace_config, workspace_config_revision, workspace_settings, estimated_tutorial_time) values (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?);" {
		t.Error("\nJourney Unit To SQL Native Failed\n    Error: wrong statement")
		return
	}

	if len(statement[len(statement)-1].Values) != 18 {
		t.Error("\nJourney Unit To SQL Native Failed\n    Error: wrong number of values")
		return
	}

	t.Log("\nJourney Unit To SQL Native Succeeded")
}

func TestJourneyUnitAttemptFromSQLNative(t *testing.T) {
	ctx, span := otel.Tracer("gigo-core").Start(context.Background(), "gigo-lib")
	defer span.End()

	callerName := "JourneyUnitAttemptFromSQLNativeTest"

	db, err := ti.CreateDatabase("gigo-dev-tidb", "4000", "mysql", "gigo-dev", "gigo-dev", "gigo_dev_test")
	if err != nil {
		t.Error("\nInitialize Database Failed\n    Error: ", err)
	}

	defer db.DB.Exec("delete from journey_unit_attempts")

	journey, err := CreateJourneyUnitAttempt(
		1,
		69420,
		2,
		"test-title",
		UnitFocusFromString("Fullstack"),
		[]ProgrammingLanguage{Go, Python, JavaScript},
		"this a test description",
		69420,
		time.Now(),
		time.Now(),
		nil,
		Tier1,
		&DefaultWorkspaceSettings,
		69420,
		0,
		nil,
	)
	if err != nil {
		t.Error("\nCreate Journey Unit Failed\n    Error: ", err)
		return
	}

	statements, err := journey.ToSQLNative()
	if err != nil {
		t.Error("\nCreate Journey Unit Failed\n    Error: ", err)
		return
	}

	for _, statement := range statements {
		stmt, err := db.DB.Prepare(statement.Statement)
		if err != nil {
			t.Error("\nPrepare Statement Failed\n    Error: ", err)
		}

		_, err = stmt.Exec(statement.Values...)
		if err != nil {
			t.Error("\nExecute Statement Failed\n    Error: ", err)
		}
	}

	rows, err := db.QueryContext(ctx, &span, &callerName, "select * from journey_unit_attempts where _id = ?", journey.ID)
	if err != nil {
		t.Error("\nCreate Journey Unit Failed\n    Error: ", err)
		return
	}

	defer rows.Close()

	var j *JourneyUnitAttempt

	for rows.Next() {
		j, err = JourneyUnitAttemptFromSQLNative(db, rows)
		if err != nil {
			t.Error("\nJourney Unit From SQL Native Failed\n    Error: ", err)
			return
		}
	}

	if j == nil {
		t.Error("\nJourney Unit From SQL Native Failed\n    Error: creation returned nil")
		return
	}

	if j.ID != 1 {
		t.Error("\nJourney Unit From SQL Native Failed\n    Error: wrong id")
		return
	}

	if j.Title != "test-title" {
		t.Error("\nJourney Unit From SQL Native Failed\n    Error: wrong name")
		return
	}

	if j.UnitFocus != Fullstack {
		t.Error("\nJourney Unit From SQL Native Failed\n    Error: wrong unit focus")
		return
	}

	if !reflect.DeepEqual(journey.LanguageList, []*JourneyUnitLanguages{
		{journey.ID,
			Go.String(),
			true,
		}, {
			journey.ID,
			Python.String(),
			true,
		},
		{
			journey.ID,
			JavaScript.String(),
			true,
		},
	}) {
		t.Error("\nJourney Unit From SQL Native Failed\n    Error: wrong language list")
		return
	}

	if j.Description != "this a test description" {
		t.Error("\nJourney Unit From SQL Native Failed\n    Error: wrong description")
		return
	}

	if j.WorkspaceConfig != 69420 {
		t.Errorf("\nJourney Unit From SQL Native Failed\n    Error: wrong workspace config: %v", j.WorkspaceConfig)
		return
	}

	t.Log("\nJourney Unit From SQL Native Succeeded")
}
