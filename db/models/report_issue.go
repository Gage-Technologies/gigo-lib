package models

import (
	"database/sql"
	"errors"
	"fmt"
	ti "github.com/gage-technologies/gigo-lib/db"
	"github.com/kisielk/sqlstruct"
	"time"
)

type ReportIssue struct {
	Date   time.Time `json:"date" sql:"date"`
	UserId int64     `json:"user_id" sql:"user_id"`
	Page   string    `json:"page" sql:"page"`
	Issue  string    `json:"issue" sql:"issue"`
	Id     int64     `json:"id" sql:"id"`
}

type ReportIssueSQL struct {
	Date   time.Time `json:"date" sql:"date"`
	UserId int64     `json:"user_id" sql:"user_id"`
	Page   string    `json:"page" sql:"page"`
	Issue  string    `json:"issue" sql:"issue"`
	Id     int64     `json:"id" sql:"id"`
}

type ReportIssueFrontend struct {
	Date   time.Time `json:"date" sql:"date"`
	UserId string    `json:"user_id" sql:"user_id"`
	Page   string    `json:"page" sql:"page"`
	Issue  string    `json:"issue" sql:"issue"`
	Id     string    `json:"id" sql:"id"`
}

func CreateReportIssue(userId int64, page string, issue string, id int64) (*ReportIssue, error) {
	return &ReportIssue{
		Date:   time.Now(),
		UserId: userId,
		Page:   page,
		Issue:  issue,
		Id:     id,
	}, nil
}

func ReportIssueFromSQLNative(db *ti.Database, rows *sql.Rows) (*ReportIssue, error) {
	// create new user stats object to load into
	reportIssueSQL := new(ReportIssueSQL)

	for rows.Next() {
		// scan row into user stats object
		err := sqlstruct.Scan(reportIssueSQL, rows)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("Error scanning user stats in first scan: %v", err))
		}
	}

	// create new user stats
	reportIssue := &ReportIssue{
		Id:     reportIssueSQL.Id,
		Date:   reportIssueSQL.Date,
		UserId: reportIssueSQL.UserId,
		Page:   reportIssueSQL.Page,
		Issue:  reportIssueSQL.Issue,
	}

	return reportIssue, nil
}

func (i *ReportIssue) ToFrontend() *ReportIssueFrontend {
	return &ReportIssueFrontend{
		Id:     fmt.Sprintf("%d", i.Id),
		UserId: fmt.Sprintf("%d", i.UserId),
		Date:   i.Date,
		Page:   i.Page,
		Issue:  i.Issue,
	}
}

func (i *ReportIssue) ToSQLNative() *SQLInsertStatement {

	sqlStatements := &SQLInsertStatement{
		Statement: "insert into report_issue(_id, user_id, date, page, issue) values(?, ?, ?, ?, ?);",
		Values:    []interface{}{i.Id, i.UserId, i.Date, i.Page, i.Issue},
	}

	// create insertion statement and return
	return sqlStatements
}
