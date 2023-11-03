package models

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/kisielk/sqlstruct"
)

type JourneyUnit struct {
	ID                      int64                 `json:"_id" sql:"_id"`
	Title                   string                `json:"title" sql:"title"`
	UnitFocus               string                `json:"unit_focus" sql:"unit_focus"`
	LanguageList            []ProgrammingLanguage `json:"language_list" sql:"language_list"`
	Description             string                `json:"description" sql:"description"`
	WorkspaceConfig         int64                 `json:"workspace_config" sql:"workspace_config"`
	WorkspaceConfigRevision int                   `json:"workspace_config_revision" sql:"workspace_config_revision"`
	WorkspaceSettings       *WorkspaceSettings    `json:"workspace_settings" sql:"workspace_settings"`
}

type JourneyUnitSQL struct {
	ID                      int64                 `json:"_id" sql:"_id"`
	Title                   string                `json:"title" sql:"title"`
	UnitFocus               string                `json:"unit_focus" sql:"unit_focus"`
	LanguageList            []ProgrammingLanguage `json:"language_list" sql:"language_list"`
	Description             string                `json:"description" sql:"description"`
	WorkspaceConfig         int64                 `json:"workspace_config" sql:"workspace_config"`
	WorkspaceConfigRevision int                   `json:"workspace_config_revision" sql:"workspace_config_revision"`
	WorkspaceSettings       *WorkspaceSettings    `json:"workspace_settings" sql:"workspace_settings"`
}

type JourneyUnitFrontend struct {
	ID           string                `json:"_id" sql:"_id"`
	Title        string                `json:"title" sql:"title"`
	UnitFocus    string                `json:"unit_focus" sql:"unit_focus"`
	LanguageList []ProgrammingLanguage `json:"language_list" sql:"language_list"`
	Description  string                `json:"description" sql:"description"`
}

func CreateJourneyUnit(id int64, title string, unitFocus string, languageList []ProgrammingLanguage, description string, workspaceConfig int64, workspaceSettings *WorkspaceSettings) (*JourneyUnit, error) {
	return &JourneyUnit{
		ID:                id,
		Title:             title,
		UnitFocus:         unitFocus,
		LanguageList:      languageList,
		Description:       description,
		WorkspaceConfig:   workspaceConfig,
		WorkspaceSettings: workspaceSettings,
	}, nil
}

func JourneyUnitFromSQLNative(rows *sql.Rows) (*JourneyUnit, error) {
	var journeyUnitSQL JourneyUnitSQL

	err := sqlstruct.Scan(&journeyUnitSQL, rows)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Error scanning journey unit in first scan: %v", err))
	}

	return &JourneyUnit{
		ID:                journeyUnitSQL.ID,
		Title:             journeyUnitSQL.Title,
		UnitFocus:         journeyUnitSQL.UnitFocus,
		LanguageList:      journeyUnitSQL.LanguageList,
		Description:       journeyUnitSQL.Description,
		WorkspaceConfig:   journeyUnitSQL.WorkspaceConfig,
		WorkspaceSettings: journeyUnitSQL.WorkspaceSettings,
	}, nil
}

func (i *JourneyUnit) ToFrontend() *JourneyUnitFrontend {
	return &JourneyUnitFrontend{
		ID:           fmt.Sprintf("%d", i.ID),
		Title:        i.Title,
		UnitFocus:    i.UnitFocus,
		LanguageList: i.LanguageList,
		Description:  i.Description,
	}
}

func (i *JourneyUnit) ToSQLNative() []*SQLInsertStatement {
	sqlStatements := make([]*SQLInsertStatement, 0)

	sqlStatements = append(sqlStatements, &SQLInsertStatement{
		Statement: "insert ignore into journey_unit (_id, title, unit_focus, language_list, description, workspace_config, workspace_settings) values (?,?,?,?,?,?,?);",
		Values:    []interface{}{i.ID, i.Title, i.UnitFocus, i.LanguageList, i.Description, i.WorkspaceConfig, i.WorkspaceSettings},
	})

	return sqlStatements
}
