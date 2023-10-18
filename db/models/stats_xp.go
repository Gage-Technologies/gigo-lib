package models

import (
	"time"
)

type StatsXP struct {
	StatsID    int64     `json:"stats_id" sql:"stats_id"`
	Expiration time.Time `json:"expiration" sql:"expiration"`
}

func CreateStatsXP(statsID int64, expiration time.Time) (*StatsXP, error) {
	return &StatsXP{
		StatsID:    statsID,
		Expiration: expiration,
	}, nil
}

func (i *StatsXP) ToSQLNative() []*SQLInsertStatement {
	sqlStatements := make([]*SQLInsertStatement, 0)

	sqlStatements = append(sqlStatements, &SQLInsertStatement{
		Statement: "insert ignore into stats_xp(stats_id, expiration) values(?,?);",
		Values:    []interface{}{i.StatsID, i.Expiration},
	})

	return sqlStatements
}
