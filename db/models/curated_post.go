package models

import (
	"context"
	"database/sql"
	"fmt"
	ti "github.com/gage-technologies/gigo-lib/db"
	"github.com/kisielk/sqlstruct"
	"go.opentelemetry.io/otel"
)

type ProficiencyType int

const (
	Beginner = iota
	Intermediate
	Advanced
	Any
)

func (p ProficiencyType) String() string {
	switch p {
	case Beginner:
		return "Beginner"
	case Intermediate:
		return "Intermediate"
	case Advanced:
		return "Advanced"
	case Any:
		return "Any"
	default:
		return "Any"
	}
}

type CuratedPost struct {
	ID                int64               `json:"_id" sql:"_id"`
	PostID            int64               `json:"post_id" sql:"post_id"`
	ProficiencyLevels []ProficiencyType   `json:"proficiency_type" sql:"proficiency_type"`
	PostLanguage      ProgrammingLanguage `json:"post_language" sql:"post_language"`
}

type CuratedPostSQL struct {
	ID           int64               `json:"_id" sql:"_id"`
	PostID       int64               `json:"post_id" sql:"post_id"`
	PostLanguage ProgrammingLanguage `json:"post_language" sql:"post_language"`
}

type CuratedPostFrontend struct {
	ID                string              `json:"_id" sql:"_id"`
	PostID            string              `json:"post_id" sql:"post_id"`
	ProficiencyLevels []ProficiencyType   `json:"proficiency_type" sql:"proficiency_type"`
	PostLanguage      ProgrammingLanguage `json:"post_language" sql:"post_language"`
}

func CreateCuratedPost(id int64, postId int64, proficiencyLevels []ProficiencyType, postLanguage ProgrammingLanguage) (*CuratedPost, error) {
	return &CuratedPost{
		ID:                id,
		PostID:            postId,
		ProficiencyLevels: proficiencyLevels,
		PostLanguage:      postLanguage,
	}, nil
}

func CuratedPostFromSQLNative(db *ti.Database, rows *sql.Rows) (*CuratedPost, error) {
	// create new recommended post object to load into
	recSql := new(CuratedPostSQL)

	// scan row into recommended post object
	err := sqlstruct.Scan(recSql, rows)
	if err != nil {
		return nil, err
	}

	ctx, span := otel.Tracer("gigo-core").Start(context.Background(), "gigo-lib")
	callerName := "CuratedPostFromSQLNative"

	// query link table to get award ids
	proficiencyRows, err := db.QueryContext(ctx, &span, &callerName, "select proficiency_type from curated_post_type where curated_id = ?", recSql.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to query curated post proficiency link table: %v", err)
	}

	// defer closure of lang rows
	defer proficiencyRows.Close()

	// slice to hold proficiency levels
	proficiencyLevels := make([]ProficiencyType, 0)

	// iterate cursor scanning proficiency ids and saving to the slice created above
	for proficiencyRows.Next() {
		var prof ProficiencyType
		err = proficiencyRows.Scan(&prof)
		if err != nil {
			return nil, fmt.Errorf("failed to scan prof id from link table cursor: %v", err)
		}
		proficiencyLevels = append(proficiencyLevels, prof)
	}

	return &CuratedPost{
		ID:                recSql.ID,
		PostID:            recSql.PostID,
		ProficiencyLevels: proficiencyLevels,
		PostLanguage:      recSql.PostLanguage,
	}, nil
}

func (i *CuratedPost) ToSQLNative() []*SQLInsertStatement {
	sqlStatements := []*SQLInsertStatement{
		{
			Statement: "insert ignore into curated_post(_id, post_id, post_language) values(?, ?, ?);",
			Values:    []interface{}{i.ID, i.PostID, i.PostLanguage},
		},
	}

	// iterate over proficiency levels and add them to the sql statement
	if len(i.ProficiencyLevels) > 0 {
		for _, a := range i.ProficiencyLevels {
			awardStatement := SQLInsertStatement{
				Statement: "insert ignore into curated_post_type(curated_id, proficiency_type) values(?, ?);",
				Values:    []interface{}{i.ID, a},
			}
			sqlStatements = append(sqlStatements, &awardStatement)
		}
	}
	// create insertion statement and return
	return sqlStatements
}

func (i *CuratedPost) ToFrontend() *CuratedPostFrontend {
	// create new recommended post frontend
	mf := &CuratedPostFrontend{
		ID:                fmt.Sprintf("%d", i.ID),
		PostID:            fmt.Sprintf("%d", i.PostID),
		ProficiencyLevels: i.ProficiencyLevels,
		PostLanguage:      i.PostLanguage,
	}

	return mf
}
