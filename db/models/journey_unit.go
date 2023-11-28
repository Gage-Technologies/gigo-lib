package models

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	ti "github.com/gage-technologies/gigo-lib/db"
	"github.com/kisielk/sqlstruct"
	"go.opentelemetry.io/otel"
	"time"
)

type UnitFocus int

const (
	Frontend UnitFocus = iota
	Backend
	Fullstack
)

func (u UnitFocus) String() string {
	switch u {
	case Frontend:
		return "frontend"
	case Backend:
		return "backend"
	case Fullstack:
		return "fullstack"
	default:
		return "unknown"
	}
}

func UnitFocusFromString(s string) UnitFocus {
	switch s {
	case "Frontend":
		return Frontend
	case "Backend":
		return Backend
	case "Fullstack":
		return Fullstack
	default:
		return Fullstack
	}
}

type JourneyUnit struct {
	ID                      int64                   `json:"_id" sql:"_id"`
	Title                   string                  `json:"title" sql:"title"`
	UnitFocus               UnitFocus               `json:"unit_focus" sql:"unit_focus"`
	AuthorID                int64                   `json:"author_id" sql:"author_id"`
	Visibility              PostVisibility          `json:"visibility" sql:"visibility"`
	LanguageList            []*JourneyUnitLanguages `json:"language_list" sql:"language_list"`
	Description             string                  `json:"description" sql:"description"`
	RepoID                  int64                   `json:"repo_id" sql:"repo_id"`
	CreatedAt               time.Time               `json:"created_at" sql:"created_at"`
	UpdatedAt               time.Time               `json:"updated_at" sql:"updated_at"`
	ChallengeCost           *string                 `json:"challenge_cost" sql:"challenge_cost"`
	Tags                    []*JourneyTags          `json:"tags" sql:"tags"`
	Completions             int64                   `json:"completions" sql:"completions"`
	Attempts                int64                   `json:"attempts" sql:"attempts"`
	Tier                    TierType                `json:"tier" sql:"tier"`
	Embedded                bool                    `json:"embedded" sql:"embedded"`
	WorkspaceConfig         int64                   `json:"workspace_config" sql:"workspace_config"`
	WorkspaceConfigRevision int                     `json:"workspace_config_revision" sql:"workspace_config_revision"`
	WorkspaceSettings       *WorkspaceSettings      `json:"workspace_settings" sql:"workspace_settings"`
	EstimatedTutorialTime   *time.Duration          `json:"estimated_tutorial_time,omitempty" sql:"estimated_tutorial_time"`
	Deleted                 bool                    `json:"deleted" sql:"deleted"`
}

type JourneyUnitSQL struct {
	ID                      int64          `json:"_id" sql:"_id"`
	Title                   string         `json:"title" sql:"title"`
	UnitFocus               string         `json:"unit_focus" sql:"unit_focus"`
	AuthorID                int64          `json:"author_id" sql:"author_id"`
	Visibility              int            `json:"visibility" sql:"visibility"`
	Description             string         `json:"description" sql:"description"`
	RepoID                  int64          `json:"repo_id" sql:"repo_id"`
	CreatedAt               time.Time      `json:"created_at" sql:"created_at"`
	UpdatedAt               time.Time      `json:"updated_at" sql:"updated_at"`
	ChallengeCost           *string        `json:"challenge_cost" sql:"challenge_cost"`
	Completions             int64          `json:"completions" sql:"completions"`
	Attempts                int64          `json:"attempts" sql:"attempts"`
	Tier                    TierType       `json:"tier" sql:"tier"`
	Embedded                bool           `json:"embedded" sql:"embedded"`
	WorkspaceConfig         int64          `json:"workspace_config" sql:"workspace_config"`
	WorkspaceConfigRevision int            `json:"workspace_config_revision" sql:"workspace_config_revision"`
	WorkspaceSettings       []byte         `json:"workspace_settings" sql:"workspace_settings"`
	EstimatedTutorialTime   *time.Duration `json:"estimated_tutorial_time,omitempty" sql:"estimated_tutorial_time"`
	Deleted                 bool           `json:"deleted" sql:"deleted"`
}

type JourneyUnitFrontend struct {
	ID                    string   `json:"_id" sql:"_id"`
	Title                 string   `json:"title" sql:"title"`
	UnitFocus             string   `json:"unit_focus" sql:"unit_focus"`
	Visibility            int      `json:"visibility" sql:"visibility"`
	LanguageList          []string `json:"language_list" sql:"language_list"`
	Description           string   `json:"description" sql:"description"`
	Tags                  []string `json:"tags" sql:"tags"`
	EstimatedTutorialTime *int64   `json:"estimated_tutorial_time,omitempty" sql:"estimated_tutorial_time"`
}

func CreateJourneyUnit(id int64, title string, unitFocus UnitFocus, authorID int64, visibility PostVisibility, languageList []ProgrammingLanguage, description string, repoID int64, createdAt time.Time, updatedAt time.Time, tags []string, tier TierType, workspaceSettings *WorkspaceSettings, workspaceConfig int64, estimatedTutorialTime *time.Duration) (*JourneyUnit, error) {
	jTags := make([]*JourneyTags, 0)

	for _, t := range tags {
		jTags = append(jTags, CreateJourneyTags(id, t, JourneyUnitType))
	}

	jLanguages := make([]*JourneyUnitLanguages, 0)

	for _, l := range languageList {
		jLanguages = append(jLanguages, CreateJourneyUnitLanguages(id, l, false))
	}

	return &JourneyUnit{
		ID:                    id,
		Title:                 title,
		UnitFocus:             unitFocus,
		AuthorID:              authorID,
		Visibility:            visibility,
		LanguageList:          jLanguages,
		Description:           description,
		RepoID:                repoID,
		CreatedAt:             createdAt,
		UpdatedAt:             updatedAt,
		Tags:                  jTags,
		Tier:                  tier,
		WorkspaceSettings:     workspaceSettings,
		WorkspaceConfig:       workspaceConfig,
		EstimatedTutorialTime: estimatedTutorialTime,
	}, nil
}

func JourneyUnitFromSQLNative(db *ti.Database, rows *sql.Rows) (*JourneyUnit, error) {
	var journeyUnitSQL JourneyUnitSQL

	err := sqlstruct.Scan(&journeyUnitSQL, rows)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Error scanning journey unit in first scan: %v", err))
	}

	if &journeyUnitSQL == nil {
		return nil, errors.New(fmt.Sprintf("Error scanning journey unit in first scan, journey returned nil: %v", journeyUnitSQL))
	}

	ctx, span := otel.Tracer("gigo-core").Start(context.Background(), "gigo-lib")
	defer span.End()

	callerName := "JourneyUnitFromSQLNative"

	// query tag link table to get tab ids
	tagRows, err := db.QueryContext(ctx, &span, &callerName, "select * from journey_tags where journey_id = ? and type = ?", journeyUnitSQL.ID, JourneyUnitType)
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
	languageRows, err := db.QueryContext(ctx, &span, &callerName, "select * from journey_unit_languages where unit_id = ? and is_attempt = ?", journeyUnitSQL.ID, false)
	if err != nil {
		return nil, fmt.Errorf("failed to query tag link table for tag ids: %v", err)
	}

	// defer closure of tag rows
	defer languageRows.Close()

	// create slice to hold tag ids
	languages := make([]*JourneyUnitLanguages, 0)

	// iterate cursor scanning tag ids and saving the to the slice created above
	for languageRows.Next() {
		JourneyUnitLanguages, err := JourneyUnitLanguagesFromSQLNative(languageRows)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("Error scanning language link table: %v", err))
		}

		languages = append(languages, JourneyUnitLanguages)

	}

	// create workspace settings to unmarshall into
	var workspaceSettings *WorkspaceSettings
	if journeyUnitSQL.WorkspaceSettings != nil {
		var ws WorkspaceSettings
		err = json.Unmarshal(journeyUnitSQL.WorkspaceSettings, &ws)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshall workspace settings: %v", err)
		}
		workspaceSettings = &ws
	}

	return &JourneyUnit{
		ID:                      journeyUnitSQL.ID,
		Title:                   journeyUnitSQL.Title,
		UnitFocus:               UnitFocusFromString(journeyUnitSQL.UnitFocus),
		AuthorID:                journeyUnitSQL.AuthorID,
		Visibility:              PostVisibility(journeyUnitSQL.Visibility),
		LanguageList:            languages,
		Description:             journeyUnitSQL.Description,
		RepoID:                  journeyUnitSQL.RepoID,
		CreatedAt:               journeyUnitSQL.CreatedAt,
		UpdatedAt:               journeyUnitSQL.UpdatedAt,
		Tags:                    tags,
		Tier:                    journeyUnitSQL.Tier,
		WorkspaceSettings:       workspaceSettings,
		WorkspaceConfig:         journeyUnitSQL.WorkspaceConfig,
		EstimatedTutorialTime:   journeyUnitSQL.EstimatedTutorialTime,
		Embedded:                journeyUnitSQL.Embedded,
		ChallengeCost:           journeyUnitSQL.ChallengeCost,
		Completions:             journeyUnitSQL.Completions,
		Attempts:                journeyUnitSQL.Attempts,
		WorkspaceConfigRevision: journeyUnitSQL.WorkspaceConfigRevision,
		Deleted:                 journeyUnitSQL.Deleted,
	}, nil
}

func (i *JourneyUnit) ToFrontend() *JourneyUnitFrontend {
	languages := make([]string, 0)
	tags := make([]string, 0)

	for _, l := range i.LanguageList {
		languages = append(languages, l.Value)
	}

	for _, t := range i.Tags {
		tags = append(tags, t.Value)
	}

	// conditionally load the estimated time
	var estimatedTime *int64
	if i.EstimatedTutorialTime != nil {
		millis := i.EstimatedTutorialTime.Milliseconds()
		estimatedTime = &millis
	}

	return &JourneyUnitFrontend{
		ID:                    fmt.Sprintf("%d", i.ID),
		Title:                 i.Title,
		UnitFocus:             i.UnitFocus.String(),
		Visibility:            int(i.Visibility),
		LanguageList:          languages,
		Description:           i.Description,
		Tags:                  tags,
		EstimatedTutorialTime: estimatedTime,
	}
}

func (i *JourneyUnit) ToSQLNative() ([]*SQLInsertStatement, error) {
	sqlStatements := make([]*SQLInsertStatement, 0)
	var buf []byte
	if i.WorkspaceSettings != nil {
		b, err := json.Marshal(i.WorkspaceSettings)
		if err != nil {
			return nil, fmt.Errorf("failed to marshall workspace settings: %v", err)
		}
		buf = b
	}

	for _, lang := range i.LanguageList {
		s := lang.ToSQLNative()
		sqlStatements = append(sqlStatements, s...)
	}
	for _, tag := range i.Tags {
		s := tag.ToSQLNative()
		sqlStatements = append(sqlStatements, s...)
	}

	sqlStatements = append(sqlStatements, &SQLInsertStatement{
		Statement: "insert ignore into journey_units (_id, title, unit_focus, author_id, visibility, description, repo_id, created_at, updated_at, challenge_cost, completions, attempts, tier, embedded, workspace_config, workspace_config_revision, workspace_settings, estimated_tutorial_time, deleted) values (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?);",
		Values: []interface{}{i.ID, i.Title, i.UnitFocus.String(), i.AuthorID, i.Visibility, i.Description, i.RepoID, i.CreatedAt,
			i.UpdatedAt, i.ChallengeCost, i.Completions, i.Attempts, i.Tier, i.Embedded, i.WorkspaceConfig,
			i.WorkspaceConfigRevision, buf, i.EstimatedTutorialTime, i.Deleted},
	})

	return sqlStatements, nil
}
