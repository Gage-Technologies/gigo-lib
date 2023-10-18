package models

import (
	"database/sql"
	"fmt"
	"github.com/kisielk/sqlstruct"
)

type UpVote struct {
	ID             int64             `json:"_id" sql:"_id"`
	DiscussionType CommunicationType `json:"discussion_type" sql:"discussion_type"`
	DiscussionId   int64             `json:"discussion_id" sql:"discussion_id"`
	UserId         int64             `json:"user_id" sql:"user_id"`
}

type UpVoteSQL struct {
	ID             int64             `json:"_id" sql:"_id"`
	DiscussionType CommunicationType `json:"discussion_type" sql:"discussion_type"`
	DiscussionId   int64             `json:"discussion_id" sql:"discussion_id"`
	UserId         int64             `json:"user_id" sql:"user_id"`
}

type UpVoteFrontend struct {
	ID             string            `json:"_id" sql:"_id"`
	DiscussionType CommunicationType `json:"discussion_type" sql:"discussion_type"`
	DiscussionId   string            `json:"discussion_id" sql:"discussion_id"`
	UserId         string            `json:"user_id" sql:"user_id"`
}

func CreateUpVote(id int64, discussionType CommunicationType, discussionId int64, userId int64) *UpVote {
	return &UpVote{
		ID:             id,
		DiscussionType: discussionType,
		DiscussionId:   discussionId,
		UserId:         userId,
	}
}

func UpVoteFromSQLNative(rows *sql.Rows) (*UpVote, error) {
	// create new tag object to load into
	upVoteSql := new(UpVote)

	// scan row into tag object
	err := sqlstruct.Scan(upVoteSql, rows)
	if err != nil {
		return nil, err
	}

	return &UpVote{
		ID:             upVoteSql.ID,
		DiscussionType: upVoteSql.DiscussionType,
		DiscussionId:   upVoteSql.DiscussionId,
		UserId:         upVoteSql.UserId,
	}, nil
}

func (u *UpVote) ToFrontend() *UpVoteFrontend {
	return &UpVoteFrontend{
		ID:             fmt.Sprintf("%d", u.ID),
		DiscussionType: u.DiscussionType,
		DiscussionId:   fmt.Sprintf("%d", u.DiscussionId),
		UserId:         fmt.Sprintf("%d", u.UserId),
	}
}

func (u *UpVote) ToSQLNative() []*SQLInsertStatement {
	return []*SQLInsertStatement{
		{
			Statement: "insert ignore into up_vote(_id, discussion_type, discussion_id, user_id) values (?, ?, ?, ?)",
			Values: []interface{}{
				u.ID, u.DiscussionType, u.DiscussionId, u.UserId,
			},
		},
	}
}
