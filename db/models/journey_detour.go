package models

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/kisielk/sqlstruct"
)

type JourneyDetour struct {
	DetourUnitID int64     `json:"detour_unit_id" sql:"detour_unit_id"`
	UserID       int64     `json:"user_id" sql:"user_id"`
	TaskID       int64     `json:"task_id" sql:"task_id"`
	StartedAt    time.Time `json:"started_at" sql:"started_at"`
}

type JourneyDetourSQL struct {
	DetourUnitID int64     `json:"detour_unit_id" sql:"detour_unit_id"`
	UserID       int64     `json:"user_id" sql:"user_id"`
	TaskID       int64     `json:"task_id" sql:"task_id"`
	StartedAt    time.Time `json:"started_at" sql:"started_at"`
}

type JourneyDetourFrontend struct {
	DetourUnitID string    `json:"detour_unit_id" sql:"detour_unit_id"`
	UserID       string    `json:"user_id" sql:"user_id"`
	TaskID       string    `json:"task_id" sql:"task_id"`
	StartedAt    time.Time `json:"started_at" sql:"started_at"`
}

type JourneyDetourSearch struct {
	DetourUnitID int64     `json:"detour_unit_id" sql:"detour_unit_id"`
	UserID       int64     `json:"user_id" sql:"user_id"`
	TaskID       int64     `json:"task_id" sql:"task_id"`
	StartedAt    time.Time `json:"started_at" sql:"started_at"`
}

func CreateJourneyDetour(detourUnitId int64, userId int64, taskId int64, startedAt time.Time) (*JourneyDetour, error) {
	return &JourneyDetour{
		DetourUnitID: detourUnitId,
		UserID:       userId,
		TaskID:       taskId,
		StartedAt:    startedAt,
	}, nil
}

func JourneyDetourFromSQLNative(rows *sql.Rows) (*JourneyDetour, error) {
	JourneyDetourSQL := new(JourneyDetourSQL)
	err := sqlstruct.Scan(JourneyDetourSQL, rows)
	if err != nil {
		return nil, fmt.Errorf("error scanning JourneyDetour info in first scan: %v", err)
	}

	return &JourneyDetour{
		DetourUnitID: JourneyDetourSQL.DetourUnitID,
		UserID:       JourneyDetourSQL.UserID,
		TaskID:       JourneyDetourSQL.TaskID,
		StartedAt:    JourneyDetourSQL.StartedAt,
	}, nil
}

func (b *JourneyDetour) ToFrontend() *JourneyDetourFrontend {
	return &JourneyDetourFrontend{
		DetourUnitID: fmt.Sprintf("%d", b.DetourUnitID),
		UserID:       fmt.Sprintf("%d", b.UserID),
		TaskID:       fmt.Sprintf("%d", b.TaskID),
		StartedAt:    b.StartedAt,
	}
}

func (b *JourneyDetour) ToSQLNative() ([]*SQLInsertStatement, error) {
	sqlStatements := make([]*SQLInsertStatement, 0)

	sqlStatements = append(sqlStatements, &SQLInsertStatement{
		Statement: "insert ignore into journey_detour(detour_unit_id, user_id, task_id, started_at) values(?,?,?,?);",
		Values:    []interface{}{b.DetourUnitID, b.UserID, b.TaskID, b.StartedAt},
	})

	return sqlStatements, nil
}
