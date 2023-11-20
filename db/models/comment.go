package models

import (
	"context"
	"database/sql"
	"fmt"
	ti "github.com/gage-technologies/gigo-lib/db"
	"github.com/kisielk/sqlstruct"
	"go.opentelemetry.io/otel"
	"time"
)

type Comment struct {
	ID              int64             `json:"_id" sql:"_id"`
	Body            string            `json:"body" sql:"body"`
	Author          string            `json:"author" sql:"author"`
	AuthorID        int64             `json:"author_id" sql:"author_id"`
	CreatedAt       time.Time         `json:"created_at" sql:"created_at"`
	AuthorTier      TierType          `json:"author_tier" sql:"author_tier"`
	Awards          []int64           `json:"awards" sql:"awards"`
	Coffee          uint64            `json:"coffee" sql:"coffee"`
	DiscussionId    int64             `json:"discussion_id" sql:"discussion_id"`
	Leads           bool              `json:"leads" sql:"leads"`
	Revision        int               `json:"revision" sql:"revision"`
	DiscussionLevel CommunicationType `json:"discussion_level" sql:"discussion_level"`
}

type CommentSQL struct {
	ID              int64             `json:"_id" sql:"_id"`
	Body            string            `json:"body" sql:"body"`
	Author          string            `json:"author" sql:"author"`
	AuthorID        int64             `json:"author_id" sql:"author_id"`
	CreatedAt       time.Time         `json:"created_at" sql:"created_at"`
	AuthorTier      TierType          `json:"author_tier" sql:"author_tier"`
	Coffee          uint64            `json:"coffee" sql:"coffee"`
	DiscussionId    int64             `json:"discussion_id" sql:"discussion_id"`
	Leads           bool              `json:"leads" sql:"leads"`
	Revision        int               `json:"revision" sql:"revision"`
	DiscussionLevel CommunicationType `json:"discussion_level" sql:"discussion_level"`
}

type CommentFrontend struct {
	ID              string            `json:"_id" sql:"_id"`
	Body            string            `json:"body" sql:"body"`
	Author          string            `json:"author" sql:"author"`
	AuthorID        string            `json:"author_id" sql:"author_id"`
	CreatedAt       time.Time         `json:"created_at" sql:"created_at"`
	AuthorTier      TierType          `json:"author_tier" sql:"author_tier"`
	Awards          []string          `json:"awards" sql:"awards"`
	Coffee          string            `json:"coffee" sql:"coffee"`
	DiscussionId    string            `json:"discussion_id" sql:"discussion_id"`
	Leads           bool              `json:"leads" sql:"leads"`
	Revision        int               `json:"revision" sql:"revision"`
	DiscussionLevel CommunicationType `json:"discussion_level" sql:"discussion_level"`
	Thumbnail       string            `json:"thumbnail"`
}

func CreateComment(id int64, body string, author string, authorId int64, createdAt time.Time, authorTier TierType, awards []int64, coffee uint64, discussionId int64, leads bool, revision int, discussionLevel CommunicationType) (*Comment, error) {

	return &Comment{
		ID:              id,
		Body:            body,
		Author:          author,
		AuthorID:        authorId,
		CreatedAt:       createdAt,
		Awards:          awards,
		AuthorTier:      authorTier,
		Coffee:          coffee,
		DiscussionId:    discussionId,
		Leads:           leads,
		Revision:        revision,
		DiscussionLevel: discussionLevel,
	}, nil
}

func CommentFromSQLNative(db *ti.Database, rows *sql.Rows) (*Comment, error) {
	// create new discussion object to load into
	commentSQL := new(CommentSQL)

	// scan row into comment object
	err := sqlstruct.Scan(commentSQL, rows)
	if err != nil {
		return nil, err
	}

	ctx, span := otel.Tracer("gigo-core").Start(context.Background(), "gigo-lib")
	callerName := "CommentFromSQLNative"
	awardRows, err := db.QueryContext(ctx, &span, &callerName, "select award_id from comment_awards where comment_id = ?", commentSQL.ID)
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

	// create new comment for the output
	comment := &Comment{
		ID:              commentSQL.ID,
		Body:            commentSQL.Body,
		Author:          commentSQL.Author,
		AuthorID:        commentSQL.AuthorID,
		CreatedAt:       commentSQL.CreatedAt,
		AuthorTier:      commentSQL.AuthorTier,
		Coffee:          commentSQL.Coffee,
		Awards:          awards,
		DiscussionId:    commentSQL.DiscussionId,
		Leads:           commentSQL.Leads,
		Revision:        commentSQL.Revision,
		DiscussionLevel: commentSQL.DiscussionLevel,
	}

	return comment, nil
}

func (i *Comment) ToFrontend() *CommentFrontend {
	awards := make([]string, 0)

	for _, b := range i.Awards {
		awards = append(awards, fmt.Sprintf("%d", b))
	}

	// create new comment frontend
	mf := &CommentFrontend{
		ID:              fmt.Sprintf("%d", i.ID),
		Body:            i.Body,
		Author:          i.Author,
		AuthorID:        fmt.Sprintf("%d", i.AuthorID),
		CreatedAt:       i.CreatedAt,
		AuthorTier:      i.AuthorTier,
		Awards:          awards,
		DiscussionId:    fmt.Sprintf("%d", i.DiscussionId),
		Coffee:          fmt.Sprintf("%d", i.Coffee),
		Leads:           i.Leads,
		Revision:        i.Revision,
		DiscussionLevel: i.DiscussionLevel,
		Thumbnail:       fmt.Sprintf("/static/user/pfp/%v", i.AuthorID),
	}

	return mf
}

func (i *Comment) ToSQLNative() []*SQLInsertStatement {

	sqlStatements := make([]*SQLInsertStatement, 0)

	if len(i.Awards) > 0 {
		for _, b := range i.Awards {
			awardStatement := SQLInsertStatement{
				Statement: "insert ignore into comment_awards(comment_id, award_id, revision) values(?, ?, ?);",
				Values:    []interface{}{i.ID, b, i.Revision},
			}

			sqlStatements = append(sqlStatements, &awardStatement)
		}
	}

	sqlStatements = append(sqlStatements, &SQLInsertStatement{
		Statement: "insert ignore into comment(_id, body, author, author_id, created_at, author_tier, coffee, discussion_id, leads, revision, discussion_level) values(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);",
		Values: []interface{}{i.ID, i.Body, i.Author, i.AuthorID, i.CreatedAt, i.AuthorTier,
			i.Coffee, i.DiscussionId, i.Leads, i.Revision, i.DiscussionLevel},
	})

	// create insertion statement and return
	return sqlStatements
}
