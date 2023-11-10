package models

import (
	"context"
	ti "github.com/gage-technologies/gigo-lib/db"
	"go.opentelemetry.io/otel"
	"testing"
)

func TestCreateJourneyUnitLanguages(t *testing.T) {
	lang := CreateJourneyUnitLanguages(
		1,
		ProgrammingLanguageFromString("Python"),
		false,
	)

	if lang == nil {
		t.Error("\nCreate Journey Unit Languages Failed\n    Error: creation returned nil")
		return
	}

	if lang.UnitID != 1 {
		t.Error("\nCreate Journey Unit Languages Failed\n    Error: wrong id")
		return
	}

	if lang.Value != "Python" {
		t.Error("\nCreate Journey Unit Languages Failed\n    Error: wrong language value id")
		return
	}

	if lang.IsAttempt != false {
		t.Error("\nCreate Journey Unit Languages Failed\n    Error: wrong parent id")
		return
	}

	t.Log("\nCreate Journey Unit Languages Succeeded")
}

func TestJourneyUnitLanguages_ToSQLNative(t *testing.T) {
	journey := CreateJourneyUnitLanguages(
		1,
		ProgrammingLanguageFromString("test-title"),
		false,
	)

	statement := journey.ToSQLNative()

	if statement[len(statement)-1].Statement != "insert ignore into journey_unit_languages(unit_id, value, is_attempt) values (?, ?, ?)" {
		t.Error("\nJourney Unit Languages To SQL Native Failed\n    Error: wrong statement")
		return
	}

	if len(statement[len(statement)-1].Values) != 3 {
		t.Error("\nJourney Unit Languages To SQL Native Failed\n    Error: wrong number of values")
		return
	}

	t.Log("\nJourney Unit Languages To SQL Native Succeeded")
}

func TestJourneyUnitLanguagesFromSQLNative(t *testing.T) {
	ctx, span := otel.Tracer("gigo-core").Start(context.Background(), "gigo-lib")
	defer span.End()

	callerName := "JourneyUnitLanguagesFromSQLNativeTest"

	db, err := ti.CreateDatabase("gigo-dev-tidb", "4000", "mysql", "gigo-dev", "gigo-dev", "gigo_dev_test")
	if err != nil {
		t.Error("\nInitialize Database Failed\n    Error: ", err)
	}

	defer db.DB.Exec("delete from journey_unit_languages")

	lang := CreateJourneyUnitLanguages(
		1,
		ProgrammingLanguageFromString("Python"),
		false,
	)
	if err != nil {
		t.Error("\nCreate Journey Unit Languages Failed\n    Error: ", err)
		return
	}

	statements := lang.ToSQLNative()

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

	rows, err := db.QueryContext(ctx, &span, &callerName, "select * from journey_unit_languages where unit_id = ?", lang.UnitID)
	if err != nil {
		t.Error("\nCreate Journey Unit Languages Failed\n    Error: ", err)
		return
	}

	defer rows.Close()

	var j *JourneyUnitLanguages

	for rows.Next() {
		j, err = JourneyUnitLanguagesFromSQLNative(rows)
		if err != nil {
			t.Error("\nJourney Unit Languages From SQL Native Failed\n    Error: ", err)
			return
		}
	}

	if j == nil {
		t.Error("\nJourney Unit Languages From SQL Native Failed\n    Error: creation returned nil")
		return
	}

	if j.UnitID != 1 {
		t.Error("\nJourney Unit Languages From SQL Native Failed\n    Error: wrong id")
		return
	}

	if j.Value != "Python" {
		t.Error("\nJourney Unit Languages From SQL Native Failed\n    Error: wrong name")
		return
	}

	if j.IsAttempt != false {
		t.Error("\nJourney Unit Languages From SQL Native Failed\n    Error: wrong unit focus")
		return
	}

	t.Log("\nJourney Unit Languages From SQL Native Succeeded")
}
