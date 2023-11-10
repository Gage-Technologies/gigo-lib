package models

import (
	ti "github.com/gage-technologies/gigo-lib/db"
	"testing"
)

func TestCreateJourneyTags(t *testing.T) {
	journey := CreateJourneyTags(69420, "test-tag", JourneyUnitType)

	if journey == nil {
		t.Error("\nCreate Journey Tags Failed\n    Error: creation returned nil")
		return
	}

	if journey.JourneyID != 69420 {
		t.Error("\nCreate Journey Tags Failed\n    Error: wrong id")
		return
	}

	if journey.Value != "test-tag" {
		t.Error("\nCreate Journey Tags Failed\n    Error: wrong user id")
		return
	}

	if journey.Type != JourneyUnitType {
		t.Error("\nCreate Journey Tags Failed\n    Error: wrong learning goal")
		return
	}

	t.Log("\nCreate Journey Tags Succeeded")

}

func TestJourneyTags_ToSQLNative(t *testing.T) {
	journey := CreateJourneyTags(69420, "test-tag", JourneyUnitType)

	statement := journey.ToSQLNative()

	if statement[0].Statement != "insert ignore into journey_tags(journey_id, value, type) values (?, ?, ?)" {
		t.Error("\nCreate Journey Tags Failed\n    Error: wrong statement")
		return
	}

	if len(statement[0].Values) != 3 {
		t.Error("\nCreate Journey Tags Failed\n    Error: wrong number of values")
	}

	t.Log("\nCreate Journey Tags Succeeded")
}

func TestJourneyTagsFromSQLNative(t *testing.T) {
	db, err := ti.CreateDatabase("gigo-dev-tidb", "4000", "mysql", "gigo-dev", "gigo-dev", "gigo_dev_test")
	if err != nil {
		t.Error("\nInitialize Database Failed\n    Error: ", err)
	}

	defer db.DB.Exec("delete from journey_tags")

	journey := CreateJourneyTags(69420, "test-tag", JourneyUnitType)
	if err != nil {
		t.Error("\nCreate Journey Tags Failed")
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

	rows, err := db.DB.Query("select * from journey_tags where journey_id = ?", journey.JourneyID)
	if err != nil {
		t.Error("\nCreate Journey Tags Failed\n    Error: ", err)
		return
	}

	if !rows.Next() {
		t.Error("\nCreate Journey Tags Failed\n    Error: no rows found")
		return
	}

	j, err := JourneyTagsFromSQLNative(rows)
	if err != nil {
		t.Error("\nJourney Tags From SQL Native Failed\n    Error: ", err)
		return
	}

	if j == nil {
		t.Error("\nJourney Tags From SQL Native Failed\n    Error: creation returned nil")
		return
	}

	if j.JourneyID != 69420 {
		t.Errorf("\nJourney Tags From SQL Native Failed\n    Error: wrong id, got: %d", j.JourneyID)
		return
	}

	if j.Value != "test-tag" {
		t.Error("\nJourney Tags From SQL Native Failed\n    Error: wrong user id")
		return
	}

	if j.Type != JourneyUnitType {
		t.Error("\nJourney Tags From SQL Native Failed\n    Error: wrong learning goal")
		return
	}

	t.Log("\nJourney Tags From SQL Native Succeeded")

}
