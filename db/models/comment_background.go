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

type CommentBackground struct {
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
	RewardID        *int64            `json:"reward_id" sql:"reward_id"`
	Name            *string           `json:"name" sql:"name"`
	ColorPalette    *string           `json:"color_palette" sql:"color_palette"`
	RenderInFront   *bool             `json:"render_in_front" sql:"render_in_front"`
	UserStatus      UserStatus        `json:"user_status" sql:"user_status"`
}

type CommentBackgroundSQL struct {
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
	RewardID        *int64            `json:"reward_id" sql:"reward_id"`
	Name            *string           `json:"name" sql:"name"`
	ColorPalette    *string           `json:"color_palette" sql:"color_palette"`
	RenderInFront   *bool             `json:"render_in_front" sql:"render_in_front"`
	UserStatus      UserStatus        `json:"user_status" sql:"user_status"`
}

type CommentBackgroundFrontend struct {
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
	RewardID        *string           `json:"reward_id" sql:"reward_id"`
	Name            *string           `json:"name" sql:"name"`
	ColorPalette    *string           `json:"color_palette" sql:"color_palette"`
	RenderInFront   *bool             `json:"render_in_front" sql:"render_in_front"`
	UserStatus      UserStatus        `json:"user_status" sql:"user_status"`
}

func CommentBackgroundFromSQLNative(db *ti.Database, rows *sql.Rows) (*CommentBackground, error) {
	// create new discussion object to load into
	commentSQL := new(CommentBackgroundSQL)

	// scan row into comment object
	err := sqlstruct.Scan(commentSQL, rows)
	if err != nil {
		return nil, err
	}

	ctx, span := otel.Tracer("gigo-core").Start(context.Background(), "gigo-lib")
	callerName := "CommentBackgroundFromSQLNative"
	awardRows, err := db.QueryContext(ctx, &span, &callerName, "select award_id from comment_awards where comment_id = ? and revision = ?", commentSQL.ID, commentSQL.Revision)
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
	comment := &CommentBackground{
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
		RewardID:        commentSQL.RewardID,
		Name:            commentSQL.Name,
		ColorPalette:    commentSQL.ColorPalette,
		RenderInFront:   commentSQL.RenderInFront,
		UserStatus:      commentSQL.UserStatus,
	}

	return comment, nil
}

func (i *CommentBackground) ToFrontend() *CommentBackgroundFrontend {
	awards := make([]string, 0)

	for b := range i.Awards {
		awards = append(awards, fmt.Sprintf("%d", b))
	}

	var rewardId *string = nil
	if i.RewardID != nil {
		reward := fmt.Sprintf("%d", *i.RewardID)
		rewardId = &reward
	}

	var colorPalette *string = nil
	if i.ColorPalette != nil {
		colorPalette = i.ColorPalette
	}

	var renderInFront *bool = nil
	if i.RenderInFront != nil {
		renderInFront = i.RenderInFront
	}

	var name *string = nil
	if i.Name != nil {
		name = i.Name
	}

	// create new comment frontend
	mf := &CommentBackgroundFrontend{
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
		RewardID:        rewardId,
		Name:            name,
		ColorPalette:    colorPalette,
		RenderInFront:   renderInFront,
		UserStatus:      i.UserStatus,
	}

	return mf
}
