package models

import (
	"database/sql"
	"errors"
	"fmt"
	ti "github.com/gage-technologies/gigo-lib/db"
	"github.com/kisielk/sqlstruct"
	"time"
)

type XPReason struct {
	ID     int64      `json:"id" sql:"_id"`
	UserID int64      `json:"user_id" sql:"user_id"`
	Date   *time.Time `json:"date" sql:"date"`
	Reason string     `json:"reason" sql:"reason"`
	XP     int64      `json:"xp" sql:"xp"`
}

type XPReasonSQL struct {
	ID     int64      `json:"id" sql:"_id"`
	UserID int64      `json:"user_id" sql:"user_id"`
	Date   *time.Time `json:"date" sql:"date"`
	Reason string     `json:"reason" sql:"reason"`
	XP     int64      `json:"xp" sql:"xp"`
}

type XPReasonFrontend struct {
	ID     string     `json:"id" sql:"_id"`
	UserID string     `json:"user_id" sql:"user_id"`
	Date   *time.Time `json:"date" sql:"date"`
	Reason string     `json:"reason" sql:"reason"`
	XP     int64      `json:"xp" sql:"xp"`
}

func CreateXPReason(id int64, userId int64, date *time.Time, reason string, xp int64) *XPReason {
	return &XPReason{
		ID:     id,
		UserID: userId,
		Date:   date,
		Reason: reason,
		XP:     xp,
	}
}

func XPReasonFromSQLNative(db *ti.Database, rows *sql.Rows) (*XPReason, error) {
	// create new user stats object to load into
	xPReasonSQL := new(XPReasonSQL)

	for rows.Next() {
		// scan row into user stats object
		err := sqlstruct.Scan(xPReasonSQL, rows)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("Error scanning user stats in first scan: %v", err))
		}
	}

	// create new user stats
	xPReason := &XPReason{
		ID:     xPReasonSQL.ID,
		UserID: xPReasonSQL.UserID,
		Date:   xPReasonSQL.Date,
		Reason: xPReasonSQL.Reason,
		XP:     xPReasonSQL.XP,
	}

	return xPReason, nil
}

func (i *XPReason) ToFrontend() *XPReasonFrontend {
	return &XPReasonFrontend{
		ID:     fmt.Sprintf("%d", i.ID),
		UserID: fmt.Sprintf("%d", i.UserID),
		Date:   i.Date,
		Reason: i.Reason,
		XP:     i.XP,
	}
}

func (i *XPReason) ToSQLNative() []*SQLInsertStatement {
	sqlStatements := make([]*SQLInsertStatement, 0)

	sqlStatements = append(sqlStatements, &SQLInsertStatement{
		Statement: "insert ignore into xp_reasons(_id, user_id, date, reason, xp) values (?, ?, ?, ?, ?);",
		Values:    []interface{}{i.ID, i.UserID, i.Date, i.Reason, i.XP},
	})

	// create insertion statement and return
	return sqlStatements
}
