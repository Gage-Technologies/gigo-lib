package models

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/kisielk/sqlstruct"
)

type ByteAttempts struct {
	ID            int64  `json:"_id" sql:"_id"`
	ByteID        int64  `json:"byte_id" sql:"byte_id"`
	AuthorID      int64  `json:"author_id" sql:"author_id"`
	ContentEasy   string `json:"content_easy" sql:"content_easy"`
	ContentMedium string `json:"content_medium" sql:"content_medium"`
	ContentHard   string `json:"content_hard" sql:"content_hard"`
	Modified      bool   `json:"modified" sql:"modified"`
}

type ByteAttemptsSQL struct {
	ID            int64  `json:"_id" sql:"_id"`
	ByteID        int64  `json:"byte_id" sql:"byte_id"`
	AuthorID      int64  `json:"author_id" sql:"author_id"`
	ContentEasy   string `json:"content_easy" sql:"content_easy"`
	ContentMedium string `json:"content_medium" sql:"content_medium"`
	ContentHard   string `json:"content_hard" sql:"content_hard"`
	Modified      bool   `json:"modified" sql:"modified"`
}

type ByteAttemptsFrontend struct {
	ID            string `json:"_id" sql:"_id"`
	ByteID        string `json:"byte_id" sql:"byte_id"`
	AuthorID      string `json:"author_id" sql:"author_id"`
	ContentEasy   string `json:"content_easy" sql:"content_easy"`
	ContentMedium string `json:"content_medium" sql:"content_medium"`
	ContentHard   string `json:"content_hard" sql:"content_hard"`
	Modified      bool   `json:"modified" sql:"modified"`
}

func CreateByteAttempts(id int64, byteID int64, authorID int64, easyContent string, mediumContent string, hardContent string) (*ByteAttempts, error) {
	return &ByteAttempts{
		ID:            id,
		ByteID:        byteID,
		AuthorID:      authorID,
		ContentEasy:   easyContent,
		ContentMedium: mediumContent,
		ContentHard:   hardContent,
	}, nil
}

func ByteAttemptsFromSQLNative(rows *sql.Rows) (*ByteAttempts, error) {
	byteAttemptsSQL := new(ByteAttemptsSQL)
	err := sqlstruct.Scan(byteAttemptsSQL, rows)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Error scanning ByteAttempts info in first scan: %v", err))
	}

	return &ByteAttempts{
		ID:            byteAttemptsSQL.ID,
		ByteID:        byteAttemptsSQL.ByteID,
		AuthorID:      byteAttemptsSQL.AuthorID,
		ContentEasy:   byteAttemptsSQL.ContentEasy,
		ContentMedium: byteAttemptsSQL.ContentMedium,
		ContentHard:   byteAttemptsSQL.ContentHard,
		Modified:      byteAttemptsSQL.Modified,
	}, nil
}

func (b *ByteAttempts) ToFrontend() *ByteAttemptsFrontend {
	return &ByteAttemptsFrontend{
		ID:            fmt.Sprintf("%d", b.ID),
		ByteID:        fmt.Sprintf("%d", b.ByteID),
		AuthorID:      fmt.Sprintf("%d", b.AuthorID),
		ContentEasy:   b.ContentEasy,
		ContentMedium: b.ContentMedium,
		ContentHard:   b.ContentHard,
		Modified:      b.Modified,
	}
}

func (b *ByteAttempts) ToSQLNative() ([]*SQLInsertStatement, error) {
	sqlStatements := make([]*SQLInsertStatement, 0)

	sqlStatements = append(sqlStatements, &SQLInsertStatement{
		Statement: "insert ignore into byte_attempts(_id, byte_id, author_id, content_easy, content_medium, content_hard, modified) values(?,?,?,?,?,?,?);",
		Values:    []interface{}{b.ID, b.ByteID, b.AuthorID, b.ContentEasy, b.ContentMedium, b.ContentHard, b.Modified},
	})

	return sqlStatements, nil
}
