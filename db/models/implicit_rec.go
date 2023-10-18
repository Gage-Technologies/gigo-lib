package models

import (
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"github.com/kisielk/sqlstruct"
	"time"
)

type ImplicitAction int

const (
	ImplicitTypeAttemptStart = iota
	ImplicitTypeAttemptEnd
	ImplicitTypeChallengeStart
	ImplicitTypeChallengeEnd
	ImplicitTypeInteractiveStart
	ImplicitTypeClicked
	ImplicitTypeClickedOff
	ImplicitTypeClickedOwnedProject
	ImplicitTypeClickedOffOwnedProject
	ImplicitTypeOwnedProjectStart
	ImplicitTypeOwnedProjectEnd
	ImplicitTypeInteractiveEnd
)

func (r ImplicitAction) String() string {
	switch r {
	case ImplicitTypeAttemptStart:
		return "AttemptStart"
	case ImplicitTypeAttemptEnd:
		return "AttemptEnd"
	case ImplicitTypeChallengeStart:
		return "ChallengeStart"
	case ImplicitTypeChallengeEnd:
		return "ChallengeEnd"
	case ImplicitTypeInteractiveStart:
		return "InteractiveStart"
	case ImplicitTypeClicked:
		return "Clicked"
	case ImplicitTypeClickedOff:
		return "ClickedOff"
	case ImplicitTypeClickedOwnedProject:
		return "ClickedOwnedProject"
	case ImplicitTypeClickedOffOwnedProject:
		return "ClickedOffOwnedProject"
	case ImplicitTypeOwnedProjectStart:
		return "OwnedProjectStart"
	case ImplicitTypeOwnedProjectEnd:
		return "OwnedProjectEnd"
	case ImplicitTypeInteractiveEnd:
		return "InteractiveEnd"
	default:
		return "Unknown"
	}
}

type ImplicitRec struct {
	ID               int64          `sql:"_id"`
	UserID           int64          `sql:"user_id"`
	PostID           int64          `sql:"post_id"`
	SessionID        uuid.UUID      `sql:"session_id"`
	ImplicitAction   ImplicitAction `sql:"implicit_action"`
	CreatedAt        time.Time      `sql:"created_at"`
	UserTierAtAction TierType       `json:"user_tier_at_action" sql:"user_tier_at_action"`
}

type ImplicitRecSQL struct {
	ID               int64          `sql:"_id"`
	UserID           int64          `sql:"user_id"`
	PostID           int64          `sql:"post_id"`
	SessionID        uuid.UUID      `sql:"session_id"`
	ImplicitAction   ImplicitAction `sql:"implicit_action"`
	CreatedAt        time.Time      `sql:"created_at"`
	UserTierAtAction TierType       `json:"user_tier_at_action" sql:"user_tier_at_action"`
}

type ImplicitRecFrontend struct {
	ID               string         `sql:"_id"`
	UserID           string         `sql:"user_id"`
	PostID           string         `sql:"post_id"`
	SessionID        string         `sql:"session_id"`
	ImplicitAction   ImplicitAction `sql:"implicit_action"`
	CreatedAt        time.Time      `sql:"created_at"`
	UserTierAtAction TierType       `json:"user_tier_at_action" sql:"user_tier_at_action"`
}

func CreateImplicitRec(id int64, userId int64, postId int64, sessionId uuid.UUID, action ImplicitAction,
	createdAt time.Time, userTierAtAction TierType) *ImplicitRec {
	return &ImplicitRec{
		ID:               id,
		UserID:           userId,
		PostID:           postId,
		SessionID:        sessionId,
		ImplicitAction:   action,
		CreatedAt:        createdAt,
		UserTierAtAction: userTierAtAction,
	}
}

func ImplicitRecFromSQLNative(rows *sql.Rows) (*ImplicitRec, error) {
	// create new Implicit post object to load into
	recSql := new(ImplicitRecSQL)

	// scan row into Implicit post object
	err := sqlstruct.Scan(recSql, rows)
	if err != nil {
		return nil, err
	}

	return &ImplicitRec{
		ID:               recSql.ID,
		UserID:           recSql.UserID,
		PostID:           recSql.PostID,
		SessionID:        recSql.SessionID,
		ImplicitAction:   recSql.ImplicitAction,
		CreatedAt:        recSql.CreatedAt,
		UserTierAtAction: TierType(recSql.UserTierAtAction),
	}, nil
}

func (i *ImplicitRec) ToFrontend() *ImplicitRecFrontend {
	// create new Implicit post frontend
	mf := &ImplicitRecFrontend{
		ID:               fmt.Sprintf("%d", i.ID),
		UserID:           fmt.Sprintf("%d", i.UserID),
		PostID:           fmt.Sprintf("%d", i.PostID),
		SessionID:        i.SessionID.String(),
		CreatedAt:        i.CreatedAt,
		UserTierAtAction: i.UserTierAtAction,
		ImplicitAction:   i.ImplicitAction,
	}

	return mf
}

func (i *ImplicitRec) ToSQLNative() *SQLInsertStatement {
	sqlStatements := &SQLInsertStatement{
		Statement: "insert into implicit_rec(_id, user_id, post_id, session_id, implicit_action, created_at, user_tier_at_action) " +
			"values(?, ?, ?, uuid_to_bin(?), ?, ?, ?) on duplicate key update created_at = if(implicit_action in (1, 3, 6, 8, 10, 11), values(created_at), created_at);",
		Values: []interface{}{i.ID, i.UserID, i.PostID, i.SessionID, i.ImplicitAction, i.CreatedAt, i.UserTierAtAction},
	}

	// create insertion statement and return
	return sqlStatements
}
