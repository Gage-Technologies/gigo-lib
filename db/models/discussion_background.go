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

type DiscussionBackground struct {
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
	RewardID        *int64            `json:"reward_id" sql:"reward_id"`
	Name            *string           `json:"name" sql:"name"`
	ColorPalette    *string           `json:"color_palette" sql:"color_palette"`
	RenderInFront   *bool             `json:"render_in_front" sql:"render_in_front"`
	UserStatus      UserStatus        `json:"user_status" sql:"user_status"`
}

type DiscussionBackgroundSQL struct {
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
	RewardID        *int64            `json:"reward_id" sql:"reward_id"`
	Name            *string           `json:"name" sql:"name"`
	ColorPalette    *string           `json:"color_palette" sql:"color_palette"`
	RenderInFront   *bool             `json:"render_in_front" sql:"render_in_front"`
	UserStatus      UserStatus        `json:"user_status" sql:"user_status"`
}

type DiscussionBackgroundFrontend struct {
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
	RewardID        *string           `json:"reward_id" sql:"reward_id"`
	Name            *string           `json:"name" sql:"name"`
	ColorPalette    *string           `json:"color_palette" sql:"color_palette"`
	RenderInFront   *bool             `json:"render_in_front" sql:"render_in_front"`
	UserStatus      UserStatus        `json:"user_status" sql:"user_status"`
}

func DiscussionBackgroundFromSQLNative(db *ti.Database, rows *sql.Rows) (*DiscussionBackground, error) {
	// create new discussion object to load into
	discussionSQL := new(DiscussionBackgroundSQL)

	// scan row into discussion object
	err := sqlstruct.Scan(discussionSQL, rows)
	if err != nil {
		return nil, err
	}

	ctx, span := otel.Tracer("gigo-core").Start(context.Background(), "gigo-lib")
	callerName := "DiscussionBackgroundFromSQLNative"
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
	discussion := &DiscussionBackground{
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
		RewardID:        discussionSQL.RewardID,
		Name:            discussionSQL.Name,
		ColorPalette:    discussionSQL.ColorPalette,
		RenderInFront:   discussionSQL.RenderInFront,
		UserStatus:      discussionSQL.UserStatus,
	}

	return discussion, nil
}

func (i *DiscussionBackground) ToFrontend() *DiscussionBackgroundFrontend {
	awards := make([]string, 0)

	for _, b := range i.Awards {
		awards = append(awards, fmt.Sprintf("%d", b))
	}

	tags := make([]string, 0)

	for _, b := range i.Tags {
		tags = append(tags, fmt.Sprintf("%d", b))
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

	// create new discussion frontend
	mf := &DiscussionBackgroundFrontend{
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
		RewardID:        rewardId,
		Name:            name,
		ColorPalette:    colorPalette,
		RenderInFront:   renderInFront,
		UserStatus:      i.UserStatus,
	}

	return mf
}
