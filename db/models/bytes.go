package models

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/kisielk/sqlstruct"
)

type Bytes struct {
	ID             int64  `json:"_id" sql:"_id"`
	Name           string `json:"name" sql:"name"`
	Description    string `json:"description" sql:"description"`
	OutlineContent string `json:"outline_content" sql:"outline_content"`
	DevSteps       string `json:"dev_steps" sql:"dev_steps"`
}

type BytesSQL struct {
	ID             int64  `json:"_id" sql:"_id"`
	Name           string `json:"name" sql:"name"`
	Description    string `json:"description" sql:"description"`
	OutlineContent string `json:"outline_content" sql:"outline_content"`
	DevSteps       string `json:"dev_steps" sql:"dev_steps"`
}

type BytesFrontend struct {
	ID             string `json:"_id" sql:"_id"`
	Name           string `json:"name" sql:"name"`
	Description    string `json:"description" sql:"description"`
	OutlineContent string `json:"outline_content" sql:"outline_content"`
	DevSteps       string `json:"dev_steps" sql:"dev_steps"`
}

func CreateBytes(id int64, name string, description string, outlineContent string, devSteps string) (*Bytes, error) {
	return &Bytes{
		ID:             id,
		Name:           name,
		Description:    description,
		OutlineContent: outlineContent,
		DevSteps:       devSteps,
	}, nil
}

func BytesFromSQLNative(rows *sql.Rows) (*Bytes, error) {
	bytesSQL := new(BytesSQL)
	err := sqlstruct.Scan(bytesSQL, rows)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Error scanning Bytes info in first scan: %v", err))
	}

	return &Bytes{
		ID:             bytesSQL.ID,
		Name:           bytesSQL.Name,
		Description:    bytesSQL.Description,
		OutlineContent: bytesSQL.OutlineContent,
		DevSteps:       bytesSQL.DevSteps,
	}, nil
}

func (b *Bytes) ToFrontend() *BytesFrontend {
	return &BytesFrontend{
		ID:             fmt.Sprintf("%d", b.ID),
		Name:           b.Name,
		Description:    b.Description,
		OutlineContent: b.OutlineContent,
		DevSteps:       b.DevSteps,
	}
}

func (b *Bytes) ToSQLNative() ([]*SQLInsertStatement, error) {
	sqlStatements := make([]*SQLInsertStatement, 0)

	sqlStatements = append(sqlStatements, &SQLInsertStatement{
		Statement: "insert ignore into bytes(_id, name, description, outline_content, dev_steps) values(?,?,?,?,?);",
		Values:    []interface{}{b.ID, b.Name, b.Description, b.OutlineContent, b.DevSteps},
	})

	return sqlStatements, nil
}
