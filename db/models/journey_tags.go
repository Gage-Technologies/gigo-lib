package models

import (
	"database/sql"
	"github.com/kisielk/sqlstruct"
)

type JourneySourceType int

const (
	JourneyUnitType JourneySourceType = iota
	JourneyUnitProjectType
	JourneyUnitAttemptType
	JourneyUnitProjectAttemptType
)

type JourneyTags struct {
	JourneyID int64             `json:"journey_id" sql:"journey_id"`
	Value     string            `json:"value" sql:"value"`
	Type      JourneySourceType `json:"type" sql:"type"`
}

type JourneyTagsSQL struct {
	JourneyID int64  `json:"journey_id" sql:"journey_id"`
	Value     string `json:"value" sql:"value"`
	Type      int    `json:"type" sql:"type"`
}

type JourneyTagsFrontend struct {
	Value string `json:"value" sql:"value"`
}

func CreateJourneyTags(id int64, value string, journeyType JourneySourceType) *JourneyTags {
	return &JourneyTags{
		JourneyID: id,
		Value:     value,
		Type:      journeyType,
	}
}

func JourneyTagsFromSQLNative(rows *sql.Rows) (*JourneyTags, error) {
	// create new tag object to load into
	tagSQL := new(JourneyTags)

	// scan row into tag object
	err := sqlstruct.Scan(tagSQL, rows)
	if err != nil {
		return nil, err
	}

	return &JourneyTags{
		JourneyID: tagSQL.JourneyID,
		Value:     tagSQL.Value,
		Type:      JourneySourceType(tagSQL.Type),
	}, nil
}

func (t *JourneyTags) ToFrontend() *JourneyTagsFrontend {
	return &JourneyTagsFrontend{
		Value: t.Value,
	}
}

func (t *JourneyTags) ToSQLNative() []*SQLInsertStatement {
	return []*SQLInsertStatement{
		{
			Statement: "insert ignore into journey_tags(journey_id, value, type) values (?, ?, ?)",
			Values: []interface{}{
				t.JourneyID, t.Value, int(t.Type),
			},
		},
	}
}
