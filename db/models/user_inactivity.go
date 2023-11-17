package models

import (
	"database/sql"
	"github.com/kisielk/sqlstruct"
	"time"
)

type UserInactivity struct {
	UserId       int64     `json:"user_id" sql:"user_id"`
	LastLogin    time.Time `json:"last_login" sql:"last_login"`
	LastNotified time.Time `json:"last_notified" sql:"last_notified"`
	SendWeek     bool      `json:"send_week" sql:"send_week"`
	SendMonth    bool      `json:"send_month" sql:"send_month"`
}

type UserInactivitySQL struct {
	UserId       int64     `json:"user_id" sql:"user_id"`
	LastLogin    time.Time `json:"last_login" sql:"last_login"`
	LastNotified time.Time `json:"last_notified" sql:"last_notified"`
	SendWeek     bool      `json:"send_week" sql:"send_week"`
	SendMonth    bool      `json:"send_month" sql:"send_month"`
}

func CreateUserInactivity(userId int64, lastLogin time.Time, lastNotified time.Time, sendWeek bool, sendMonth bool) (*UserInactivity, error) {
	return &UserInactivity{
		UserId:       userId,
		LastLogin:    lastLogin,
		LastNotified: lastNotified,
		SendWeek:     sendWeek,
		SendMonth:    sendMonth,
	}, nil
}

func (i *UserInactivity) ToSQLNative() []*SQLInsertStatement {
	sqlStatements := make([]*SQLInsertStatement, 0)

	sqlStatements = append(sqlStatements, &SQLInsertStatement{
		Statement: "insert ignore into user_inactivity(user_id, last_login, last_notified, send_week, send_month) values(?, ?, ?, ?, ?);",
		Values:    []interface{}{i.UserId, i.LastLogin, i.LastNotified, i.SendWeek, i.SendMonth},
	})

	return sqlStatements
}

func UserInactivityFromSQLNative(rows *sql.Rows) (*UserInactivity, error) {
	// create new UserInactivity object to load into
	inactivitySQL := new(UserInactivitySQL)

	// scan row into userInactivity object
	err := sqlstruct.Scan(inactivitySQL, rows)
	if err != nil {
		return nil, err
	}

	// create new UserInactivity for the output
	inactive := &UserInactivity{
		UserId:       inactivitySQL.UserId,
		LastLogin:    inactivitySQL.LastLogin,
		LastNotified: inactivitySQL.LastNotified,
		SendWeek:     inactivitySQL.SendWeek,
		SendMonth:    inactivitySQL.SendMonth,
	}

	return inactive, nil
}
