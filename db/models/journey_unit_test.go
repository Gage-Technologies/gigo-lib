package models

import (
	"context"
	"fmt"
	ti "github.com/gage-technologies/gigo-lib/db"
	"go.opentelemetry.io/otel"
	"reflect"
	"testing"
	"time"
)

func TestCreateJourneyUnit(t *testing.T) {
	journey, err := CreateJourneyUnit(
		1,
		"Journey",
		UnitFocusFromString("Fullstack"),
		[]ProgrammingLanguage{Go, Python, JavaScript},
		"this is a test description",
		69420,
		time.Now(),
		time.Now(),
		nil,
		Tier1,
		&DefaultWorkspaceSettings,
		69420,
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

	if journey.Title != "Journey" {
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
			false,
		},
		{
			1,
			"Python",
			false,
		},
		{
			1,
			"JavaScript",
			false,
		},
	}) {
		t.Errorf("\nCreate Journey Unit Failed\n    Error: wrong language list: \n%v\n%v\n%v", journey.LanguageList[0], journey.LanguageList[1], journey.LanguageList[2])
		return
	}

	if journey.Description != "this is a test description" {
		t.Error("\nCreate Journey Unit Failed\n    Error: wrong description")
		return
	}

	if journey.WorkspaceConfig != 69420 {
		t.Error("\nCreate Journey Unit Failed\n    Error: wrong workspace config")
		return
	}

	t.Log("\nCreate Journey Unit Succeeded")
}

func TestJourneyUnit_ToSQLNative(t *testing.T) {
	journey, err := CreateJourneyUnit(
		1,
		"Journey",
		UnitFocusFromString("Fullstack"),
		[]ProgrammingLanguage{Go, Python, JavaScript},
		"this is a test description",
		69420,
		time.Now(),
		time.Now(),
		nil,
		Tier1,
		&DefaultWorkspaceSettings,
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

	if statement[len(statement)-1].Statement != "insert ignore into journey_units (_id, title, unit_focus, description, repo_id, created_at, updated_at, challenge_cost, completions, attempts, tier, embedded, workspace_config, workspace_config_revision, workspace_settings, estimated_tutorial_time) values (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?);" {
		t.Error("\nJourney Unit To SQL Native Failed\n    Error: wrong statement")
		return
	}

	if len(statement[len(statement)-1].Values) != 16 {
		t.Error("\nJourney Unit To SQL Native Failed\n    Error: wrong number of values")
		return
	}

	t.Log("\nJourney Unit To SQL Native Succeeded")
}

func TestJourneyUnitFromSQLNative(t *testing.T) {
	ctx, span := otel.Tracer("gigo-core").Start(context.Background(), "gigo-lib")
	defer span.End()

	callerName := "JourneyUnitFromSQLNativeTest"

	db, err := ti.CreateDatabase("gigo-dev-tidb", "4000", "mysql", "gigo-dev", "gigo-dev", "gigo_dev_test")
	if err != nil {
		t.Error("\nInitialize Database Failed\n    Error: ", err)
	}

	defer db.DB.Exec("delete from journey_units")

	journey, err := CreateJourneyUnit(
		1,
		"Journey",
		UnitFocusFromString("Fullstack"),
		[]ProgrammingLanguage{Go, Python, JavaScript},
		"this is a test description",
		69420,
		time.Now(),
		time.Now(),
		nil,
		Tier1,
		&DefaultWorkspaceSettings,
		69420,
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

	rows, err := db.QueryContext(ctx, &span, &callerName, "select * from journey_units where _id = ?", journey.ID)
	if err != nil {
		t.Error("\nCreate Journey Unit Failed\n    Error: ", err)
		return
	}

	defer rows.Close()

	var j *JourneyUnit

	for rows.Next() {
		j, err = JourneyUnitFromSQLNative(db, rows)
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

	if j.Title != "Journey" {
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
			false,
		}, {
			journey.ID,
			Python.String(),
			false,
		},
		{
			journey.ID,
			JavaScript.String(),
			false,
		},
	}) {
		t.Error("\nJourney Unit From SQL Native Failed\n    Error: wrong language list")
		return
	}

	if j.Description != "this is a test description" {
		t.Error("\nJourney Unit From SQL Native Failed\n    Error: wrong description")
		return
	}

	fmt.Println("journey test wsconfig: ", journey)
	fmt.Println("journey test wsconfig revision: ", j)
	if j.WorkspaceConfig != 69420 {
		t.Errorf("\nJourney Unit From SQL Native Failed\n    Error: wrong workspace config: %v", j.WorkspaceConfig)
		return
	}

	t.Log("\nJourney Unit From SQL Native Succeeded")
}
