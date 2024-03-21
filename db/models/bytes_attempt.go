package models

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gage-technologies/gigo-lib/types"
	"github.com/kisielk/sqlstruct"
)

type ByteAttempts struct {
	ID              int64            `json:"_id" sql:"_id"`
	ByteID          int64            `json:"byte_id" sql:"byte_id"`
	AuthorID        int64            `json:"author_id" sql:"author_id"`
	FilesEasy       []types.CodeFile `json:"files_easy" sql:"files_easy"`
	FilesMedium     []types.CodeFile `json:"files_medium" sql:"files_medium"`
	FilesHard       []types.CodeFile `json:"files_hard" sql:"files_hard"`
	Modified        bool             `json:"modified" sql:"modified"`
	CompletedEasy   bool             `json:"completed_easy" sql:"completed_easy"`
	CompletedMedium bool             `json:"completed_medium" sql:"completed_medium"`
	CompletedHard   bool             `json:"completed_hard" sql:"completed_hard"`
}

type ByteAttemptsSQL struct {
	ID              int64  `json:"_id" sql:"_id"`
	ByteID          int64  `json:"byte_id" sql:"byte_id"`
	AuthorID        int64  `json:"author_id" sql:"author_id"`
	FilesEasy       []byte `json:"files_easy" sql:"files_easy"`
	FilesMedium     []byte `json:"files_medium" sql:"files_medium"`
	FilesHard       []byte `json:"files_hard" sql:"files_hard"`
	Modified        bool   `json:"modified" sql:"modified"`
	CompletedEasy   bool   `json:"completed_easy" sql:"completed_easy"`
	CompletedMedium bool   `json:"completed_medium" sql:"completed_medium"`
	CompletedHard   bool   `json:"completed_hard" sql:"completed_hard"`
}

type ByteAttemptsFrontend struct {
	ID              string           `json:"_id" sql:"_id"`
	ByteID          string           `json:"byte_id" sql:"byte_id"`
	AuthorID        string           `json:"author_id" sql:"author_id"`
	FilesEasy       []types.CodeFile `json:"files_easy" sql:"files_easy"`
	FilesMedium     []types.CodeFile `json:"files_medium" sql:"files_medium"`
	FilesHard       []types.CodeFile `json:"files_hard" sql:"files_hard"`
	Modified        bool             `json:"modified" sql:"modified"`
	CompletedEasy   bool             `json:"completed_easy" sql:"completed_easy"`
	CompletedMedium bool             `json:"completed_medium" sql:"completed_medium"`
	CompletedHard   bool             `json:"completed_hard" sql:"completed_hard"`
}

func CreateByteAttempts(id int64, byteID int64, authorID int64, easyFiles []types.CodeFile, mediumFiles []types.CodeFile, hardFiles []types.CodeFile) (*ByteAttempts, error) {
	return &ByteAttempts{
		ID:          id,
		ByteID:      byteID,
		AuthorID:    authorID,
		FilesEasy:   easyFiles,
		FilesMedium: mediumFiles,
		FilesHard:   hardFiles,
	}, nil
}

func ByteAttemptsFromSQLNative(rows *sql.Rows) (*ByteAttempts, error) {
	byteAttemptsSQL := new(ByteAttemptsSQL)
	err := sqlstruct.Scan(byteAttemptsSQL, rows)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Error scanning ByteAttempts info in first scan: %v", err))
	}

	// unmarshall the files from byte buffers
	var contentEasy []types.CodeFile
	if len(byteAttemptsSQL.FilesEasy) > 0 {
		err = json.Unmarshal(byteAttemptsSQL.FilesEasy, &contentEasy)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("Error unmarshalling Files Easy JSON into slice of bytes: %v", err))
		}
	}
	var contentMedium []types.CodeFile
	if len(byteAttemptsSQL.FilesMedium) > 0 {
		err = json.Unmarshal(byteAttemptsSQL.FilesMedium, &contentMedium)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("Error unmarshalling Files Medium JSON into slice of bytes: %v", err))
		}
	}
	var contentHard []types.CodeFile
	if len(byteAttemptsSQL.FilesHard) > 0 {
		err = json.Unmarshal(byteAttemptsSQL.FilesHard, &contentHard)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("Error unmarshalling Files Hard JSON into slice of bytes: %v", err))
		}
	}

	return &ByteAttempts{
		ID:              byteAttemptsSQL.ID,
		ByteID:          byteAttemptsSQL.ByteID,
		AuthorID:        byteAttemptsSQL.AuthorID,
		FilesEasy:       contentEasy,
		FilesMedium:     contentMedium,
		FilesHard:       contentHard,
		Modified:        byteAttemptsSQL.Modified,
		CompletedEasy:   byteAttemptsSQL.CompletedEasy,
		CompletedMedium: byteAttemptsSQL.CompletedMedium,
		CompletedHard:   byteAttemptsSQL.CompletedHard,
	}, nil
}

func (b *ByteAttempts) ToFrontend() *ByteAttemptsFrontend {
	return &ByteAttemptsFrontend{
		ID:              fmt.Sprintf("%d", b.ID),
		ByteID:          fmt.Sprintf("%d", b.ByteID),
		AuthorID:        fmt.Sprintf("%d", b.AuthorID),
		FilesEasy:       b.FilesEasy,
		FilesMedium:     b.FilesMedium,
		FilesHard:       b.FilesHard,
		Modified:        b.Modified,
		CompletedEasy:   b.CompletedEasy,
		CompletedMedium: b.CompletedMedium,
		CompletedHard:   b.CompletedHard,
	}
}

func (b *ByteAttempts) ToSQLNative() ([]*SQLInsertStatement, error) {
	sqlStatements := make([]*SQLInsertStatement, 0)

	// marshall the files into byte buffers
	filesEasy, err := json.Marshal(b.FilesEasy)
	if err != nil {
		return nil, fmt.Errorf("error marshaling Files Easy JSON: %v", err)
	}
	filesMedium, err := json.Marshal(b.FilesMedium)
	if err != nil {
		return nil, fmt.Errorf("error marshaling Files Medium JSON: %v", err)
	}
	filesHard, err := json.Marshal(b.FilesHard)
	if err != nil {
		return nil, fmt.Errorf("error marshaling Files Hard JSON: %v", err)
	}

	sqlStatements = append(sqlStatements, &SQLInsertStatement{
		Statement: "insert ignore into byte_attempts(_id, byte_id, author_id, files_easy, files_medium, files_hard, modified, completed_easy, completed_medium, completed_hard) values(?,?,?,?,?,?,?,?,?,?);",
		Values:    []interface{}{b.ID, b.ByteID, b.AuthorID, filesEasy, filesMedium, filesHard, b.Modified, b.CompletedEasy, b.CompletedMedium, b.CompletedHard},
	})

	return sqlStatements, nil
}
