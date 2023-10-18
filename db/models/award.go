package models

import (
	"database/sql"
	"fmt"
	"github.com/kisielk/sqlstruct"
)

type ContentType int

const (
	PostType ContentType = iota
	DiscussionType
	CommentType
	AttemptType
	ThreadCommentType
)

func (s ContentType) String() string {
	switch s {
	case PostType:
		return "Post"
	case DiscussionType:
		return "Discussion"
	case CommentType:
		return "Comment"
	case AttemptType:
		return "Attempt"
	default:
		return "Unknown"
	}
}

type Award struct {
	ID    int64       `json:"_id" sql:"_id"`
	Award string      `json:"award" sql:"award"`
	Types ContentType `json:"type" sql:"type"`
}

type AwardSQL struct {
	ID    int64       `json:"_id" sql:"_id"`
	Award string      `json:"award" sql:"award"`
	Types ContentType `json:"type" sql:"type"`
}

type AwardFrontend struct {
	ID    string      `json:"_id" sql:"_id"`
	Award string      `json:"award" sql:"award"`
	Types ContentType `json:"type" sql:"type"`
}

func CreateAward(id int64, award string, types ContentType) (*Award, error) {

	return &Award{
		ID:    id,
		Award: award,
		Types: types,
	}, nil
}

func AwardFromSQLNative(rows *sql.Rows) (*Award, error) {
	// create new coffee object to load into
	coffeeSQL := new(AwardSQL)

	// scan row into coffee object
	err := sqlstruct.Scan(coffeeSQL, rows)
	if err != nil {
		return nil, err
	}

	// create new coffee for the output
	award := &Award{
		ID:    coffeeSQL.ID,
		Award: coffeeSQL.Award,
		Types: coffeeSQL.Types,
	}

	return award, nil
}

func (i *Award) ToFrontend() *AwardFrontend {

	// create new coffee frontend
	mf := &AwardFrontend{
		ID:    fmt.Sprintf("%d", i.ID),
		Award: i.Award,
		Types: i.Types,
	}

	return mf
}

func (i *Award) ToSQLNative() *SQLInsertStatement {

	// create insertion statement and return
	return &SQLInsertStatement{
		Statement: "insert ignore into award(_id, award, types) values(?, ?, ?);",
		Values:    []interface{}{i.ID, i.Award, i.Types},
	}
}
