package models

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/kisielk/sqlstruct"
)

type ByteAttempts struct {
	ID       int64  `json:"_id" sql:"_id"`
	ByteID   int64  `json:"byte_id" sql:"byte_id"`
	AuthorID int64  `json:"author_id" sql:"author_id"`
	Content  string `json:"content" sql:"content"`
}

type ByteAttemptsSQL struct {
	ID       int64  `json:"_id" sql:"_id"`
	ByteID   int64  `json:"byte_id" sql:"byte_id"`
	AuthorID int64  `json:"author_id" sql:"author_id"`
	Content  string `json:"content" sql:"content"`
}

type ByteAttemptsFrontend struct {
	ID       string `json:"_id" sql:"_id"`
	ByteID   string `json:"byte_id" sql:"byte_id"`
	AuthorID string `json:"author_id" sql:"author_id"`
	Content  string `json:"content" sql:"content"`
}

func CreateByteAttempts(id int64, byteID int64, authorID int64, content string) (*ByteAttempts, error) {
	return &ByteAttempts{
		ID:       id,
		ByteID:   byteID,
		AuthorID: authorID,
		Content:  content,
	}, nil
}

func ByteAttemptsFromSQLNative(rows *sql.Rows) (*ByteAttempts, error) {
	byteAttemptsSQL := new(ByteAttemptsSQL)
	err := sqlstruct.Scan(byteAttemptsSQL, rows)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Error scanning ByteAttempts info in first scan: %v", err))
	}

	return &ByteAttempts{
		ID:       byteAttemptsSQL.ID,
		ByteID:   byteAttemptsSQL.ByteID,
		AuthorID: byteAttemptsSQL.AuthorID,
		Content:  byteAttemptsSQL.Content,
	}, nil
}

func (b *ByteAttempts) ToFrontend() *ByteAttemptsFrontend {
	return &ByteAttemptsFrontend{
		ID:       fmt.Sprintf("%d", b.ID),
		ByteID:   fmt.Sprintf("%d", b.ByteID),
		AuthorID: fmt.Sprintf("%d", b.AuthorID),
		Content:  b.Content,
	}
}

func (b *ByteAttempts) ToSQLNative() ([]*SQLInsertStatement, error) {
	sqlStatements := make([]*SQLInsertStatement, 0)

	sqlStatements = append(sqlStatements, &SQLInsertStatement{
		Statement: "insert ignore into byte_attempts(_id, byte_id, author_id, content) values(?,?,?,?);",
		Values:    []interface{}{b.ID, b.ByteID, b.AuthorID, b.Content},
	})

	return sqlStatements, nil
}
