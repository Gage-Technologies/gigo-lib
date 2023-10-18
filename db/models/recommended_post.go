package models

import (
	"database/sql"
	"fmt"
	"github.com/kisielk/sqlstruct"
	"time"
)

type RecommendationType int

const (
	RecommendationTypeSematic = iota
	RecommendationTypeCollaborative
	RecommendationTypeHybrid
)

func (r RecommendationType) String() string {
	switch r {
	case RecommendationTypeSematic:
		return "Sematic"
	case RecommendationTypeCollaborative:
		return "Collaborative"
	case RecommendationTypeHybrid:
		return "Hybrid"
	default:
		return "Unknown"
	}
}

type RecommendedPost struct {
	ID            int64              `json:"_id" sql:"_id"`
	UserID        int64              `json:"user_id" sql:"user_id"`
	PostID        int64              `json:"post_id" sql:"post_id"`
	Type          RecommendationType `json:"type" sql:"type"`
	ReferenceID   int64              `json:"reference_id" sql:"reference_id"`
	Score         float32            `json:"score" sql:"score"`
	CreatedAt     time.Time          `json:"created_at" sql:"created_at"`
	ExpiresAt     time.Time          `json:"expires_at" sql:"expires_at"`
	ReferenceTier TierType           `json:"reference_tier" sql:"reference_tier"`
	Accepted      bool               `json:"accepted" sql:"accepted"`
	Views         int64              `json:"views" sql:"views"`
}

type RecommendedPostSQL struct {
	ID            int64              `sql:"_id"`
	UserID        int64              `sql:"user_id"`
	PostID        int64              `sql:"post_id"`
	Type          RecommendationType `sql:"type"`
	ReferenceID   int64              `sql:"reference_id"`
	Score         float32            `sql:"score"`
	CreatedAt     time.Time          `sql:"created_at"`
	ExpiresAt     time.Time          `sql:"expires_at"`
	ReferenceTier TierType           `json:"reference_tier" sql:"reference_tier"`
	Accepted      bool               `json:"accepted" sql:"accepted"`
	Views         int64              `json:"views" sql:"views"`
}

type RecommendedPostFrontend struct {
	ID            string             `json:"_id"`
	UserID        string             `json:"user_id"`
	PostID        string             `json:"post_id"`
	Type          RecommendationType `json:"type"`
	TypeString    string             `json:"type_string"`
	ReferenceID   string             `json:"reference_id"`
	Score         float32            `json:"score"`
	CreatedAt     time.Time          `json:"created_at"`
	ExpiresAt     time.Time          `json:"expires_at"`
	ReferenceTier TierType           `json:"reference_tier" sql:"reference_tier"`
}

func CreateRecommendedPost(id int64, userId int64, postId int64, recType RecommendationType, referenceId int64,
	score float32, createdAt time.Time, expiresAt time.Time, referenceTier TierType) (*RecommendedPost, error) {
	return &RecommendedPost{
		ID:            id,
		UserID:        userId,
		PostID:        postId,
		Type:          recType,
		ReferenceID:   referenceId,
		Score:         score,
		CreatedAt:     createdAt,
		ExpiresAt:     expiresAt,
		ReferenceTier: referenceTier,
	}, nil
}

func RecommendedPostFromSQLNative(rows *sql.Rows) (*RecommendedPost, error) {
	// create new recommended post object to load into
	recSql := new(RecommendedPostSQL)

	// scan row into recommended post object
	err := sqlstruct.Scan(recSql, rows)
	if err != nil {
		return nil, err
	}

	return &RecommendedPost{
		ID:            recSql.ID,
		UserID:        recSql.UserID,
		PostID:        recSql.PostID,
		Type:          recSql.Type,
		ReferenceID:   recSql.ReferenceID,
		Score:         recSql.Score,
		CreatedAt:     recSql.CreatedAt,
		ExpiresAt:     recSql.ExpiresAt,
		ReferenceTier: recSql.ReferenceTier,
		Accepted:      recSql.Accepted,
		Views:         recSql.Views,
	}, nil
}

func (i *RecommendedPost) ToFrontend() *RecommendedPostFrontend {
	// create new recommended post frontend
	mf := &RecommendedPostFrontend{
		ID:            fmt.Sprintf("%d", i.ID),
		UserID:        fmt.Sprintf("%d", i.UserID),
		PostID:        fmt.Sprintf("%d", i.PostID),
		Type:          i.Type,
		TypeString:    i.Type.String(),
		ReferenceID:   fmt.Sprintf("%d", i.ReferenceID),
		Score:         i.Score,
		CreatedAt:     i.CreatedAt,
		ExpiresAt:     i.ExpiresAt,
		ReferenceTier: i.ReferenceTier,
	}

	return mf
}

func (i *RecommendedPost) ToSQLNative() *SQLInsertStatement {
	sqlStatements := &SQLInsertStatement{
		Statement: "insert ignore into recommended_post(_id, user_id, post_id, type, reference_id, score, created_at, expires_at, reference_tier, accepted, views) values(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);",
		Values:    []interface{}{i.ID, i.UserID, i.PostID, i.Type, i.ReferenceID, i.Score, i.CreatedAt, i.ExpiresAt, i.ReferenceTier, i.Accepted, i.Views},
	}

	// create insertion statement and return
	return sqlStatements
}
