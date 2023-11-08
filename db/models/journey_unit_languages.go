package models

import (
	"database/sql"
	"github.com/kisielk/sqlstruct"
)

type JourneyUnitLanguages struct {
	UnitID    int64  `json:"unit_id" sql:"unit_id"`
	Value     string `json:"value" sql:"value"`
	IsAttempt bool   `json:"is_attempt" sql:"is_attempt"`
}

type JourneyUnitLanguagesSQL struct {
	UnitID    int64  `json:"unit_id" sql:"unit_id"`
	Value     string `json:"value" sql:"value"`
	IsAttempt bool   `json:"is_attempt" sql:"is_attempt"`
}

type JourneyUnitLanguagesFrontend struct {
	Value string `json:"value" sql:"value"`
}

func CreateJourneyUnitLanguages(id int64, value ProgrammingLanguage, isAttempt bool) *JourneyUnitLanguages {
	return &JourneyUnitLanguages{
		UnitID:    id,
		Value:     value.String(),
		IsAttempt: isAttempt,
	}
}

func JourneyUnitLanguagesFromSQLNative(rows *sql.Rows) (*JourneyUnitLanguages, error) {
	// create new tag object to load into
	tagSQL := new(JourneyUnitLanguages)

	// scan row into tag object
	err := sqlstruct.Scan(tagSQL, rows)
	if err != nil {
		return nil, err
	}

	return &JourneyUnitLanguages{
		UnitID:    tagSQL.UnitID,
		Value:     tagSQL.Value,
		IsAttempt: tagSQL.IsAttempt,
	}, nil
}

func (t *JourneyUnitLanguages) ToFrontend() *JourneyUnitLanguagesFrontend {
	return &JourneyUnitLanguagesFrontend{
		Value: t.Value,
	}
}

func (t *JourneyUnitLanguages) ToSQLNative() []*SQLInsertStatement {
	return []*SQLInsertStatement{
		{
			Statement: "insert ignore into journey_unit_languages(unit_id, value, is_attempt) values (?, ?, ?)",
			Values: []interface{}{
				t.UnitID, t.Value, t.IsAttempt,
			},
		},
	}
}
