package models

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	ti "github.com/gage-technologies/gigo-lib/db"
	"go.opentelemetry.io/otel/trace"

	"github.com/kisielk/sqlstruct"
)

type JourneyUnit struct {
	ID          int64                 `json:"_id" sql:"_id"`
	Name        string                `json:"name" sql:"name"`
	UnitAbove   *int64                `json:"unit_above" sql:"unit_above"`
	UnitBelow   *int64                `json:"unit_below" sql:"unit_below"`
	Description string                `json:"description" sql:"description"`
	Langs       []ProgrammingLanguage `json:"langs" sql:"langs"`
	Tags        []string              `json:"tags" sql:"tags"`
	Published   bool                  `json:"published" sql:"published"`
	Color       string                `json:"color" sql:"color"`
}

type JourneyUnitSQL struct {
	ID          int64    `json:"_id" sql:"_id"`
	Name        string   `json:"name" sql:"name"`
	UnitAbove   *int64   `json:"unit_above" sql:"unit_above"`
	UnitBelow   *int64   `json:"unit_below" sql:"unit_below"`
	Description string   `json:"description" sql:"description"`
	Langs       []byte   `json:"langs" sql:"langs"`
	Tags        []string `json:"tags" sql:"tags"`
	Published   bool     `json:"published" sql:"published"`
	Color       string   `json:"color" sql:"color"`
}

type JourneyUnitFrontend struct {
	ID          string   `json:"_id" sql:"_id"`
	Name        string   `json:"name" sql:"name"`
	UnitAbove   *string  `json:"unit_above" sql:"unit_above"`
	UnitBelow   *string  `json:"unit_below" sql:"unit_below"`
	Description string   `json:"description" sql:"description"`
	Langs       []string `json:"langs" sql:"langs"`
	Tags        []string `json:"tags" sql:"tags"`
	Published   bool     `json:"published" sql:"published"`
	Color       string   `json:"color" sql:"color"`
}

func CreateJourneyUnit(id int64, name string, unitAbove *int64, unitBelow *int64,
	description string, langs []ProgrammingLanguage, tags []string, published bool, color string) (*JourneyUnit, error) {
	return &JourneyUnit{
		ID:          id,
		Name:        name,
		UnitAbove:   unitAbove,
		UnitBelow:   unitBelow,
		Description: description,
		Langs:       langs,
		Tags:        tags,
		Published:   published,
		Color:       color,
	}, nil
}

func JourneyUnitFromSQLNative(ctx context.Context, span *trace.Span, tidb *ti.Database, rows *sql.Rows) (*JourneyUnit, error) {
	JourneyUnitSQL := new(JourneyUnitSQL)
	err := sqlstruct.Scan(JourneyUnitSQL, rows)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed to marshal rows into JourneyUnitSQL, err: %v", err))
	}

	// create empty variable to hold workspace ports data
	var langs []ProgrammingLanguage

	// conditionally unmarshall json for workspace settings
	if JourneyUnitSQL.Langs != nil {
		err = json.Unmarshal(JourneyUnitSQL.Langs, &langs)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshall langs: %v", err)
		}
	}

	jUnit := JourneyUnit{
		ID:          JourneyUnitSQL.ID,
		Name:        JourneyUnitSQL.Name,
		UnitAbove:   JourneyUnitSQL.UnitAbove,
		UnitBelow:   JourneyUnitSQL.UnitBelow,
		Description: JourneyUnitSQL.Description,
		Langs:       langs,
		Published:   JourneyUnitSQL.Published,
		Color:       JourneyUnitSQL.Color,
	}

	callerName := "JourneyUnitFromSQLNative"

	res, err := tidb.QueryContext(ctx, span, &callerName, "select value from journey_unit_tags where unit_id = ?", jUnit.ID)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed to query for tags in JourneyUnitSQL, err: %v", err))
	}
	defer res.Close()

	tags := make([]string, 0)

	for res.Next() {
		var tag string
		err = res.Scan(&tag)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("failed to marshal tag rows into JourneyUnitSQL, err: %v", err))
		}
		tags = append(tags, tag)
	}

	jUnit.Tags = tags

	return &jUnit, nil
}

func (b *JourneyUnit) ToFrontend() *JourneyUnitFrontend {
	var unitAbove *string

	if b.UnitAbove != nil {
		unitStr := fmt.Sprintf("%d", *b.UnitAbove)
		unitAbove = &unitStr
	}

	var unitBelow *string

	if b.UnitBelow != nil {
		unitStr := fmt.Sprintf("%d", *b.UnitBelow)
		unitBelow = &unitStr
	}

	langs := make([]string, 0)

	for _, l := range b.Langs {
		langs = append(langs, l.String())
	}

	return &JourneyUnitFrontend{
		ID:          fmt.Sprintf("%d", b.ID),
		Name:        b.Name,
		Description: b.Description,
		UnitAbove:   unitAbove,
		UnitBelow:   unitBelow,
		Tags:        b.Tags,
		Langs:       langs,
		Published:   b.Published,
		Color:       b.Color,
	}
}

func (b *JourneyUnit) ToSQLNative() ([]*SQLInsertStatement, error) {
	sqlStatements := make([]*SQLInsertStatement, 0)

	var bytes []byte
	if b.Langs != nil {
		var err error
		bytes, err = json.Marshal(b.Langs)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal langs: %v", err)
		}
	}

	for _, t := range b.Tags {
		sqlStatements = append(sqlStatements, &SQLInsertStatement{
			Statement: "insert ignore into journey_unit_tags(unit_id, value) values(?,?);",
			Values:    []interface{}{b.ID, t},
		})
	}

	sqlStatements = append(sqlStatements, &SQLInsertStatement{
		Statement: "insert ignore into journey_units(_id, name, description, unit_above, unit_below, langs, published, color) values(?,?,?,?,?,?,?,?);",
		Values:    []interface{}{b.ID, b.Name, b.Description, b.UnitAbove, b.UnitBelow, bytes, b.Published, b.Color},
	})

	return sqlStatements, nil
}
