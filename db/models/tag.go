package models

import (
	"database/sql"
	"fmt"
	"github.com/kisielk/sqlstruct"
)

type Tag struct {
	ID         int64  `json:"_id" sql:"_id"`
	Value      string `json:"value" sql:"value"`
	Official   bool   `json:"official" sql:"official"`
	UsageCount int64  `json:"usage_count" sql:"usage_count"`
}

type TagSQL struct {
	ID         int64  `json:"_id" sql:"_id"`
	Value      string `json:"value" sql:"value"`
	Official   bool   `json:"official" sql:"official"`
	UsageCount int64  `json:"usage_count" sql:"usage_count"`
}

type TagSearch struct {
	ID       int64  `json:"_id" sql:"_id"`
	Value    string `json:"value" sql:"value"`
	Official bool   `json:"official" sql:"official"`
}

type TagFrontend struct {
	ID         string `json:"_id" sql:"_id"`
	Value      string `json:"value" sql:"value"`
	Official   bool   `json:"official" sql:"official"`
	UsageCount int64  `json:"usage_count" sql:"usage_count"`
}

func CreateTag(id int64, value string) *Tag {
	return &Tag{
		ID:         id,
		Value:      value,
		UsageCount: 1,
	}
}

func TagFromSQLNative(rows *sql.Rows) (*Tag, error) {
	// create new tag object to load into
	tagSQL := new(Tag)

	// scan row into tag object
	err := sqlstruct.Scan(tagSQL, rows)
	if err != nil {
		return nil, err
	}

	return &Tag{
		ID:         tagSQL.ID,
		Value:      tagSQL.Value,
		Official:   tagSQL.Official,
		UsageCount: tagSQL.UsageCount,
	}, nil
}

func (t *Tag) ToSearch() *TagSearch {
	return &TagSearch{
		ID:       t.ID,
		Value:    t.Value,
		Official: t.Official,
	}
}

func (t *Tag) ToFrontend() *TagFrontend {
	return &TagFrontend{
		ID:         fmt.Sprintf("%d", t.ID),
		Value:      t.Value,
		Official:   t.Official,
		UsageCount: t.UsageCount,
	}
}

func (t *Tag) ToSQLNative() []*SQLInsertStatement {
	return []*SQLInsertStatement{
		{
			Statement: "insert ignore into tag(_id, value, official, usage_count) values (?, ?, ?, ?)",
			Values: []interface{}{
				t.ID, t.Value, t.Official, t.UsageCount,
			},
		},
	}
}
