package models

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	ti "github.com/gage-technologies/gigo-lib/db"
	"go.opentelemetry.io/otel/trace"
	"strings"
	"time"

	"github.com/kisielk/sqlstruct"
)

type JourneyUserMap struct {
	UserID int64         `json:"user_id" sql:"user_id"`
	Units  []JourneyUnit `json:"units" sql:"units"`
}

type JourneyUserMapSQL struct {
	UserID      int64     `json:"user_id" sql:"user_id"`
	Unit        int64     `json:"unit_id" sql:"unit_id"`
	DateStarted time.Time `json:"date_started" sql:"date_started"`
}

type JourneyUserMapFrontend struct {
	UserID string                 `json:"user_id" sql:"user_id"`
	Units  []*JourneyUnitFrontend `json:"units" sql:"units"`
}

func CreateJourneyUserMap(userId int64, units []JourneyUnit) (*JourneyUserMap, error) {
	return &JourneyUserMap{
		UserID: userId,
		Units:  units,
	}, nil
}

func JourneyUserMapFromSQLNative(ctx context.Context, span *trace.Span, tidb *ti.Database, rows *sql.Rows) (*JourneyUserMap, error) {
	var userId int64

	unitIDs := make([]interface{}, 0)
	paramSlots := make([]string, 0)
	journeyUnits := make([]JourneyUnit, 0)

	for rows.Next() {
		tempJ := new(JourneyUserMapSQL)
		err := sqlstruct.Scan(tempJ, rows)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("error scanning JourneyUserMap info in first scan: %v", err))
		}

		unitIDs = append(unitIDs, tempJ.Unit)
		userId = tempJ.UserID
		paramSlots = append(paramSlots, "?")
	}

	callerName := "JourneyUserMapFromSQLNative"

	res, err := tidb.QueryContext(ctx, span, &callerName, fmt.Sprintf("select * from journey_units where _id in (%s)", strings.Join(paramSlots, ",")), unitIDs...)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed to query for full JourneyUnits inside JourneyUserMapFromSQLNative query: %v, err: %v", fmt.Sprintf("select * from journey_units where _id in (%s)", strings.Join(paramSlots, ",")), err))
	}

	defer res.Close()

	for res.Next() {
		j, err := JourneyUnitFromSQLNative(ctx, span, tidb, res)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("failed to query for full JourneyUnits inside JourneyUserMapFromSQLNative, err: %v", err))
		}

		if j != nil {
			journeyUnits = append(journeyUnits, *j)
		} else {
			return nil, errors.New(fmt.Sprintf("journey unit is null from SQL native in JourneyUserMapFromSQLNative"))
		}
	}

	if len(journeyUnits) < 1 {
		return nil, errors.New(fmt.Sprintf("no units returned from user map with query: %v and params: %v", fmt.Sprintf("select * from journey_units where _id in (%s)", strings.Join(paramSlots, ",")), unitIDs))
	}

	return &JourneyUserMap{
		UserID: userId,
		Units:  journeyUnits,
	}, nil
}

func (b *JourneyUserMap) ToFrontend() *JourneyUserMapFrontend {

	units := make([]*JourneyUnitFrontend, 0)

	for _, u := range b.Units {
		units = append(units, u.ToFrontend())
	}

	return &JourneyUserMapFrontend{
		UserID: fmt.Sprintf("%d", b.UserID),
		Units:  units,
	}
}

func (b *JourneyUserMap) ToSQLNative() ([]*SQLInsertStatement, error) {
	sqlStatements := make([]*SQLInsertStatement, 0)

	for _, u := range b.Units {
		sqlStatements = append(sqlStatements, &SQLInsertStatement{
			Statement: "insert ignore into journey_user_map(user_id, unit_id, started_at) values(?,?,?);",
			Values:    []interface{}{b.UserID, u.ID, time.Now()},
		})
	}

	return sqlStatements, nil
}
