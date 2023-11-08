package models

import (
	"database/sql"
	"fmt"
	"github.com/kisielk/sqlstruct"
)

type JourneyDependencies struct {
	ProjectID int64 `json:"project_id" sql:"project_id"`
	DepID     int64 `json:"dependency_id" sql:"dependency_id"`
}

type JourneyDependenciesSQL struct {
	ProjectID int64 `json:"project_id" sql:"project_id"`
	DepID     int64 `json:"dependency_id" sql:"dependency_id"`
}

type JourneyDependenciesFrontend struct {
	ProjectID string `json:"project_id" sql:"project_id"`
	DepID     string `json:"dependency_id" sql:"dependency_id"`
}

func CreateJourneyDependencies(id int64, depID int64) *JourneyDependencies {
	return &JourneyDependencies{
		ProjectID: id,
		DepID:     depID,
	}
}

func JourneyDependenciesFromSQLNative(rows *sql.Rows) (*JourneyDependencies, error) {
	// create new tag object to load into
	tagSQL := new(JourneyDependencies)

	// scan row into tag object
	err := sqlstruct.Scan(tagSQL, rows)
	if err != nil {
		return nil, err
	}

	return &JourneyDependencies{
		ProjectID: tagSQL.ProjectID,
		DepID:     tagSQL.DepID,
	}, nil
}

func (t *JourneyDependencies) ToFrontend() *JourneyDependenciesFrontend {
	return &JourneyDependenciesFrontend{
		ProjectID: fmt.Sprintf("%d", t.ProjectID),
		DepID:     fmt.Sprintf("%d", t.DepID),
	}
}

func (t *JourneyDependencies) ToSQLNative() []*SQLInsertStatement {
	return []*SQLInsertStatement{
		{
			Statement: "insert ignore into journey_unit_project_dependencies(project_id, dependency_id) values (?, ?)",
			Values: []interface{}{
				t.ProjectID, t.DepID,
			},
		},
	}
}
