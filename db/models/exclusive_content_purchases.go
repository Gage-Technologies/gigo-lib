package models

import (
	"database/sql"
	"fmt"
	ti "github.com/gage-technologies/gigo-lib/db"
	"github.com/kisielk/sqlstruct"
	"time"
)

type ExclusiveContentPurchases struct {
	UserId int64     `json:"user_id" sql:"user_id"`
	Post   int64     `json:"post" sql:"post"`
	Date   time.Time `json:"date" sql:"date"`
}

type ExclusiveContentPurchasesSQL struct {
	UserId int64     `json:"user_id" sql:"user_id"`
	Post   int64     `json:"post" sql:"post"`
	Date   time.Time `json:"date" sql:"date"`
}

type ExclusiveContentPurchasesFrontend struct {
	UserId string    `json:"user_id" sql:"user_id"`
	Post   string    `json:"post" sql:"post"`
	Date   time.Time `json:"date" sql:"date"`
}

func CreateExclusiveContentPurchases(userId int64, post int64) (*ExclusiveContentPurchases, error) {

	return &ExclusiveContentPurchases{
		UserId: userId,
		Post:   post,
		Date:   time.Now(),
	}, nil
}

func ExclusiveContentPurchasesFromSQLNative(db *ti.Database, rows *sql.Rows) (*ExclusiveContentPurchases, error) {
	// create new attempt object to load into
	attemptSQL := new(ExclusiveContentPurchasesSQL)

	// scan row into attempt object
	err := sqlstruct.Scan(attemptSQL, rows)
	if err != nil {
		return nil, err
	}

	// create new attempt for the output
	attempt := &ExclusiveContentPurchases{
		UserId: attemptSQL.UserId,
		Post:   attemptSQL.Post,
		Date:   attemptSQL.Date,
	}

	return attempt, nil
}

func (i *ExclusiveContentPurchases) ToFrontend() *ExclusiveContentPurchasesFrontend {

	// create new attempt frontend
	mf := &ExclusiveContentPurchasesFrontend{
		UserId: fmt.Sprintf("%d", i.UserId),
		Post:   fmt.Sprintf("%d", i.Post),
		Date:   i.Date,
	}

	return mf
}

func (i *ExclusiveContentPurchases) ToSQLNative() *SQLInsertStatement {

	sqlStatements := &SQLInsertStatement{
		Statement: "insert ignore into exclusive_content_purchases(user_id, post, date) values(?, ?, ?);",
		Values:    []interface{}{i.UserId, i.Post, i.Date},
	}

	// create insertion statement and return
	return sqlStatements
}
