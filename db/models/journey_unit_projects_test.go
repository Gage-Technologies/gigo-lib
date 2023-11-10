package models

import (
	"fmt"
	ti "github.com/gage-technologies/gigo-lib/db"
	"reflect"
	"testing"
	"time"
)

func TestCreateJourneyUnitProjects(t *testing.T) {
	estimatedCompletion := time.Minute * 8
	tags := []string{"test-tag", "test-tag-2"}
	dependencies := []int64{3, 4}

	journey, err := CreateJourneyUnitProjects(
		1,
		69420,
		"/codebase",
		1,
		"learn frontend lesson 1",
		"a test project",
		Go,
		tags,
		dependencies,
		&estimatedCompletion,
	)
	if err != nil {
		t.Errorf("failed to create journey unit projects: %v", err)
		return
	}

	if journey == nil {
		t.Errorf("failed to create journey unit projects")
		return
	}

	if journey.ID != 1 {
		t.Errorf("failed to create journey unit projects with id: %d", journey.ID)
		return
	}

	if journey.UnitID != 69420 {
		t.Errorf("failed to create journey unit projects with unit id: %d", journey.UnitID)
		return
	}

	if journey.Title != "learn frontend lesson 1" {
		t.Errorf("failed to create journey unit projects with title: %s", journey.Title)
		return
	}

	if journey.Description != "a test project" {
		t.Errorf("failed to create journey unit projects with description: %s", journey.Description)
		return
	}

	if !reflect.DeepEqual(journey.ProjectLanguage, Go) {
		t.Errorf("failed to create journey unit projects with project language: %s", journey.ProjectLanguage)
		return
	}

	if journey.EstimatedTutorialTime != &estimatedCompletion {
		t.Errorf("failed to create journey unit projects with estimated duration: %s", journey.EstimatedTutorialTime)
		return
	}

	if len(journey.Tags) != 2 {
		t.Errorf("failed to create journey unit projects with tags: %d", len(journey.Tags))
		return
	}

	if len(journey.Dependencies) != 2 {
		t.Errorf("failed to create journey unit projects with dependencies: %d", len(journey.Dependencies))
		return
	}

	t.Log("\ncreated journey unit projects")

}

func TestJourneyUnitProjects_ToSQLNative(t *testing.T) {
	estimatedCompletion := time.Minute * 8
	tags := []string{"test-tag", "test-tag-2"}
	dependencies := []int64{3, 4}

	journey, err := CreateJourneyUnitProjects(
		1,
		69420,
		"/codebase",
		1,
		"learn frontend",
		"a test project",
		Go,
		tags,
		dependencies,
		&estimatedCompletion,
	)
	if err != nil {
		t.Errorf("failed to create journey unit projects: %v", err)
	}

	if journey == nil {
		t.Errorf("failed to create journey unit projects")
	}

	statement := journey.ToSQLNative()

	if statement[len(statement)-1].Statement != "insert ignore into journey_unit_projects (_id, unit_id, completions, working_directory, title, description, project_language, estimated_tutorial_time) values (?,?,?,?,?,?,?,?);" {
		t.Errorf("\nCreate Journey Unit Projects Failed\n    Error: wrong statement: \n%v", statement[len(statement)-1].Statement)
		return
	}

	if len(statement[0].Values) != 2 {
		t.Error("\nCreate Journey Unit Projects Failed\n    Error: wrong number of values")
		return
	}

	t.Log("\ncreated journey unit projects")
}

func TestJourneyUnitProjectsFromSQLNative(t *testing.T) {
	db, err := ti.CreateDatabase("gigo-dev-tidb", "4000", "mysql", "gigo-dev", "gigo-dev", "gigo_dev_test")
	if err != nil {
		t.Error("\nInitialize Database Failed\n    Error: ", err)
	}

	defer db.DB.Exec("delete from journey_unit_projects")

	estimatedCompletion := time.Minute * 8
	tags := []string{"test-tag", "test-tag-2"}
	dependencies := []int64{3, 4}

	journey, err := CreateJourneyUnitProjects(
		1,
		69420,
		"/codebase",
		1,
		"learn frontend",
		"a test project",
		Go,
		tags,
		dependencies,
		&estimatedCompletion,
	)
	if err != nil {
		t.Errorf("failed to create journey unit projects: %v", err)
	}

	statements := journey.ToSQLNative()

	fmt.Println(statements[len(statements)-1].Statement, statements[len(statements)-1].Values)

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

	rows, err := db.DB.Query("select * from journey_unit_projects where _id = ?", journey.ID)
	if err != nil {
		t.Error("\nCreate Journey Unit Projects Failed\n    Error: ", err)
		return
	}

	if !rows.Next() {
		t.Error("\nCreate Journey Unit Projects Failed\n    Error: no rows found")
		return
	}

	j, err := JourneyUnitProjectsFromSQLNative(db, rows)
	if err != nil {
		t.Error("\nCreate Journey Unit Projects Failed\n    Error: ", err)
		return
	}

	if j == nil {
		t.Error("\nCreate Journey Unit Projects Failed\n    Error: creation returned nil")
		return
	}

	if j.ID != 1 {
		t.Errorf("\nCreate Journey Unit Projects Failed\n    Error: wrong id, got: %d", j.ID)
		return
	}

	if j.UnitID != 69420 {
		t.Errorf("\nCreate Journey Unit Projects Failed\n    Error: wrong unit id, got: %d", j.UnitID)
		return
	}

	if j.Title != "learn frontend" {
		t.Errorf("\nCreate Journey Unit Projects Failed\n    Error: wrong title, got: %s", j.Title)
		return
	}

	if j.Description != "a test project" {
		t.Errorf("\nCreate Journey Unit Projects Failed\n    Error: wrong description, got: %s", j.Description)
		return
	}

	if !reflect.DeepEqual(j.ProjectLanguage, Go) {
		t.Errorf("\nCreate Journey Unit Projects Failed\n    Error: wrong project language, got: %s", j.ProjectLanguage)
		return
	}

	if *j.EstimatedTutorialTime != estimatedCompletion {
		t.Errorf("\nCreate Journey Unit Projects Failed\n    Error: wrong estimated time completion, got: %s", j.EstimatedTutorialTime)
		return
	}

	if len(j.Tags) != 2 {
		t.Errorf("\nCreate Journey Unit Projects Failed\n    Error: wrong number of tags, got: %d", len(j.Tags))
		return
	}

	if len(j.Dependencies) != 2 {
		t.Errorf("\nCreate Journey Unit Projects Failed\n    Error: wrong number of dependencies, got: %d", len(j.Dependencies))
		return
	}

	t.Log("\ncreated journey unit projects")
}
