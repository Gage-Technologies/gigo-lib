package models

import (
	"database/sql"
	"fmt"
	"github.com/kisielk/sqlstruct"
	"time"
)

type ThreadComment struct {
	ID              int64             `json:"_id" sql:"_id"`
	Body            string            `json:"body" sql:"body"`
	Author          string            `json:"author" sql:"author"`
	AuthorID        int64             `json:"author_id" sql:"author_id"`
	CreatedAt       time.Time         `json:"created_at" sql:"created_at"`
	AuthorTier      TierType          `json:"author_tier" sql:"author_tier"`
	Coffee          uint64            `json:"coffee" sql:"coffee"`
	CommentId       int64             `json:"comment_id" sql:"comment_id"`
	Leads           bool              `json:"leads" sql:"leads"`
	Revision        int               `json:"revision" sql:"revision"`
	DiscussionLevel CommunicationType `json:"discussion_level" sql:"discussion_level"`
}

type ThreadCommentSQL struct {
	ID              int64             `json:"_id" sql:"_id"`
	Body            string            `json:"body" sql:"body"`
	Author          string            `json:"author" sql:"author"`
	AuthorID        int64             `json:"author_id" sql:"author_id"`
	CreatedAt       time.Time         `json:"created_at" sql:"created_at"`
	AuthorTier      TierType          `json:"author_tier" sql:"author_tier"`
	Coffee          uint64            `json:"coffee" sql:"coffee"`
	CommentId       int64             `json:"comment_id" sql:"comment_id"`
	Leads           bool              `json:"leads" sql:"leads"`
	Revision        int               `json:"revision" sql:"revision"`
	DiscussionLevel CommunicationType `json:"discussion_level" sql:"discussion_level"`
}

type ThreadCommentFrontend struct {
	ID              string            `json:"_id" sql:"_id"`
	Body            string            `json:"body" sql:"body"`
	Author          string            `json:"author" sql:"author"`
	AuthorID        string            `json:"author_id" sql:"author_id"`
	CreatedAt       time.Time         `json:"created_at" sql:"created_at"`
	AuthorTier      TierType          `json:"author_tier" sql:"author_tier"`
	Coffee          string            `json:"coffee" sql:"coffee"`
	CommentId       string            `json:"comment_id" sql:"comment_id"`
	Leads           bool              `json:"leads" sql:"leads"`
	Revision        int               `json:"revision" sql:"revision"`
	DiscussionLevel CommunicationType `json:"discussion_level" sql:"discussion_level"`
	Thumbnail       string            `json:"thumbnail"`
}

func CreateThreadComment(id int64, body string, author string, authorid int64, createdat time.Time, authortier TierType, coffee uint64, commentId int64, leads bool, revision int, discussionLevel CommunicationType) (*ThreadComment, error) {

	return &ThreadComment{
		ID:              id,
		Body:            body,
		Author:          author,
		AuthorID:        authorid,
		CreatedAt:       createdat,
		AuthorTier:      authortier,
		Coffee:          coffee,
		CommentId:       commentId,
		Leads:           leads,
		Revision:        revision,
		DiscussionLevel: discussionLevel,
	}, nil
}

func ThreadCommentFromSQLNative(rows *sql.Rows) (*ThreadComment, error) {
	// create new discussion object to load into
	commentSQL := new(ThreadCommentSQL)

	// scan row into comment object
	err := sqlstruct.Scan(commentSQL, rows)
	if err != nil {
		return nil, err
	}

	// create new comment for the output
	comment := &ThreadComment{
		ID:              commentSQL.ID,
		Body:            commentSQL.Body,
		Author:          commentSQL.Author,
		AuthorID:        commentSQL.AuthorID,
		CreatedAt:       commentSQL.CreatedAt,
		AuthorTier:      commentSQL.AuthorTier,
		Coffee:          commentSQL.Coffee,
		CommentId:       commentSQL.CommentId,
		Leads:           commentSQL.Leads,
		Revision:        commentSQL.Revision,
		DiscussionLevel: commentSQL.DiscussionLevel,
	}

	return comment, nil
}

func (i *ThreadComment) ToFrontend() *ThreadCommentFrontend {

	// create new comment frontend
	mf := &ThreadCommentFrontend{
		ID:              fmt.Sprintf("%d", i.ID),
		Body:            i.Body,
		Author:          i.Author,
		AuthorID:        fmt.Sprintf("%d", i.AuthorID),
		CreatedAt:       i.CreatedAt,
		AuthorTier:      i.AuthorTier,
		Coffee:          fmt.Sprintf("%d", i.Coffee),
		CommentId:       fmt.Sprintf("%d", i.CommentId),
		Leads:           i.Leads,
		Revision:        i.Revision,
		DiscussionLevel: i.DiscussionLevel,
		Thumbnail:       fmt.Sprintf("/static/user/pfp/%v", i.AuthorID),
	}

	return mf
}

func (i *ThreadComment) ToSQLNative() []*SQLInsertStatement {

	sqlStatements := make([]*SQLInsertStatement, 0)

	sqlStatements = append(sqlStatements, &SQLInsertStatement{
		Statement: "insert ignore into thread_comment(_id, body, author, author_id, created_at, author_tier, coffee, comment_id, leads, revision, discussion_level) values(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);",
		Values: []interface{}{i.ID, i.Body, i.Author, i.AuthorID, i.CreatedAt, i.AuthorTier,
			i.Coffee, i.CommentId, i.Leads, i.Revision, i.DiscussionLevel},
	})

	// create insertion statement and return
	return sqlStatements
}
