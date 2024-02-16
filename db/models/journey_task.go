package models

import (
	"database/sql"
	"fmt"

	"github.com/kisielk/sqlstruct"
)

type JourneyTask struct {
	ID             int64               `json:"_id" sql:"_id"`
	Name           string              `json:"name" sql:"name"`
	Description    string              `json:"description" sql:"description"`
	JourneyUnitID  int64               `json:"journey_id" sql:"journey_id"`
	NodeAbove      *int64              `json:"node_above" sql:"node_above"`
	NodeBelow      *int64              `json:"node_below" sql:"node_below"`
	CodeSourceId   int64               `json:"code_source_id" sql:"code_source_id"`
	CodeSourceType CodeSource          `json:"code_source" sql:"code_source"`
	Lang           ProgrammingLanguage `json:"lang" sql:"lang"`
	Published      bool                `json:"published" sql:"published"`
}

type JourneyTaskSQL struct {
	ID             int64               `json:"_id" sql:"_id"`
	Name           string              `json:"name" sql:"name"`
	Description    string              `json:"description" sql:"description"`
	JourneyUnitID  int64               `json:"journey_id" sql:"journey_id"`
	NodeAbove      *int64              `json:"node_above" sql:"node_above"`
	NodeBelow      *int64              `json:"node_below" sql:"node_below"`
	CodeSourceId   int64               `json:"code_source_id" sql:"code_source_id"`
	CodeSourceType CodeSource          `json:"code_source" sql:"code_source"`
	Lang           ProgrammingLanguage `json:"lang" sql:"lang"`
	Published      bool                `json:"published" sql:"published"`
}

type JourneyTaskFrontend struct {
	ID             string              `json:"_id" sql:"_id"`
	Name           string              `json:"name" sql:"name"`
	Description    string              `json:"description" sql:"description"`
	JourneyUnitID  string              `json:"journey_id" sql:"journey_id"`
	NodeAbove      *string             `json:"node_above" sql:"node_above"`
	NodeBelow      *string             `json:"node_below" sql:"node_below"`
	CodeSourceId   string              `json:"code_source_id" sql:"code_source_id"`
	CodeSourceType CodeSource          `json:"code_source" sql:"code_source"`
	Lang           ProgrammingLanguage `json:"lang" sql:"lang"`
	Published      bool                `json:"published" sql:"published"`
}

func CreateJourneyTask(id int64, name string, description string, journeyUnitID int64, nodeAbove *int64,
	nodeBelow *int64, codeSourceId int64, codeSourceType CodeSource,
	lang ProgrammingLanguage, published bool) (*JourneyTask, error) {
	return &JourneyTask{
		ID:             id,
		Name:           name,
		Description:    description,
		JourneyUnitID:  journeyUnitID,
		NodeAbove:      nodeAbove,
		NodeBelow:      nodeBelow,
		CodeSourceId:   codeSourceId,
		CodeSourceType: codeSourceType,
		Lang:           lang,
		Published:      published,
	}, nil
}

func JourneyTaskFromSQLNative(rows *sql.Rows) (*JourneyTask, error) {
	JourneyTaskSQL := new(JourneyTaskSQL)
	err := sqlstruct.Scan(JourneyTaskSQL, rows)
	if err != nil {
		return nil, fmt.Errorf("error scanning JourneyTask info in first scan: %v", err)
	}

	return &JourneyTask{
		ID:             JourneyTaskSQL.ID,
		Name:           JourneyTaskSQL.Name,
		Description:    JourneyTaskSQL.Description,
		JourneyUnitID:  JourneyTaskSQL.JourneyUnitID,
		NodeAbove:      JourneyTaskSQL.NodeAbove,
		NodeBelow:      JourneyTaskSQL.NodeBelow,
		CodeSourceId:   JourneyTaskSQL.CodeSourceId,
		CodeSourceType: JourneyTaskSQL.CodeSourceType,
		Lang:           JourneyTaskSQL.Lang,
		Published:      JourneyTaskSQL.Published,
	}, nil
}

func (b *JourneyTask) ToFrontend() *JourneyTaskFrontend {

	var nodeAbove *string
	if b.NodeAbove != nil {
		nodeStr := fmt.Sprintf("%v", b.NodeAbove)
		nodeAbove = &nodeStr
	}

	var nodeBelow *string

	if b.NodeBelow != nil {
		nodeStr := fmt.Sprintf("%v", b.NodeBelow)
		nodeBelow = &nodeStr
	}

	return &JourneyTaskFrontend{
		ID:             fmt.Sprintf("%d", b.ID),
		Name:           b.Name,
		Description:    b.Description,
		JourneyUnitID:  fmt.Sprintf("%v", b.JourneyUnitID),
		NodeAbove:      nodeAbove,
		NodeBelow:      nodeBelow,
		CodeSourceId:   fmt.Sprintf("%v", b.CodeSourceId),
		CodeSourceType: b.CodeSourceType,
		Lang:           b.Lang,
		Published:      b.Published,
	}
}

func (b *JourneyTask) ToSQLNative() (*SQLInsertStatement, error) {

	sqlStatements := &SQLInsertStatement{
		Statement: "insert ignore into journey_tasks(_id, name, description, journey_unit_id, node_above, node_below, code_source_id, code_source_type, lang, published) values(?,?,?,?,?,?,?,?,?,?);",
		Values:    []interface{}{b.ID, b.Name, b.Description, b.JourneyUnitID, b.NodeAbove, b.NodeBelow, b.CodeSourceId, b.CodeSourceType, b.Lang, b.Published},
	}

	return sqlStatements, nil
}
