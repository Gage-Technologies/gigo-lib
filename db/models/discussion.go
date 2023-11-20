// TODO finish
package models

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/gage-technologies/gigo-lib/db"
	"github.com/kisielk/sqlstruct"
	"go.opentelemetry.io/otel"
	"time"
)

type Discussion struct {
	ID              int64             `json:"_id" sql:"_id"`
	Body            string            `json:"body" sql:"body"`
	Author          string            `json:"author" sql:"author"`
	AuthorID        int64             `json:"author_id" sql:"author_id"`
	CreatedAt       time.Time         `json:"created_at" sql:"created_at"`
	UpdatedAt       time.Time         `json:"updated_at" sql:"updated_at"`
	AuthorTier      TierType          `json:"author_tier" sql:"author_tier"`
	Awards          []int64           `json:"awards" sql:"awards"`
	Coffee          uint64            `json:"coffee" sql:"coffee"`
	PostId          int64             `json:"post_id" sql:"post_id"`
	Title           string            `json:"title" sql:"title"`
	Tags            []int64           `json:"tags" sql:"tags"`
	Leads           bool              `json:"leads" sql:"leads"`
	Revision        int               `json:"revision" sql:"revision"`
	DiscussionLevel CommunicationType `json:"discussion_level" sql:"discussion_level"`
}

type DiscussionSQL struct {
	ID              int64             `json:"_id" sql:"_id"`
	Body            string            `json:"body" sql:"body"`
	Author          string            `json:"author" sql:"author"`
	AuthorID        int64             `json:"author_id" sql:"author_id"`
	CreatedAt       time.Time         `json:"created_at" sql:"created_at"`
	UpdatedAt       time.Time         `json:"updated_at" sql:"updated_at"`
	AuthorTier      TierType          `json:"author_tier" sql:"author_tier"`
	Coffee          uint64            `json:"coffee" sql:"coffee"`
	PostId          int64             `json:"post_id" sql:"post_id"`
	Title           string            `json:"title" sql:"title"`
	Revision        int               `json:"revision" sql:"revision"`
	Leads           bool              `json:"leads" sql:"leads"`
	DiscussionLevel CommunicationType `json:"discussion_level" sql:"discussion_level"`
}

type DiscussionFrontend struct {
	ID              string            `json:"_id" sql:"_id"`
	Body            string            `json:"body" sql:"body"`
	Author          string            `json:"author" sql:"author"`
	AuthorID        string            `json:"author_id" sql:"author_id"`
	CreatedAt       time.Time         `json:"created_at" sql:"created_at"`
	UpdatedAt       time.Time         `json:"updated_at" sql:"updated_at"`
	AuthorTier      TierType          `json:"author_tier" sql:"author_tier"`
	Awards          []string          `json:"awards" sql:"awards"`
	Coffee          string            `json:"coffee" sql:"coffee"`
	PostId          string            `json:"post_id" sql:"post_id"`
	Title           string            `json:"title" sql:"title"`
	Tags            []string          `json:"tags" sql:"tags"`
	Leads           bool              `json:"leads" sql:"leads"`
	Revision        int               `json:"revision" sql:"revision"`
	DiscussionLevel CommunicationType `json:"discussion_level" sql:"discussion_level"`
	Thumbnail       string            `json:"thumbnail"`
}

func CreateDiscussion(id int64, body string, author string, authorID int64, createdAt time.Time, updatedAt time.Time,
	authorTier TierType, awards []int64, coffee uint64, postID int64, title string, tags []int64, leads bool, revision int, discussionLevel CommunicationType) (*Discussion, error) {

	return &Discussion{
		ID:              id,
		Body:            body,
		Author:          author,
		AuthorID:        authorID,
		CreatedAt:       createdAt,
		UpdatedAt:       updatedAt,
		Awards:          awards,
		AuthorTier:      authorTier,
		Coffee:          coffee,
		PostId:          postID,
		Title:           title,
		Tags:            tags,
		Leads:           leads,
		Revision:        revision,
		DiscussionLevel: discussionLevel,
	}, nil
}

func DiscussionFromSQLNative(db *ti.Database, rows *sql.Rows) (*Discussion, error) {
	// create new discussion object to load into
	discussionSQL := new(DiscussionSQL)

	// scan row into discussion object
	err := sqlstruct.Scan(discussionSQL, rows)
	if err != nil {
		return nil, err
	}

	ctx, span := otel.Tracer("gigo-core").Start(context.Background(), "gigo-lib")
	callerName := "DiscussionFromSQLNative"
	awardRows, err := db.QueryContext(ctx, &span, &callerName, "select award_id from discussion_awards where discussion_id = ?", discussionSQL.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to query award link table for award ids: %v", err)
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

	ctx, span = otel.Tracer("gigo-core").Start(context.Background(), "gigo-lib")
	tagRows, err := db.QueryContext(ctx, &span, &callerName, "select tag_id from discussion_tags where discussion_id = ?", discussionSQL.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to query tag link table for tag ids: %v", err)
	}

	defer tagRows.Close()

	tags := make([]int64, 0)

	for tagRows.Next() {
		var tag int64
		err = tagRows.Scan(&tag)
		if err != nil {
			return nil, fmt.Errorf("failed to scan tag id from link tbale cursor: %v", err)
		}
		tags = append(tags, tag)
	}

	// create new discussion for the output
	discussion := &Discussion{
		ID:              discussionSQL.ID,
		Body:            discussionSQL.Body,
		Author:          discussionSQL.Author,
		AuthorID:        discussionSQL.AuthorID,
		CreatedAt:       discussionSQL.CreatedAt,
		UpdatedAt:       discussionSQL.UpdatedAt,
		AuthorTier:      discussionSQL.AuthorTier,
		Coffee:          discussionSQL.Coffee,
		Awards:          awards,
		PostId:          discussionSQL.PostId,
		Title:           discussionSQL.Title,
		Tags:            tags,
		Leads:           discussionSQL.Leads,
		Revision:        discussionSQL.Revision,
		DiscussionLevel: discussionSQL.DiscussionLevel,
	}

	return discussion, nil
}

func (i *Discussion) ToFrontend() *DiscussionFrontend {
	awards := make([]string, 0)

	for _, b := range i.Awards {
		awards = append(awards, fmt.Sprintf("%d", b))
	}

	tags := make([]string, 0)

	for _, b := range i.Tags {
		tags = append(tags, fmt.Sprintf("%d", b))
	}

	// create new discussion frontend
	mf := &DiscussionFrontend{
		ID:              fmt.Sprintf("%d", i.ID),
		Body:            i.Body,
		Author:          i.Author,
		AuthorID:        fmt.Sprintf("%d", i.AuthorID),
		CreatedAt:       i.CreatedAt,
		UpdatedAt:       i.UpdatedAt,
		AuthorTier:      i.AuthorTier,
		Awards:          awards,
		PostId:          fmt.Sprintf("%d", i.PostId),
		Coffee:          fmt.Sprintf("%d", i.Coffee),
		Title:           i.Title,
		Tags:            tags,
		Leads:           i.Leads,
		Revision:        i.Revision,
		DiscussionLevel: i.DiscussionLevel,
		Thumbnail:       fmt.Sprintf("/static/user/pfp/%v", i.AuthorID),
	}

	return mf
}

func (i *Discussion) ToSQLNative() []*SQLInsertStatement {

	sqlStatements := make([]*SQLInsertStatement, 0)

	if len(i.Awards) > 0 {
		for _, b := range i.Awards {
			awardStatement := SQLInsertStatement{
				Statement: "insert ignore into discussion_awards(discussion_id, award_id, revision) values(?, ?, ?);",
				Values:    []interface{}{i.ID, b, i.Revision},
			}

			sqlStatements = append(sqlStatements, &awardStatement)
		}
	}

	if len(i.Tags) > 0 {
		for _, b := range i.Tags {
			tagStatement := SQLInsertStatement{
				Statement: "insert ignore into discussion_tags(discussion_id, tag_id, revision) values(?, ?, ?);",
				Values:    []interface{}{i.ID, b, i.Revision},
			}

			sqlStatements = append(sqlStatements, &tagStatement)
		}
	}

	sqlStatements = append(sqlStatements, &SQLInsertStatement{
		Statement: "insert ignore into discussion(_id, body, author, author_id, created_at, updated_at, author_tier, coffee, post_id, title, leads, revision, discussion_level) values(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);",
		Values: []interface{}{i.ID, i.Body, i.Author, i.AuthorID, i.CreatedAt, i.UpdatedAt, i.AuthorTier,
			i.Coffee, i.PostId, i.Title, i.Leads, i.Revision, i.DiscussionLevel},
	})

	// create insertion statement and return
	return sqlStatements
}
