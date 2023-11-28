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
	ID                    int64                  `json:"_id" sql:"_id"`
	UnitID                int64                  `json:"unit_id" sql:"unit_id"`
	WorkingDirectory      string                 `json:"working_directory" sql:"working_directory"`
	Completions           int64                  `json:"completions" sql:"completions"`
	Title                 string                 `json:"title" sql:"title"`
	Description           string                 `json:"description" sql:"description"`
	ProjectLanguage       ProgrammingLanguage    `json:"project_language" sql:"project_language"`
	Tags                  []*JourneyTags         `json:"tags" sql:"tags"`
	Dependencies          []*JourneyDependencies `json:"dependencies" sql:"dependencies"`
	EstimatedTutorialTime *time.Duration         `json:"estimated_tutorial_time,omitempty" sql:"estimated_tutorial_time"`
	Deleted               bool                   `json:"deleted,omitempty" sql:"deleted"`
}

type JourneyUnitProjectsSQL struct {
	ID                    int64               `json:"_id" sql:"_id"`
	UnitID                int64               `json:"unit_id" sql:"unit_id"`
	WorkingDirectory      string              `json:"working_directory" sql:"working_directory"`
	Completions           int64               `json:"completions" sql:"completions"`
	Title                 string              `json:"title" sql:"title"`
	Description           string              `json:"description" sql:"description"`
	ProjectLanguage       ProgrammingLanguage `json:"project_language" sql:"project_language"`
	EstimatedTutorialTime *time.Duration      `json:"estimated_tutorial_time,omitempty" sql:"estimated_tutorial_time"`
	Deleted               bool                `json:"deleted,omitempty" sql:"deleted"`
}

type JourneyUnitProjectsFrontend struct {
	ID                    string   `json:"_id" sql:"_id"`
	UnitID                string   `json:"unit_id" sql:"unit_id"`
	WorkingDirectory      string   `json:"working_directory" sql:"working_directory"`
	Title                 string   `json:"title" sql:"title"`
	Description           string   `json:"description" sql:"description"`
	ProjectLanguage       string   `json:"project_language" sql:"project_language"`
	Tags                  []string `json:"tags" sql:"tags"`
	Dependencies          []string `json:"dependencies" sql:"dependencies"`
	EstimatedTutorialTime *string  `json:"estimated_tutorial_time,omitempty" sql:"estimated_tutorial_time"`
}

func CreateJourneyUnitProjects(id int64, unitID int64, workingDirectory string, completions int64, title string, description string, projectLanguage ProgrammingLanguage, tags []string, dependencies []int64, estimatedTutorialTime *time.Duration) (*JourneyUnitProjects, error) {
	jTags := make([]*JourneyTags, 0)

	for _, t := range tags {
		jTags = append(jTags, CreateJourneyTags(id, t, JourneyUnitProjectType))
	}

	jDependencies := make([]*JourneyDependencies, 0)

	for _, l := range dependencies {
		jDependencies = append(jDependencies, CreateJourneyDependencies(id, l))
	}

	return &JourneyUnitProjects{
		ID:                    id,
		UnitID:                unitID,
		WorkingDirectory:      workingDirectory,
		Completions:           completions,
		Title:                 title,
		Description:           description,
		ProjectLanguage:       projectLanguage,
		Tags:                  jTags,
		Dependencies:          jDependencies,
		EstimatedTutorialTime: estimatedTutorialTime,
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
	tagRows, err := db.QueryContext(ctx, &span, &callerName, "select * from journey_tags where journey_id = ? and type = ?", journeyUnitProjects.ID, JourneyUnitProjectType)
	if err != nil {
		return nil, fmt.Errorf("failed to query tag link table for tag ids: %v", err)
	}

	// defer closure of tag rows
	defer tagRows.Close()

	// create slice to hold tag ids
	tags := make([]*JourneyTags, 0)

	// iterate cursor scanning tag ids and saving the to the slice created above
	for tagRows.Next() {
		t, err := JourneyTagsFromSQLNative(tagRows)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("Error scanning tag link table: %v", err))
		}
		tags = append(tags, t)
	}

	// query dependency link table to get dependency ids
	dependencyRows, err := db.QueryContext(ctx, &span, &callerName, "select * from journey_unit_project_dependencies where project_id = ?", journeyUnitProjects.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to query tag link table for tag ids: %v", err)
	}

	// defer closure of tag rows
	defer dependencyRows.Close()

	// create slice to hold tag ids
	deps := make([]*JourneyDependencies, 0)

	// iterate cursor scanning tag ids and saving the to the slice created above
	for dependencyRows.Next() {
		journeyDeps, err := JourneyDependenciesFromSQLNative(dependencyRows)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("Error scanning language link table: %v", err))
		}

		deps = append(deps, journeyDeps)

	}

	return &JourneyUnitProjects{
		ID:                    journeyUnitProjects.ID,
		UnitID:                journeyUnitProjects.UnitID,
		WorkingDirectory:      journeyUnitProjects.WorkingDirectory,
		Title:                 journeyUnitProjects.Title,
		Description:           journeyUnitProjects.Description,
		ProjectLanguage:       journeyUnitProjects.ProjectLanguage,
		Tags:                  tags,
		Dependencies:          deps,
		EstimatedTutorialTime: journeyUnitProjects.EstimatedTutorialTime,
		Deleted:               journeyUnitProjects.Deleted,
	}, nil
}

func (i *JourneyUnitProjects) ToFrontend() *JourneyUnitProjectsFrontend {
	// create slice to hold tag ids in string form
	tags := make([]string, 0)

	// iterate tag ids formatting them to string format and saving them to the above slice
	for _, b := range i.Tags {
		tags = append(tags, fmt.Sprintf("%d", b))
	}

	// create slice to hold dependency ids in string form
	dependencies := make([]string, 0)

	// iterate tag ids formatting them to strong format and saving them to the above slice
	for d := range i.Dependencies {
		dependencies = append(dependencies, fmt.Sprintf("%d", d))
	}

	var ett *string

	if i.EstimatedTutorialTime != nil {
		ettI := i.EstimatedTutorialTime.String()
		ett = &ettI
	}

	return &JourneyUnitProjectsFrontend{
		ID:                    fmt.Sprintf("%d", i.ID),
		UnitID:                fmt.Sprintf("%d", i.UnitID),
		WorkingDirectory:      i.WorkingDirectory,
		Title:                 i.Title,
		Description:           i.Description,
		ProjectLanguage:       i.ProjectLanguage.String(),
		EstimatedTutorialTime: ett,
		Tags:                  tags,
		Dependencies:          dependencies,
	}
}

func (i *JourneyUnitProjects) ToSQLNative() []*SQLInsertStatement {
	sqlStatements := make([]*SQLInsertStatement, 0)

	for _, deps := range i.Dependencies {
		s := deps.ToSQLNative()
		sqlStatements = append(sqlStatements, s...)
	}
	for _, tag := range i.Tags {
		s := tag.ToSQLNative()
		sqlStatements = append(sqlStatements, s...)
	}

	sqlStatements = append(sqlStatements, &SQLInsertStatement{
		Statement: "insert ignore into journey_unit_projects (_id, unit_id, completions, working_directory, title, description, project_language, estimated_tutorial_time, deleted) values (?,?,?,?,?,?,?,?,?);",
		Values:    []interface{}{i.ID, i.UnitID, i.Completions, i.WorkingDirectory, i.Title, i.Description, i.ProjectLanguage, i.EstimatedTutorialTime, i.Deleted},
	})

	return sqlStatements
}
