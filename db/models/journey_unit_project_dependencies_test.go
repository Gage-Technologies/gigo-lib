package models

import (
	"context"
	ti "github.com/gage-technologies/gigo-lib/db"
	"go.opentelemetry.io/otel"
	"testing"
)

func TestCreateJourneyDependencies(t *testing.T) {
	lang := CreateJourneyDependencies(
		1,
		69420,
	)

	if lang == nil {
		t.Error("\nCreate Journey Unit Languages Failed\n    Error: creation returned nil")
		return
	}

	if lang.ProjectID != 1 {
		t.Error("\nCreate Journey Unit Languages Failed\n    Error: wrong id")
		return
	}

	if lang.DepID != 69420 {
		t.Error("\nCreate Journey Unit Languages Failed\n    Error: wrong dep id")
		return
	}

	t.Log("\nCreate Journey Unit Languages Succeeded")
}

func TestJourneyDependencies_ToSQLNative(t *testing.T) {
	journey := CreateJourneyDependencies(
		1,
		69420,
	)

	statement := journey.ToSQLNative()

	if statement[len(statement)-1].Statement != "insert ignore into journey_unit_project_dependencies(project_id, dependency_id) values (?, ?)" {
		t.Error("\nJourney Unit Languages To SQL Native Failed\n    Error: wrong statement")
		return
	}

	if len(statement[len(statement)-1].Values) != 2 {
		t.Error("\nJourney Unit Languages To SQL Native Failed\n    Error: wrong number of values")
		return
	}

	t.Log("\nJourney Unit Languages To SQL Native Succeeded")
}

func TestJourneyDependenciesFromSQLNative(t *testing.T) {
	ctx, span := otel.Tracer("gigo-core").Start(context.Background(), "gigo-lib")
	defer span.End()

	callerName := "JourneyDependenciesFromSQLNativeTest"

	db, err := ti.CreateDatabase("gigo-dev-tidb", "4000", "mysql", "gigo-dev", "gigo-dev", "gigo_dev_test")
	if err != nil {
		t.Error("\nInitialize Database Failed\n    Error: ", err)
	}

	defer db.DB.Exec("delete from journey_unit_project_dependencies")

	lang := CreateJourneyDependencies(
		1,
		69420,
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

	rows, err := db.QueryContext(ctx, &span, &callerName, "select * from journey_unit_project_dependencies where project_id = ?", lang.ProjectID)
	if err != nil {
		t.Error("\nCreate Journey Unit Languages Failed\n    Error: ", err)
		return
	}

	defer rows.Close()

	var j *JourneyDependencies

	for rows.Next() {
		j, err = JourneyDependenciesFromSQLNative(rows)
		if err != nil {
			t.Error("\nJourney Unit Languages From SQL Native Failed\n    Error: ", err)
			return
		}
	}

	if j == nil {
		t.Error("\nJourney Unit Languages From SQL Native Failed\n    Error: creation returned nil")
		return
	}

	if j.ProjectID != 1 {
		t.Error("\nJourney Unit Languages From SQL Native Failed\n    Error: wrong id")
		return
	}

	if j.DepID != 69420 {
		t.Error("\nJourney Unit Languages From SQL Native Failed\n    Error: wrong dep id")
		return
	}

	t.Log("\nJourney Unit Languages From SQL Native Succeeded")
}
