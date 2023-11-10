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

type JourneyUnitAttempt struct {
	ID                      int64                   `json:"_id" sql:"_id"`
	UserID                  int64                   `json:"user_id" sql:"user_id"`
	ParentUnit              int64                   `json:"parent_unit" sql:"parent_unit"`
	Title                   string                  `json:"title" sql:"title"`
	UnitFocus               UnitFocus               `json:"unit_focus" sql:"unit_focus"`
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
}

type JourneyUnitAttemptSQL struct {
	ID                      int64          `json:"_id" sql:"_id"`
	UserID                  int64          `json:"user_id" sql:"user_id"`
	ParentUnit              int64          `json:"parent_unit" sql:"parent_unit"`
	Title                   string         `json:"title" sql:"title"`
	UnitFocus               string         `json:"unit_focus" sql:"unit_focus"`
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
}

type JourneyUnitAttemptFrontend struct {
	ID                    string         `json:"_id" sql:"_id"`
	Title                 string         `json:"title" sql:"title"`
	UnitFocus             UnitFocus      `json:"unit_focus" sql:"unit_focus"`
	LanguageList          []string       `json:"language_list" sql:"language_list"`
	Description           string         `json:"description" sql:"description"`
	Tags                  []string       `json:"tags" sql:"tags"`
	EstimatedTutorialTime *time.Duration `json:"estimated_tutorial_time,omitempty" sql:"estimated_tutorial_time"`
}

func CreateJourneyUnitAttempt(id int64, userID int64, parentUnit int64, title string, unitFocus UnitFocus, languageList []ProgrammingLanguage, description string, repoID int64, createdAt time.Time, updatedAt time.Time, tags []string, tier TierType, workspaceSettings *WorkspaceSettings, workspaceConfig int64, workspaceConfigRevision int, estimatedTutorialTime *time.Duration) (*JourneyUnitAttempt, error) {
	jTags := make([]*JourneyTags, 0)

	for _, t := range tags {
		jTags = append(jTags, CreateJourneyTags(id, t, JourneyUnitType))
	}

	jLanguages := make([]*JourneyUnitLanguages, 0)

	for _, l := range languageList {
		jLanguages = append(jLanguages, CreateJourneyUnitLanguages(id, l, true))
	}

	return &JourneyUnitAttempt{
		ID:                      id,
		UserID:                  userID,
		ParentUnit:              parentUnit,
		Title:                   title,
		UnitFocus:               unitFocus,
		LanguageList:            jLanguages,
		Description:             description,
		RepoID:                  repoID,
		CreatedAt:               createdAt,
		UpdatedAt:               updatedAt,
		Tags:                    jTags,
		Tier:                    tier,
		WorkspaceSettings:       workspaceSettings,
		WorkspaceConfig:         workspaceConfig,
		WorkspaceConfigRevision: workspaceConfigRevision,
		EstimatedTutorialTime:   estimatedTutorialTime,
	}, nil
}

func JourneyUnitAttemptFromSQLNative(db *ti.Database, rows *sql.Rows) (*JourneyUnitAttempt, error) {
	var journeyUnitSQL JourneyUnitAttemptSQL

	err := sqlstruct.Scan(&journeyUnitSQL, rows)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Error scanning journey unit in first scan: %v", err))
	}

	ctx, span := otel.Tracer("gigo-core").Start(context.Background(), "gigo-lib")
	defer span.End()

	callerName := "JourneyUnitFromSQLNative"

	// query tag link table to get tab ids
	tagRows, err := db.QueryContext(ctx, &span, &callerName, "select * from journey_tags where journey_id = ? and type = ?", journeyUnitSQL.ID, JourneyUnitAttemptType)
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
	languageRows, err := db.QueryContext(ctx, &span, &callerName, "select * from journey_unit_languages where unit_id = ? and is_attempt = ?", journeyUnitSQL.ID, true)
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

	return &JourneyUnitAttempt{
		ID:                      journeyUnitSQL.ID,
		UserID:                  journeyUnitSQL.UserID,
		ParentUnit:              journeyUnitSQL.ParentUnit,
		Title:                   journeyUnitSQL.Title,
		UnitFocus:               UnitFocusFromString(journeyUnitSQL.UnitFocus),
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
	}, nil
}

func (i *JourneyUnitAttempt) ToFrontend() *JourneyUnitAttemptFrontend {
	languages := make([]string, 0)
	tags := make([]string, 0)

	for _, l := range i.LanguageList {
		languages = append(languages, l.Value)
	}

	for _, t := range i.Tags {
		tags = append(tags, t.Value)
	}

	return &JourneyUnitAttemptFrontend{
		ID:                    fmt.Sprintf("%d", i.ID),
		Title:                 i.Title,
		UnitFocus:             i.UnitFocus,
		LanguageList:          languages,
		Description:           i.Description,
		Tags:                  tags,
		EstimatedTutorialTime: i.EstimatedTutorialTime,
	}
}

func (i *JourneyUnitAttempt) ToSQLNative() ([]*SQLInsertStatement, error) {
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
		Statement: "insert ignore into journey_unit_attempts (_id, title, user_id, parent_unit, unit_focus, description, repo_id, created_at, updated_at, challenge_cost, completions, attempts, tier, embedded, workspace_config, workspace_config_revision, workspace_settings, estimated_tutorial_time) values (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?);",
		Values: []interface{}{i.ID, i.Title, i.UserID, i.ParentUnit, i.UnitFocus.String(), i.Description, i.RepoID, i.CreatedAt,
			i.UpdatedAt, i.ChallengeCost, i.Completions, i.Attempts, i.Tier, i.Embedded, i.WorkspaceConfig,
			i.WorkspaceConfigRevision, buf, i.EstimatedTutorialTime},
	})

	return sqlStatements, nil
}
