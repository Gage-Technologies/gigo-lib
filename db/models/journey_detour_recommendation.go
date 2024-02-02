package models

import (
	"database/sql"
	"fmt"
	"gvisor.dev/gvisor/pkg/sentry/kernel/time"

	"github.com/kisielk/sqlstruct"
)

type JourneyDetourRecommendation struct {
	ID              int64     `json:"_id" sql:"_id"`
	UserID          int64     `json:"user_id" sql:"user_id"`
	RecommendedUnit int64     `json:"recommended_unit" sql:"recommended_unit"`
	CreatedAt       time.Time `json:"created_at" sql:"created_at"`
	FromTaskID      int64     `json:"from_task_id" sql:"from_task_id"`
	Accepted        bool      `json:"accepted" sql:"accepted"`
}

type JourneyDetourRecommendationSQL struct {
	ID              int64     `json:"_id" sql:"_id"`
	UserID          int64     `json:"user_id" sql:"user_id"`
	RecommendedUnit int64     `json:"recommended_unit" sql:"recommended_unit"`
	CreatedAt       time.Time `json:"created_at" sql:"created_at"`
	FromTaskID      int64     `json:"from_task_id" sql:"from_task_id"`
	Accepted        bool      `json:"accepted" sql:"accepted"`
}

type JourneyDetourRecommendationFrontend struct {
	ID              string    `json:"_id" sql:"_id"`
	UserID          string    `json:"user_id" sql:"user_id"`
	RecommendedUnit string    `json:"recommended_unit" sql:"recommended_unit"`
	CreatedAt       time.Time `json:"created_at" sql:"created_at"`
	FromTaskID      string    `json:"from_task_id" sql:"from_task_id"`
	Accepted        bool      `json:"accepted" sql:"accepted"`
}

func CreateJourneyDetourRecommendation(id int64, userId int64, recommendedUnit int64, fromTaskID int64, accepted bool, createdAt time.Time) (*JourneyDetourRecommendation, error) {
	return &JourneyDetourRecommendation{
		ID:              id,
		UserID:          userId,
		RecommendedUnit: recommendedUnit,
		CreatedAt:       createdAt,
		FromTaskID:      fromTaskID,
		Accepted:        accepted,
	}, nil
}

func JourneyDetourRecommendationFromSQLNative(rows *sql.Rows) (*JourneyDetourRecommendation, error) {
	JourneyDetourRecommendationSQL := new(JourneyDetourRecommendationSQL)
	err := sqlstruct.Scan(JourneyDetourRecommendationSQL, rows)
	if err != nil {
		return nil, fmt.Errorf("error scanning JourneyDetourRecommendation info in first scan: %v", err)
	}

	return &JourneyDetourRecommendation{
		ID:              JourneyDetourRecommendationSQL.ID,
		UserID:          JourneyDetourRecommendationSQL.UserID,
		RecommendedUnit: JourneyDetourRecommendationSQL.RecommendedUnit,
		CreatedAt:       JourneyDetourRecommendationSQL.CreatedAt,
		FromTaskID:      JourneyDetourRecommendationSQL.FromTaskID,
		Accepted:        JourneyDetourRecommendationSQL.Accepted,
	}, nil
}

func (b *JourneyDetourRecommendation) ToFrontend() *JourneyDetourRecommendationFrontend {
	return &JourneyDetourRecommendationFrontend{
		ID:              fmt.Sprintf("%d", b.ID),
		UserID:          fmt.Sprintf("%d", b.UserID),
		RecommendedUnit: fmt.Sprintf("%d", b.RecommendedUnit),
		CreatedAt:       b.CreatedAt,
		FromTaskID:      fmt.Sprintf("%d", b.FromTaskID),
		Accepted:        b.Accepted,
	}
}

func (b *JourneyDetourRecommendation) ToSQLNative() ([]*SQLInsertStatement, error) {
	sqlStatements := make([]*SQLInsertStatement, 0)

	sqlStatements = append(sqlStatements, &SQLInsertStatement{
		Statement: "insert ignore into journey_detour_recommendation(_id, user_id, recommended_unit, created_at, from_task_id, accepted) values(?,?,?,?,?,?);",
		Values:    []interface{}{b.ID, b.UserID, b.RecommendedUnit, b.CreatedAt, b.FromTaskID, b.Accepted},
	})

	return sqlStatements, nil
}
