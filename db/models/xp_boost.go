package models

import (
	"database/sql"
	"errors"
	"fmt"
	ti "github.com/gage-technologies/gigo-lib/db"
	"github.com/kisielk/sqlstruct"
	"time"
)

type XPBoost struct {
	ID      int64      `json:"id" sql:"_id"`
	UserID  int64      `json:"user_id" sql:"user_id"`
	EndDate *time.Time `json:"end_date" sql:"end_date"`
}

type XPBoostSQL struct {
	ID      int64      `json:"id" sql:"_id"`
	UserID  int64      `json:"user_id" sql:"user_id"`
	EndDate *time.Time `json:"end_date" sql:"end_date"`
}

type XPBoostFrontend struct {
	ID      string     `json:"id" sql:"_id"`
	UserID  string     `json:"user_id" sql:"user_id"`
	EndDate *time.Time `json:"end_date" sql:"end_date"`
}

func CreateXPBoost(id int64, userId int64, endTime *time.Time) *XPBoost {
	return &XPBoost{
		ID:      id,
		UserID:  userId,
		EndDate: endTime,
	}
}

func XPBoostFromSQLNative(db *ti.Database, rows *sql.Rows) (*XPBoost, error) {
	// create new user stats object to load into
	xPBoostSQL := new(XPBoostSQL)

	for rows.Next() {
		// scan row into user stats object
		err := sqlstruct.Scan(xPBoostSQL, rows)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("Error scanning user stats in first scan: %v", err))
		}
	}

	// create new user stats
	xPBoost := &XPBoost{
		ID:      xPBoostSQL.ID,
		UserID:  xPBoostSQL.UserID,
		EndDate: xPBoostSQL.EndDate,
	}

	return xPBoost, nil
}

func (i *XPBoost) ToFrontend() *XPBoostFrontend {
	return &XPBoostFrontend{
		ID:      fmt.Sprintf("%d", i.ID),
		UserID:  fmt.Sprintf("%d", i.UserID),
		EndDate: i.EndDate,
	}
}

func (i *XPBoost) ToSQLNative() []*SQLInsertStatement {
	sqlStatements := make([]*SQLInsertStatement, 0)

	sqlStatements = append(sqlStatements, &SQLInsertStatement{
		Statement: "insert ignore into xp_boosts(_id, user_id, end_date) values (?, ?, ?);",
		Values:    []interface{}{i.ID, i.UserID, i.EndDate},
	})

	// create insertion statement and return
	return sqlStatements
}
