package models

import (
	"database/sql"
	"fmt"
	"github.com/kisielk/sqlstruct"
	"time"
)

type ThreadReply struct {
	ID              int64             `json:"_id" sql:"_id"`
	Body            string            `json:"body" sql:"body"`
	Author          string            `json:"author" sql:"author"`
	AuthorID        int64             `json:"author_id" sql:"author_id"`
	CreatedAt       time.Time         `json:"created_at" sql:"created_at"`
	AuthorTier      TierType          `json:"author_tier" sql:"author_tier"`
	Coffee          uint64            `json:"coffee" sql:"coffee"`
	ThreadCommentId int64             `json:"thread_comment_id" sql:"thread_comment_id"`
	Revision        int               `json:"revision" sql:"revision"`
	DiscussionLevel CommunicationType `json:"discussion_level" sql:"discussion_level"`
}

type ThreadReplySQL struct {
	ID              int64             `json:"_id" sql:"_id"`
	Body            string            `json:"body" sql:"body"`
	Author          string            `json:"author" sql:"author"`
	AuthorID        int64             `json:"author_id" sql:"author_id"`
	CreatedAt       time.Time         `json:"created_at" sql:"created_at"`
	AuthorTier      TierType          `json:"author_tier" sql:"author_tier"`
	Coffee          uint64            `json:"coffee" sql:"coffee"`
	ThreadCommentId int64             `json:"thread_comment_id" sql:"thread_comment_id"`
	Revision        int               `json:"revision" sql:"revision"`
	DiscussionLevel CommunicationType `json:"discussion_level" sql:"discussion_level"`
}

type ThreadReplyFrontend struct {
	ID              string            `json:"_id" sql:"_id"`
	Body            string            `json:"body" sql:"body"`
	Author          string            `json:"author" sql:"author"`
	AuthorID        string            `json:"author_id" sql:"author_id"`
	CreatedAt       time.Time         `json:"created_at" sql:"created_at"`
	AuthorTier      TierType          `json:"author_tier" sql:"author_tier"`
	Coffee          string            `json:"coffee" sql:"coffee"`
	ThreadCommentId string            `json:"thread_comment_id" sql:"thread_comment_id"`
	Revision        int               `json:"revision" sql:"revision"`
	DiscussionLevel CommunicationType `json:"discussion_level" sql:"discussion_level"`
	Thumbnail       string            `json:"thumbnail"`
}

func CreateThreadReply(id int64, body string, author string, authorid int64, createdat time.Time, authortier TierType, coffee uint64, ThreadCommentId int64, revision int, discussionLevel CommunicationType) (*ThreadReply, error) {

	return &ThreadReply{
		ID:              id,
		Body:            body,
		Author:          author,
		AuthorID:        authorid,
		CreatedAt:       createdat,
		AuthorTier:      authortier,
		Coffee:          coffee,
		ThreadCommentId: ThreadCommentId,
		Revision:        revision,
		DiscussionLevel: discussionLevel,
	}, nil
}

func ThreadReplyFromSQLNative(rows *sql.Rows) (*ThreadReply, error) {
	// create new discussion object to load into
	commentSQL := new(ThreadReplySQL)

	// scan row into comment object
	err := sqlstruct.Scan(commentSQL, rows)
	if err != nil {
		return nil, err
	}

	// create new comment for the output
	comment := &ThreadReply{
		ID:              commentSQL.ID,
		Body:            commentSQL.Body,
		Author:          commentSQL.Author,
		AuthorID:        commentSQL.AuthorID,
		CreatedAt:       commentSQL.CreatedAt,
		AuthorTier:      commentSQL.AuthorTier,
		Coffee:          commentSQL.Coffee,
		ThreadCommentId: commentSQL.ThreadCommentId,
		Revision:        commentSQL.Revision,
		DiscussionLevel: commentSQL.DiscussionLevel,
	}

	return comment, nil
}

func (i *ThreadReply) ToFrontend() *ThreadReplyFrontend {

	// create new comment frontend
	mf := &ThreadReplyFrontend{
		ID:              fmt.Sprintf("%d", i.ID),
		Body:            i.Body,
		Author:          i.Author,
		AuthorID:        fmt.Sprintf("%d", i.AuthorID),
		CreatedAt:       i.CreatedAt,
		AuthorTier:      i.AuthorTier,
		Coffee:          fmt.Sprintf("%d", i.Coffee),
		ThreadCommentId: fmt.Sprintf("%d", i.ThreadCommentId),
		Revision:        i.Revision,
		DiscussionLevel: i.DiscussionLevel,
		Thumbnail:       fmt.Sprintf("/static/user/pfp/%v", i.AuthorID),
	}

	return mf
}

func (i *ThreadReply) ToSQLNative() []*SQLInsertStatement {

	sqlStatements := make([]*SQLInsertStatement, 0)

	sqlStatements = append(sqlStatements, &SQLInsertStatement{
		Statement: "insert ignore into thread_reply(_id, body, author, author_id, created_at, author_tier, coffee, thread_comment_id, revision, discussion_level) values(?, ?, ?, ?, ?, ?, ?, ?, ?, ?);",
		Values: []interface{}{i.ID, i.Body, i.Author, i.AuthorID, i.CreatedAt, i.AuthorTier,
			i.Coffee, i.ThreadCommentId, i.Revision, i.DiscussionLevel},
	})

	// create insertion statement and return
	return sqlStatements
}
