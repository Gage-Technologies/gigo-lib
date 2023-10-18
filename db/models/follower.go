package models

import (
	"database/sql"
	"fmt"
	ti "github.com/gage-technologies/gigo-lib/db"
	"github.com/kisielk/sqlstruct"
)

type Follower struct {
	Follower  int64 `json:"follower" sql:"follower"`
	Following int64 `json:"following" sql:"following"`
}

type FollowerSQL struct {
	Follower  int64 `json:"follower" sql:"follower"`
	Following int64 `json:"following" sql:"following"`
}

type FollowerFrontend struct {
	Follower  string `json:"follower" sql:"follower"`
	Following string `json:"following" sql:"following"`
}

func CreateFollower(follower int64, following int64) (*Follower, error) {

	return &Follower{
		Follower:  follower,
		Following: following,
	}, nil
}

func FollowerFromSQLNative(db *ti.Database, rows *sql.Rows) (*Follower, error) {
	// create new attempt object to load into
	attemptSQL := new(FollowerSQL)

	// scan row into attempt object
	err := sqlstruct.Scan(attemptSQL, rows)
	if err != nil {
		return nil, err
	}

	// create new attempt for the output
	attempt := &Follower{
		Follower:  attemptSQL.Follower,
		Following: attemptSQL.Following,
	}

	return attempt, nil
}

func (i *Follower) ToFrontend() *FollowerFrontend {

	// create new attempt frontend
	mf := &FollowerFrontend{
		Follower:  fmt.Sprintf("%d", i.Follower),
		Following: fmt.Sprintf("%d", i.Following),
	}

	return mf
}

func (i *Follower) ToSQLNative() *SQLInsertStatement {

	sqlStatements := &SQLInsertStatement{
		Statement: "insert ignore into follower(follower, following) values(?, ?);",
		Values:    []interface{}{i.Follower, i.Following},
	}

	// create insertion statement and return
	return sqlStatements
}
