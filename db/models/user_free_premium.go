package models

import (
	"database/sql"
	"github.com/kisielk/sqlstruct"
	"time"
)

type UserFreePremium struct {
	Id        int64     `json:"id" sql:"id"`
	UserId    int64     `json:"user_id" sql:"user_id"`
	StartDate time.Time `json:"start_date" sql:"start_date"`
	EndDate   time.Time `json:"end_date" sql:"end_date"`
	Length    string    `json:"length" sql:"length"`
}

type UserFreePremiumSQL struct {
	Id        int64     `json:"id" sql:"id"`
	UserId    int64     `json:"user_id" sql:"user_id"`
	StartDate time.Time `json:"start_date" sql:"start_date"`
	EndDate   time.Time `json:"end_date" sql:"end_date"`
	Length    string    `json:"length" sql:"length"`
}

func CreateUserFreePremium(userId int64, startDate time.Time, endDate time.Time, length string, id int64) (*UserFreePremium, error) {
	return &UserFreePremium{
		Id:        id,
		UserId:    userId,
		StartDate: startDate,
		EndDate:   endDate,
		Length:    length,
	}, nil
}

func (i *UserFreePremium) ToSQLNative() []*SQLInsertStatement {
	sqlStatements := make([]*SQLInsertStatement, 0)

	sqlStatements = append(sqlStatements, &SQLInsertStatement{
		Statement: "insert ignore into user_free_premium(_id, user_id, start_date, end_date, length) values(?,?,?,?, ?);",
		Values:    []interface{}{i.Id, i.UserId, i.StartDate, i.EndDate, i.Length},
	})

	return sqlStatements
}

func UserFreePremiumFromSQLNative(rows *sql.Rows) (*UserFreePremium, error) {
	// create new discussion object to load into
	commentSQL := new(UserFreePremiumSQL)

	// scan row into comment object
	err := sqlstruct.Scan(commentSQL, rows)
	if err != nil {
		return nil, err
	}

	// create new comment for the output
	comment := &UserFreePremium{
		Id:        commentSQL.Id,
		UserId:    commentSQL.UserId,
		StartDate: commentSQL.StartDate,
		EndDate:   commentSQL.EndDate,
		Length:    commentSQL.Length,
	}

	return comment, nil
}
