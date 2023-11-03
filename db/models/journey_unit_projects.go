package models

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	ti "github.com/gage-technologies/gigo-lib/db"
	"github.com/kisielk/sqlstruct"
	"go.opentelemetry.io/otel"
	"time"
)

type JourneyUnitProjects struct {
	ID                      int64               `json:"_id" sql:"_id"`
	UnitID                  int64               `json:"unit_id" sql:"unit_id"`
	RepoID                  int64               `json:"repo_id" sql:"repo_id"`
	Title                   string              `json:"title" sql:"title"`
	Description             string              `json:"description" sql:"description"`
	ProjectLanguage         ProgrammingLanguage `json:"project_language" sql:"project_language"`
	EstimatedTimeCompletion *time.Duration      `json:"estimated_time_completion" sql:"estimated_time_completion"`
	Tags                    []int64             `json:"tags" sql:"tags"`
	Dependencies            []int64             `json:"dependencies" sql:"dependencies"`
}

type JourneyUnitProjectsSQL struct {
	ID                      int64               `json:"_id" sql:"_id"`
	UnitID                  int64               `json:"unit_id" sql:"unit_id"`
	RepoID                  int64               `json:"repo_id" sql:"repo_id"`
	Title                   string              `json:"title" sql:"title"`
	Description             string              `json:"description" sql:"description"`
	ProjectLanguage         ProgrammingLanguage `json:"project_language" sql:"project_language"`
	EstimatedTimeCompletion *time.Duration      `json:"estimated_time_completion" sql:"estimated_time_completion"`
	Dependencies            []int64             `json:"dependencies" sql:"dependencies"`
}

type JourneyUnitProjectsFrontend struct {
	ID                      string              `json:"_id" sql:"_id"`
	UnitID                  string              `json:"unit_id" sql:"unit_id"`
	RepoID                  string              `json:"repo_id" sql:"repo_id"`
	Title                   string              `json:"title" sql:"title"`
	Description             string              `json:"description" sql:"description"`
	ProjectLanguage         ProgrammingLanguage `json:"project_language" sql:"project_language"`
	EstimatedTimeCompletion string              `json:"estimated_time_completion" sql:"estimated_time_completion"`
	Tags                    []string            `json:"tags" sql:"tags"`
	Dependencies            []string            `json:"dependencies" sql:"dependencies"`
}

func CreateJourneyUnitProjects(id int64, unitID int64, repoID int64, title string, description string, projectLanguage ProgrammingLanguage, estimatedTimeCompletion *time.Duration, tags []int64, dependencies []int64) (*JourneyUnitProjects, error) {
	return &JourneyUnitProjects{
		ID:                      id,
		UnitID:                  unitID,
		RepoID:                  repoID,
		Title:                   title,
		Description:             description,
		ProjectLanguage:         projectLanguage,
		EstimatedTimeCompletion: estimatedTimeCompletion,
		Tags:                    tags,
		Dependencies:            dependencies,
	}, nil
}

func JourneyUnitProjectsFromSQLNative(db *ti.Database, rows *sql.Rows) (*JourneyUnitProjects, error) {
	var journeyUnitProjects JourneyUnitProjectsSQL

	err := sqlstruct.Scan(&journeyUnitProjects, rows)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Error scanning journey unit projects in first scan: %v", err))
	}

	ctx, span := otel.Tracer("gigo-core").Start(context.Background(), "gigo-lib")
	defer span.End()

	callerName := "JourneyUnitProjectsFromSQLNative"

	// query tag link table to get tab ids
	tagRows, err := db.QueryContext(ctx, &span, &callerName, "select tag_id from post_tags where post_id = ?", journeyUnitProjects.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to query tag link table for tag ids: %v", err)
	}

	// defer closure of tag rows
	defer tagRows.Close()

	// create slice to hold tag ids
	tags := make([]int64, 0)

	// iterate cursor scanning tag ids and saving the to the slice created above
	for tagRows.Next() {
		var tag int64
		err = tagRows.Scan(&tag)
		if err != nil {
			return nil, fmt.Errorf("failed to scan tag id from link tbale cursor: %v", err)
		}
		tags = append(tags, tag)
	}

	return &JourneyUnitProjects{
		ID:                      journeyUnitProjects.ID,
		UnitID:                  journeyUnitProjects.UnitID,
		RepoID:                  journeyUnitProjects.RepoID,
		Title:                   journeyUnitProjects.Title,
		Description:             journeyUnitProjects.Description,
		ProjectLanguage:         journeyUnitProjects.ProjectLanguage,
		EstimatedTimeCompletion: journeyUnitProjects.EstimatedTimeCompletion,
		Tags:                    tags,
		Dependencies:            journeyUnitProjects.Dependencies,
	}, nil
}

func (i *JourneyUnitProjects) ToFrontend() *JourneyUnitProjectsFrontend {
	// create slice to hold tag ids in string form
	tags := make([]string, 0)

	// iterate tag ids formatting them to string format and saving them to the above slice
	for b := range i.Tags {
		tags = append(tags, fmt.Sprintf("%d", b))
	}

	// create slice to hold dependency ids in string form
	dependencies := make([]string, 0)

	// iterate tag ids formatting them to strong format and saving them to the above slice
	for d := range i.Dependencies {
		dependencies = append(dependencies, fmt.Sprintf("%d", d))
	}

	return &JourneyUnitProjectsFrontend{
		ID:                      fmt.Sprintf("%d", i.ID),
		UnitID:                  fmt.Sprintf("%d", i.UnitID),
		RepoID:                  fmt.Sprintf("%d", i.RepoID),
		Title:                   i.Title,
		Description:             i.Description,
		ProjectLanguage:         i.ProjectLanguage,
		EstimatedTimeCompletion: i.EstimatedTimeCompletion.String(),
		Tags:                    tags,
		Dependencies:            dependencies,
	}
}

func (i *JourneyUnitProjects) ToSQLNative() []*SQLInsertStatement {
	sqlStatements := make([]*SQLInsertStatement, 0)

	sqlStatements = append(sqlStatements, &SQLInsertStatement{
		Statement: "insert ignore into journey_unit_projects (_id, unit_id, repo_id, title, description, project_language, estimated_time_completion, tags, dependencies) values (?,?,?,?,?,?,?,?,?,?);",
		Values:    []interface{}{i.ID, i.UnitID, i.RepoID, i.Title, i.Description, i.ProjectLanguage, i.EstimatedTimeCompletion, i.Tags, i.Dependencies},
	})

	return sqlStatements
}
