package models

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/kisielk/sqlstruct"
	"time"
)

type Nemesis struct {
	ID                        int64      `json:"id" sql:"_id"`
	AntagonistID              int64      `json:"antagonist_id" sql:"antagonist_id"`
	AntagonistName            string     `json:"antagonist_name" sql:"antagonist_name"`
	AntagonistTowersCaptured  uint64     `json:"antagonist_towers_captured" sql:"antagonist_towers_captured"`
	ProtagonistID             int64      `json:"protagonist_id" sql:"protagonist_id"`
	ProtagonistName           string     `json:"protagonist_name" sql:"protagonist_name"`
	ProtagonistTowersCaptured uint64     `json:"protagonist_towers_captured" sql:"protagonist_towers_captured"`
	TimeOfVillainy            time.Time  `json:"time_of_villainy" sql:"time_of_villainy"`
	Victor                    *int64     `json:"victor" sql:"victor"`
	IsAccepted                bool       `json:"is_accepted" sql:"is_accepted"`
	EndTime                   *time.Time `json:"end_time" sql:"end_time"`
}

type NemesisSQL struct {
	ID                        int64      `json:"id" sql:"_id"`
	AntagonistID              int64      `json:"antagonist_id" sql:"antagonist_id"`
	AntagonistName            string     `json:"antagonist_name" sql:"antagonist_name"`
	AntagonistTowersCaptured  uint64     `json:"antagonist_towers_captured" sql:"antagonist_towers_captured"`
	ProtagonistID             int64      `json:"protagonist_id" sql:"protagonist_id"`
	ProtagonistName           string     `json:"protagonist_name" sql:"protagonist_name"`
	ProtagonistTowersCaptured uint64     `json:"protagonist_towers_captured" sql:"protagonist_towers_captured"`
	TimeOfVillainy            time.Time  `json:"time_of_villainy" sql:"time_of_villainy"`
	Victor                    *int64     `json:"victor" sql:"victor"`
	IsAccepted                bool       `json:"is_accepted" sql:"is_accepted"`
	EndTime                   *time.Time `json:"end_time" sql:"end_time"`
}

type NemesisFrontend struct {
	ID                        string     `json:"id" sql:"_id"`
	AntagonistID              string     `json:"antagonist_id" sql:"antagonist_id"`
	AntagonistName            string     `json:"antagonist_name" sql:"antagonist_name"`
	AntagonistTowersCaptured  string     `json:"antagonist_towers_captured" sql:"antagonist_towers_captured"`
	ProtagonistID             string     `json:"protagonist_id" sql:"protagonist_id"`
	ProtagonistName           string     `json:"protagonist_name" sql:"protagonist_name"`
	ProtagonistTowersCaptured string     `json:"protagonist_towers_captured" sql:"protagonist_towers_captured"`
	TimeOfVillainy            time.Time  `json:"time_of_villainy" sql:"time_of_villainy"`
	Victor                    *string    `json:"victor" sql:"victor"`
	IsAccepted                bool       `json:"is_accepted" sql:"is_accepted"`
	EndTime                   *time.Time `json:"end_time" sql:"end_time"`
}

type NemesisHistory struct {
	ID                    int64     `json:"id" sql:"_id"`
	MatchID               int64     `json:"match_id" sql:"match_id"`
	AntagonistID          int64     `json:"antagonist_id" sql:"antagonist_id"`
	ProtagonistID         int64     `json:"protagonist_id" sql:"protagonist_id"`
	ProtagonistTowersHeld int64     `json:"protagonist_towers_held" sql:"protagonist_towers_held"`
	AntagonistTowersHeld  int64     `json:"antagonist_towers_held" sql:"antagonist_towers_held"`
	ProtagonistTotalXP    int64     `json:"protagonist_total_xp" sql:"protagonist_total_xp"`
	AntagonistTotalXP     int64     `json:"antagonist_total_xp" sql:"antagonist_total_xp"`
	IsAlerted             bool      `json:"is_alerted" sql:"is_alerted"`
	CreatedAt             time.Time `json:"created_at" sql:"created_at"`
}

func CreateNemesis(id int64, antagID int64, antagName string, protageID int64, protagName string, timeOfVillainy time.Time, victor *int64,
	isAccepted bool, endTime *time.Time, protagonistTowersCaptured uint64, antagonistTowersCaptured uint64) *Nemesis {
	return &Nemesis{
		ID:                        id,
		AntagonistID:              antagID,
		AntagonistName:            antagName,
		ProtagonistID:             protageID,
		ProtagonistName:           protagName,
		TimeOfVillainy:            timeOfVillainy,
		Victor:                    victor,
		IsAccepted:                isAccepted,
		EndTime:                   endTime,
		ProtagonistTowersCaptured: protagonistTowersCaptured,
		AntagonistTowersCaptured:  antagonistTowersCaptured,
	}
}

func NemesisFromSQLNative(rows *sql.Rows) (*Nemesis, error) {
	// create new user stats object to load into
	var nemesisSQL NemesisSQL

	for rows.Next() {
		//scan row into user stats object
		err := sqlstruct.Scan(&nemesisSQL, rows)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("Error scanning nemesis in first scan: %v", err))
		}
	}

	// create new user stats

	return &Nemesis{
		ID:                        nemesisSQL.ID,
		AntagonistID:              nemesisSQL.AntagonistID,
		AntagonistName:            nemesisSQL.AntagonistName,
		AntagonistTowersCaptured:  nemesisSQL.AntagonistTowersCaptured,
		ProtagonistTowersCaptured: nemesisSQL.ProtagonistTowersCaptured,
		ProtagonistID:             nemesisSQL.ProtagonistID,
		ProtagonistName:           nemesisSQL.ProtagonistName,
		TimeOfVillainy:            nemesisSQL.TimeOfVillainy,
		Victor:                    nemesisSQL.Victor,
		EndTime:                   nemesisSQL.EndTime,
	}, nil
}

func (i *Nemesis) ToFrontend() *NemesisFrontend {
	victor := fmt.Sprintf("%d", i.Victor)

	return &NemesisFrontend{
		ID:                        fmt.Sprintf("%d", i.ID),
		AntagonistID:              fmt.Sprintf("%d", i.AntagonistID),
		AntagonistName:            i.AntagonistName,
		AntagonistTowersCaptured:  fmt.Sprintf("%d", i.AntagonistTowersCaptured),
		ProtagonistTowersCaptured: fmt.Sprintf("%d", i.ProtagonistTowersCaptured),
		ProtagonistName:           i.ProtagonistName,
		ProtagonistID:             fmt.Sprintf("%d", i.ProtagonistID),
		TimeOfVillainy:            i.TimeOfVillainy,
		Victor:                    &victor,
		IsAccepted:                i.IsAccepted,
		EndTime:                   i.EndTime,
	}
}

func (i *Nemesis) ToSQLNative() []*SQLInsertStatement {
	sqlStatements := make([]*SQLInsertStatement, 0)

	sqlStatements = append(sqlStatements, &SQLInsertStatement{
		Statement: "insert ignore into nemesis(_id, antagonist_id, antagonist_name, antagonist_towers_captured, protagonist_id, protagonist_name, protagonist_towers_captured, time_of_villainy, victor, is_accepted, end_time) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);",
		Values:    []interface{}{i.ID, i.AntagonistID, i.AntagonistName, i.AntagonistTowersCaptured, i.ProtagonistID, i.ProtagonistName, i.ProtagonistTowersCaptured, i.TimeOfVillainy, i.Victor, i.IsAccepted, i.EndTime},
	})

	// create insertion statement and return
	return sqlStatements
}

func (i *NemesisHistory) ToSQLNative() []*SQLInsertStatement {
	sqlStatements := make([]*SQLInsertStatement, 0)

	sqlStatements = append(sqlStatements, &SQLInsertStatement{
		Statement: "insert ignore into nemesis_history(_id, match_id, antagonist_id, protagonist_id, protagonist_towers_held, antagonist_towers_held, protagonist_total_xp, antagonist_total_xp, is_alerted, created_at) values (?,?,?,?,?,?,?,?,?,?);",
		Values:    []interface{}{i.ID, i.MatchID, i.AntagonistID, i.ProtagonistID, i.ProtagonistTowersHeld, i.AntagonistTowersHeld, i.ProtagonistTotalXP, i.AntagonistTotalXP, i.IsAlerted, i.CreatedAt},
	})

	return sqlStatements
}
