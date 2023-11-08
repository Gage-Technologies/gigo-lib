// TODO finish
package models

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/kisielk/sqlstruct"
	"go.opentelemetry.io/otel"

	ti "github.com/gage-technologies/gigo-lib/db"
)

type Attempt struct {
	ID          int64     `json:"_id" sql:"_id"`
	PostTitle   string    `json:"post_title" sql:"post_title"`
	Description string    `json:"description" sql:"description"`
	Author      string    `json:"author" sql:"author"`
	AuthorID    int64     `json:"author_id" sql:"author_id"`
	CreatedAt   time.Time `json:"created_at" sql:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" sql:"updated_at"`
	RepoID      int64     `json:"repo_id" sql:"repo_id"`
	AuthorTier  TierType  `json:"author_tier" sql:"author_tier"`
	//Awards            []int64            `json:"awards" sql:"awards"`
	Coffee            uint64             `json:"coffee" sql:"coffee"`
	PostID            int64              `json:"post_id" sql:"post_id"`
	Closed            bool               `json:"closed" sql:"closed"`
	Success           bool               `json:"success" sql:"success"`
	ClosedDate        *time.Time         `json:"closed_date" sql:"closed_date"`
	Tier              TierType           `json:"tier" sql:"tier"`
	ParentAttempt     *int64             `json:"parent_attempt" sql:"parent_attempt"`
	WorkspaceSettings *WorkspaceSettings `json:"workspace_settings" sql:"workspace_settings"`
	PostType          ChallengeType      `json:"post_type" sql:"post_type"`
	StartTime         *time.Duration     `json:"start_time" sql:"start_time"`
}

type AttemptSQL struct {
	ID          int64     `json:"_id" sql:"_id"`
	PostTitle   string    `json:"post_title" sql:"post_title"`
	Description string    `json:"description" sql:"description"`
	Author      string    `json:"author" sql:"author"`
	AuthorID    int64     `json:"author_id" sql:"author_id"`
	CreatedAt   time.Time `json:"created_at" sql:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" sql:"updated_at"`
	RepoID      int64     `json:"repo_id" sql:"repo_id"`
	AuthorTier  TierType  `json:"author_tier" sql:"author_tier"`
	//Awards            []byte     `json:"awards" sql:"awards"`
	Coffee            uint64         `json:"coffee" sql:"coffee"`
	PostID            int64          `json:"post_id" sql:"post_id"`
	Closed            bool           `json:"closed" sql:"closed"`
	Success           bool           `json:"success" sql:"success"`
	ClosedDate        *time.Time     `json:"closed_date" sql:"closed_date"`
	Tier              TierType       `json:"tier" sql:"tier"`
	ParentAttempt     *int64         `json:"parent_attempt" sql:"parent_attempt"`
	WorkspaceSettings []byte         `json:"workspace_settings" sql:"workspace_settings"`
	PostType          ChallengeType  `json:"post_type" sql:"post_type"`
	StartTime         *time.Duration `json:"start_time" sql:"start_time"`
}

type AttemptFrontend struct {
	ID          string    `json:"_id" sql:"_id"`
	PostTitle   string    `json:"post_title" sql:"post_title"`
	Description string    `json:"description" sql:"description"`
	Author      string    `json:"author" sql:"author"`
	AuthorID    string    `json:"author_id" sql:"author_id"`
	CreatedAt   time.Time `json:"created_at" sql:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" sql:"updated_at"`
	RepoID      string    `json:"repo_id" sql:"repo_id"`
	AuthorTier  TierType  `json:"author_tier" sql:"author_tier"`
	//Awards        []string   `json:"awards" sql:"awards"`
	Coffee          string        `json:"coffee" sql:"coffee"`
	PostID          string        `json:"post_id" sql:"post_id"`
	Closed          bool          `json:"closed" sql:"closed"`
	Success         bool          `json:"success" sql:"success"`
	ClosedDate      *time.Time    `json:"closed_date" sql:"closed_date"`
	Tier            TierType      `json:"tier" sql:"tier"`
	ParentAttempt   *string       `json:"parent_attempt" sql:"parent_attempt"`
	Thumbnail       string        `json:"thumbnail"`
	PostType        ChallengeType `json:"post_type" sql:"post_type"`
	PostTypeString  string        `json:"post_type_string" sql:"post_type_string"`
	StartTimeMillis *int64        `json:"start_time_millis" sql:"start_time_millis"`
}

func CreateAttempt(id int64, postTitle string, description string, author string, authorID int64, createdAt time.Time, updatedAt time.Time,
	repoId int64, authorTier TierType, awards []int64, coffee uint64, postID int64, tier TierType, parentAttempt *int64, postType ChallengeType) (*Attempt, error) {

	return &Attempt{
		ID:          id,
		PostTitle:   postTitle,
		Description: description,
		Author:      author,
		AuthorID:    authorID,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
		RepoID:      repoId,
		AuthorTier:  authorTier,
		//Awards:        awards,
		Coffee:        coffee,
		PostID:        postID,
		Closed:        false,
		Success:       false,
		ClosedDate:    nil,
		Tier:          tier,
		ParentAttempt: parentAttempt,
		PostType:      postType,
	}, nil
}

func AttemptFromSQLNative(db *ti.Database, rows *sql.Rows) (*Attempt, error) {
	// create new attempt object to load into
	attemptSQL := new(AttemptSQL)

	// scan row into attempt object
	err := sqlstruct.Scan(attemptSQL, rows)
	if err != nil {
		return nil, err
	}

	ctx, span := otel.Tracer("gigo-core").Start(context.Background(), "gigo-lib")
	callerName := "AttemptFromSQLNative"
	awardRows, err := db.QueryContext(ctx, &span, &callerName, "select award_id from attempt_awards where attempt_id = ?", attemptSQL.ID)
	if err != nil {
		return nil, err
	}

	defer awardRows.Close()

	awards := make([]int64, 0)

	for awardRows.Next() {
		var award int64
		err = awardRows.Scan(&award)
		if err != nil {
			return nil, err
		}
		awards = append(awards, award)
	}

	// create workspace settings to unmarshall into
	var workspaceSettings *WorkspaceSettings
	if attemptSQL.WorkspaceSettings != nil {
		var ws WorkspaceSettings
		err := json.Unmarshal(attemptSQL.WorkspaceSettings, &ws)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshall workspace settings: %v", err)
		}
		workspaceSettings = &ws
	}

	// create new attempt for the output
	attempt := &Attempt{
		ID:          attemptSQL.ID,
		PostTitle:   attemptSQL.PostTitle,
		Description: attemptSQL.Description,
		Author:      attemptSQL.Author,
		AuthorID:    attemptSQL.AuthorID,
		CreatedAt:   attemptSQL.CreatedAt,
		UpdatedAt:   attemptSQL.UpdatedAt,
		RepoID:      attemptSQL.RepoID,
		AuthorTier:  attemptSQL.AuthorTier,
		Coffee:      attemptSQL.Coffee,
		//Awards:            awards,
		PostID:            attemptSQL.PostID,
		Closed:            attemptSQL.Closed,
		Success:           attemptSQL.Success,
		ClosedDate:        attemptSQL.ClosedDate,
		Tier:              attemptSQL.Tier,
		ParentAttempt:     attemptSQL.ParentAttempt,
		WorkspaceSettings: workspaceSettings,
		PostType:          attemptSQL.PostType,
		StartTime:         attemptSQL.StartTime,
	}

	return attempt, nil
}

func (i *Attempt) ToFrontend() *AttemptFrontend {
	//awards := make([]string, 0)

	//for b := range i.Awards {
	//	awards = append(awards, fmt.Sprintf("%d", b))
	//}

	// conditionally set parent attempt
	var parentAttempt *string
	if i.ParentAttempt != nil {
		pa := fmt.Sprintf("%d", *i.ParentAttempt)
		parentAttempt = &pa
	}

	// consitionally load start time to millis
	var startTime *int64
	if i.StartTime != nil {
		t := i.StartTime.Milliseconds()
		startTime = &t
	}

	// create new attempt frontend
	mf := &AttemptFrontend{
		ID:          fmt.Sprintf("%d", i.ID),
		PostTitle:   i.PostTitle,
		Description: i.Description,
		Author:      i.Author,
		AuthorID:    fmt.Sprintf("%d", i.AuthorID),
		CreatedAt:   i.CreatedAt,
		UpdatedAt:   i.UpdatedAt,
		RepoID:      fmt.Sprintf("%d", i.RepoID),
		AuthorTier:  i.AuthorTier,
		//Awards:        awards,
		PostID:          fmt.Sprintf("%d", i.PostID),
		Coffee:          fmt.Sprintf("%d", i.Coffee),
		Closed:          i.Closed,
		Success:         i.Success,
		ClosedDate:      i.ClosedDate,
		Tier:            i.Tier,
		ParentAttempt:   parentAttempt,
		Thumbnail:       fmt.Sprintf("/static/posts/t/%v", i.PostID),
		PostType:        i.PostType,
		PostTypeString:  i.PostType.String(),
		StartTimeMillis: startTime,
	}

	return mf
}

func (i *Attempt) ToSQLNative() ([]*SQLInsertStatement, error) {

	sqlStatements := make([]*SQLInsertStatement, 0)

	//if len(i.Awards) > 0 {
	//	for b := range i.Awards {
	//		awardStatement := SQLInsertStatement{
	//			Statement: "insert ignore into attempt_awards(attempt_id, award_id) values(?, ?);",
	//			Values:    []interface{}{i.ID, b},
	//		}
	//
	//		sqlStatements = append(sqlStatements, &awardStatement)
	//	}
	//}

	// marshall workspace settings to json buffer
	var buf []byte
	if i.WorkspaceSettings != nil {
		b, err := json.Marshal(i.WorkspaceSettings)
		if err != nil {
			return nil, fmt.Errorf("failed to marshall workspace settings: %v", err)
		}
		buf = b
	}

	sqlStatements = append(sqlStatements, &SQLInsertStatement{
		Statement: "insert ignore into attempt(_id, post_title, description, author, author_id, created_at, updated_at, repo_id, " +
			"author_tier, coffee, post_id, closed, success, closed_date, tier, parent_attempt, workspace_settings, post_type, start_time) values(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);",
		Values: []interface{}{i.ID, i.PostTitle, i.Description, i.Author, i.AuthorID, i.CreatedAt, i.UpdatedAt, i.RepoID, i.AuthorTier,
			i.Coffee, i.PostID, i.Closed, i.Success, i.ClosedDate, i.Tier, i.ParentAttempt, buf, i.PostType, i.StartTime},
	})

	// create insertion statement and return
	return sqlStatements, nil
}
