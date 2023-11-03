package models

import (
	ti "github.com/gage-technologies/gigo-lib/db"
	"reflect"
	"testing"
)

func TestCreateJourneyUnit(t *testing.T) {
	journey, err := CreateJourneyUnit(
		1,
		"Journey",
		"fullstack",
		[]ProgrammingLanguage{Go, Python, JavaScript},
		"this is a test description",
		69420,
		&DefaultWorkspaceSettings,
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

	if journey.UnitFocus != "fullstack" {
		t.Error("\nCreate Journey Unit Failed\n    Error: wrong unit focus")
		return
	}

	if !reflect.DeepEqual(journey.LanguageList, []ProgrammingLanguage{Go, Python, JavaScript}) {
		t.Error("\nCreate Journey Unit Failed\n    Error: wrong language list")
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
		"fullstack",
		[]ProgrammingLanguage{Go, Python, JavaScript},
		"this is a test description",
		69420,
		&DefaultWorkspaceSettings,
	)
	if err != nil {
		t.Error("\nCreate Journey Unit Failed\n    Error: ", err)
		return
	}

	statement := journey.ToSQLNative()

	if statement[0].Statement != "insert ignore into journey_unit (_id, title, unit_focus, language_list, description, workspace_config, workspace_settings) values (?,?,?,?,?,?,?);" {
		t.Error("\nJourney Unit To SQL Native Failed\n    Error: wrong statement")
		return
	}

	if len(statement[0].Values) != 7 {
		t.Error("\nJourney Unit To SQL Native Failed\n    Error: wrong id")
		return
	}

	t.Log("\nJourney Unit To SQL Native Succeeded")
}

func TestJourneyUnitFromSQLNative(t *testing.T) {
	db, err := ti.CreateDatabase("gigo-dev-tidb", "4000", "mysql", "gigo-dev", "gigo-dev", "gigo_dev_test")
	if err != nil {
		t.Error("\nInitialize Database Failed\n    Error: ", err)
	}

	defer db.DB.Exec("delete from journey_unit")

	journey, err := CreateJourneyUnit(
		1,
		"Journey",
		"fullstack",
		[]ProgrammingLanguage{Go, Python, JavaScript},
		"this is a test description",
		69420,
		&DefaultWorkspaceSettings,
	)
	if err != nil {
		t.Error("\nCreate Journey Unit Failed\n    Error: ", err)
		return
	}

	statements := journey.ToSQLNative()

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

	rows, err := db.DB.Query("select * from journey_unit where _id = ?", journey.ID)
	if err != nil {
		t.Error("\nCreate Journey Unit Failed\n    Error: ", err)
		return
	}

	j, err := JourneyUnitFromSQLNative(rows)
	if err != nil {
		t.Error("\nJourney Unit From SQL Native Failed\n    Error: ", err)
		return
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

	if j.UnitFocus != "fullstack" {
		t.Error("\nJourney Unit From SQL Native Failed\n    Error: wrong unit focus")
		return
	}

	if !reflect.DeepEqual(journey.LanguageList, []ProgrammingLanguage{Go, Python, JavaScript}) {
		t.Error("\nJourney Unit From SQL Native Failed\n    Error: wrong language list")
		return
	}

	if j.Description != "this is a test description" {
		t.Error("\nJourney Unit From SQL Native Failed\n    Error: wrong description")
		return
	}

	if j.WorkspaceConfig != 69420 {
		t.Error("\nJourney Unit From SQL Native Failed\n    Error: wrong workspace config")
		return
	}

	t.Log("\nJourney Unit From SQL Native Succeeded")
}
