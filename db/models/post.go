package models

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gage-technologies/gigo-lib/db"
	"github.com/google/uuid"
	"github.com/kisielk/sqlstruct"
	"go.opentelemetry.io/otel"
)

type Post struct {
	ID                      int64                 `json:"_id" sql:"_id"`
	Title                   string                `json:"title" sql:"title"`
	Description             string                `json:"description" sql:"description"`
	Author                  string                `json:"author" sql:"author"`
	AuthorID                int64                 `json:"author_id" sql:"author_id"`
	CreatedAt               time.Time             `json:"created_at" sql:"created_at"`
	UpdatedAt               time.Time             `json:"updated_at" sql:"updated_at"`
	RepoID                  int64                 `json:"repo_id" sql:"repo_id"`
	Tier                    TierType              `json:"tier" sql:"tier"`
	Awards                  []int64               `json:"awards" sql:"awards"`
	TopReply                *int64                `json:"top_reply,omitempty" sql:"top_reply"`
	Coffee                  uint64                `json:"coffee" sql:"coffee"`
	Tags                    []int64               `json:"tags" sql:"tags"`
	PostType                ChallengeType         `json:"post_type" sql:"post_type"`
	Views                   int64                 `json:"views" sql:"views"`
	Completions             int64                 `json:"completions" sql:"completions"`
	Attempts                int64                 `json:"attempts" sql:"attempts"`
	Languages               []ProgrammingLanguage `json:"languages" sql:"languages"`
	Published               bool                  `json:"published" sql:"published"`
	Visibility              PostVisibility        `json:"visibility" sql:"visibility"`
	StripePriceId           *string               `json:"stripe_price_id" sql:"stripe_price_id"`
	ChallengeCost           *string               `json:"challenge_cost" sql:"challenge_cost"`
	WorkspaceConfig         int64                 `json:"workspace_config" sql:"workspace_config"`
	WorkspaceConfigRevision int                   `json:"workspace_config_revision" sql:"workspace_config_revision"`
	WorkspaceSettings       *WorkspaceSettings    `json:"workspace_settings" sql:"workspace_settings"`
	Leads                   bool                  `json:"leads" sql:"leads"`
	Embedded                bool                  `json:"embedded" sql:"embedded"`
	Deleted                 bool                  `json:"deleted" sql:"deleted"`
	ExclusiveDescription    *string               `json:"exclusive_description,omitempty" sql:"exclusive_description"`
	ShareHash              *uuid.UUID               `json:"share_hash" sql:"share_hash"`
}

type PostSQL struct {
	ID                      int64          `json:"_id" sql:"_id"`
	Title                   string         `json:"title" sql:"title"`
	Description             string         `json:"description" sql:"description"`
	Author                  string         `json:"author" sql:"author"`
	AuthorID                int64          `json:"author_id" sql:"author_id"`
	CreatedAt               time.Time      `json:"created_at" sql:"created_at"`
	UpdatedAt               time.Time      `json:"updated_at" sql:"updated_at"`
	RepoID                  int64          `json:"repo_id" sql:"repo_id"`
	Tier                    TierType       `json:"tier" sql:"tier"`
	TopReply                *int64         `json:"top_reply,omitempty" sql:"top_reply"`
	Coffee                  uint64         `json:"coffee" sql:"coffee"`
	PostType                ChallengeType  `json:"post_type" sql:"post_type"`
	Views                   int64          `json:"views" sql:"views"`
	Completions             int64          `json:"completions" sql:"completions"`
	Attempts                int64          `json:"attempts" sql:"attempts"`
	Published               bool           `json:"published" sql:"published"`
	Visibility              PostVisibility `json:"visibility" sql:"visibility"`
	StripePriceId           *string        `json:"stripe_price_id" sql:"stripe_price_id"`
	ChallengeCost           *string        `json:"challenge_cost" sql:"challenge_cost"`
	WorkspaceConfig         int64          `json:"workspace_config" sql:"workspace_config"`
	WorkspaceConfigRevision int            `json:"workspace_config_revision" sql:"workspace_config_revision"`
	WorkspaceSettings       []byte         `json:"workspace_settings" sql:"workspace_settings"`
	Leads                   bool           `json:"leads" sql:"leads"`
	Embedded                bool           `json:"embedded" sql:"embedded"`
	Deleted                 bool           `json:"deleted" sql:"deleted"`
	ExclusiveDescription    *string        `json:"exclusive_description" sql:"exclusive_description"`
	ShareHash              *uuid.UUID      `json:"share_hash" sql:"share_hash"`
}

type PostFrontend struct {
	ID                      string                `json:"_id"`
	Title                   string                `json:"title"`
	Description             string                `json:"description"`
	Author                  string                `json:"author"`
	AuthorID                string                `json:"author_id"`
	CreatedAt               time.Time             `json:"created_at"`
	UpdatedAt               time.Time             `json:"updated_at"`
	RepoID                  string                `json:"repo_id"`
	Tier                    TierType              `json:"tier"`
	TierString              string                `json:"tier_string"`
	Awards                  []string              `json:"awards"`
	TopReply                *string               `json:"top_reply"`
	Coffee                  uint64                `json:"coffee"`
	PostType                ChallengeType         `json:"post_type"`
	PostTypeString          string                `json:"post_type_string"`
	Views                   int64                 `json:"views"`
	Completions             int64                 `json:"completions"`
	Attempts                int64                 `json:"attempts"`
	Languages               []ProgrammingLanguage `json:"languages"`
	LanguageStrings         []string              `json:"languages_strings"`
	Published               bool                  `json:"published"`
	Visibility              PostVisibility        `json:"visibility"`
	VisibilityString        string                `json:"visibility_string"`
	Tags                    []string              `json:"tags"`
	Thumbnail               string                `json:"thumbnail"`
	ChallengeCost           *string               `json:"challenge_cost"`
	WorkspaceConfig         string                `json:"workspace_config"`
	WorkspaceConfigRevision int                   `json:"workspace_config_revision"`
	Leads                   bool                  `json:"leads" sql:"leads"`
	Deleted                 bool                  `json:"deleted" sql:"deleted"`
	ExclusiveDescription    *string               `json:"exclusive_description"`
}

func CreatePost(id int64, title string, description string, author string, authorID int64, createdAt time.Time,
	updatedAt time.Time, repoId int64, tier TierType, awards []int64, topReply *int64, coffee uint64,
	postType ChallengeType, views int64, completions int64, attempts int64, language []ProgrammingLanguage,
	visibility PostVisibility, tags []int64, challengeCost *string, stripeId *string, workspaceCfg int64,
	workspaceCfgRevision int, workspaceSettings *WorkspaceSettings, leads bool, embedded bool, exclusiveDescription *string) (*Post, error) {

	return &Post{
		ID:                      id,
		Title:                   title,
		Description:             description,
		Author:                  author,
		AuthorID:                authorID,
		CreatedAt:               createdAt,
		UpdatedAt:               updatedAt,
		RepoID:                  repoId,
		TopReply:                topReply,
		Tier:                    tier,
		Awards:                  awards,
		Coffee:                  coffee,
		PostType:                postType,
		Views:                   views,
		Completions:             completions,
		Attempts:                attempts,
		Languages:               language,
		Visibility:              visibility,
		Tags:                    tags,
		ChallengeCost:           challengeCost,
		StripePriceId:           stripeId,
		WorkspaceConfig:         workspaceCfg,
		WorkspaceConfigRevision: workspaceCfgRevision,
		WorkspaceSettings:       workspaceSettings,
		Leads:                   leads,
		Embedded:                embedded,
		ExclusiveDescription:    exclusiveDescription,
	}, nil
}

func PostFromSQLNative(db *ti.Database, rows *sql.Rows) (*Post, error) {
	// create new post object to load into
	postSQL := new(PostSQL)

	// scan row into post object
	err := sqlstruct.Scan(postSQL, rows)
	if err != nil {
		return nil, fmt.Errorf("failed to scan post into sql object: %v", err)
	}

	ctx, span := otel.Tracer("gigo-core").Start(context.Background(), "gigo-lib")
	callerName := "PostFromSQLNative"
	// query link table to get award ids
	awardRows, err := db.QueryContext(ctx, &span, &callerName, "select award_id from post_awards where post_id = ?", postSQL.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to query post awards link table: %v", err)
	}

	// defer closure of cursor
	defer awardRows.Close()

	// create slice to hold award ids loaded from cursor
	awards := make([]int64, 0)

	// iterate cursor loading award ids and saving to id slice created abov
	for awardRows.Next() {
		var award int64
		err = awardRows.Scan(&award)
		if err != nil {
			return nil, fmt.Errorf("failed to scan award id from link table cursor: %v", err)
		}
		awards = append(awards, award)
	}

	ctx, span = otel.Tracer("gigo-core").Start(context.Background(), "gigo-lib")

	// query tag link table to get tab ids
	tagRows, err := db.QueryContext(ctx, &span, &callerName, "select tag_id from post_tags where post_id = ?", postSQL.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to query tag link table for tag ids: %v", err)
	}

	// defer closure of tag rows
	defer tagRows.Close()

	// create slice to hold tag ids
	tags := make([]int64, 0)

	// iterate cursor scanning tag ids and saving the to the slice created above
	for tagRows.Next() {
		var tag int64
		err = tagRows.Scan(&tag)
		if err != nil {
			return nil, fmt.Errorf("failed to scan tag id from link tbale cursor: %v", err)
		}
		tags = append(tags, tag)
	}

	ctx, span = otel.Tracer("gigo-core").Start(context.Background(), "gigo-lib")
	// query lang link table to get lang ids
	langRows, err := db.QueryContext(ctx, &span, &callerName, "select lang_id from post_langs where post_id =?", postSQL.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to query lang link table for lang ids: %v", err)
	}

	// defer closure of lang rows
	defer langRows.Close()

	// create slice to hold lang ids
	languages := make([]ProgrammingLanguage, 0)

	// iterate cursor scanning lang ids and saving the to the slice created above
	for langRows.Next() {
		var lang ProgrammingLanguage
		err = langRows.Scan(&lang)
		if err != nil {
			return nil, fmt.Errorf("failed to scan lang id from link table cursor: %v", err)
		}
		languages = append(languages, lang)
	}

	// create workspace settings to unmarshall into
	var workspaceSettings *WorkspaceSettings
	if postSQL.WorkspaceSettings != nil {
		var ws WorkspaceSettings
		err = json.Unmarshal(postSQL.WorkspaceSettings, &ws)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshall workspace settings: %v", err)
		}
		workspaceSettings = &ws
	}

	// create new post for the output
	post := &Post{
		ID:                      postSQL.ID,
		Title:                   postSQL.Title,
		Description:             postSQL.Description,
		Author:                  postSQL.Author,
		AuthorID:                postSQL.AuthorID,
		CreatedAt:               postSQL.CreatedAt,
		UpdatedAt:               postSQL.UpdatedAt,
		RepoID:                  postSQL.RepoID,
		TopReply:                postSQL.TopReply,
		Tier:                    postSQL.Tier,
		Awards:                  awards,
		Coffee:                  postSQL.Coffee,
		PostType:                postSQL.PostType,
		Views:                   postSQL.Views,
		Completions:             postSQL.Completions,
		Attempts:                postSQL.Attempts,
		Languages:               languages,
		Published:               postSQL.Published,
		Visibility:              postSQL.Visibility,
		Tags:                    tags,
		ChallengeCost:           postSQL.ChallengeCost,
		StripePriceId:           postSQL.StripePriceId,
		WorkspaceConfig:         postSQL.WorkspaceConfig,
		WorkspaceConfigRevision: postSQL.WorkspaceConfigRevision,
		WorkspaceSettings:       workspaceSettings,
		Leads:                   postSQL.Leads,
		Embedded:                postSQL.Embedded,
		Deleted:                 postSQL.Deleted,
		ExclusiveDescription:    postSQL.ExclusiveDescription,
		ShareHash: postSQL.ShareHash,
	}

	return post, nil
}

func (i *Post) ToFrontend() (*PostFrontend, error) {

	// create slice to hold award ids in string form
	awards := make([]string, 0)

	// iterate award ids formatting them to string format and saving them to the above slice
	for b := range i.Awards {
		awards = append(awards, fmt.Sprintf("%d", b))
	}

	// create slice to hold tag ids in string form
	tags := make([]string, 0)

	// iterate tag ids formatting them to string format and saving them to the above slice
	for b := range i.Tags {
		tags = append(tags, fmt.Sprintf("%d", b))
	}

	// conditionally format top reply id into string
	var topReply string
	if i.TopReply != nil {
		topReply = fmt.Sprintf("%d", *i.TopReply)
	}

	// create slice to hold language strings
	langStrings := make([]string, 0)

	// iterate language ids formatting them to string format and saving them to the above slice
	for _, l := range i.Languages {
		langStrings = append(langStrings, l.String())
	}

	if i.Deleted {
		return &PostFrontend{
			ID:                      fmt.Sprintf("%d", i.ID),
			Title:                   "[Removed]",
			Description:             "This post had been removed by the original author.",
			Author:                  i.Author,
			AuthorID:                fmt.Sprintf("%d", i.AuthorID),
			CreatedAt:               i.CreatedAt,
			UpdatedAt:               i.UpdatedAt,
			RepoID:                  fmt.Sprintf("%d", i.RepoID),
			Tier:                    i.Tier,
			TierString:              i.Tier.String(),
			Awards:                  awards,
			TopReply:                &topReply,
			Coffee:                  i.Coffee,
			PostType:                i.PostType,
			PostTypeString:          i.PostType.String(),
			Views:                   i.Views,
			Attempts:                i.Attempts,
			Completions:             i.Completions,
			Languages:               i.Languages,
			LanguageStrings:         langStrings,
			Published:               i.Published,
			Visibility:              i.Visibility,
			VisibilityString:        i.Visibility.String(),
			Tags:                    tags,
			ChallengeCost:           i.ChallengeCost,
			WorkspaceConfig:         fmt.Sprintf("%d", i.WorkspaceConfig),
			WorkspaceConfigRevision: i.WorkspaceConfigRevision,
			Leads:                   i.Leads,
			Thumbnail:               fmt.Sprintf("/static/posts/t/%v", i.ID),
			Deleted:                 true,
			ExclusiveDescription:    i.ExclusiveDescription,
		}, nil
	}

	// // hash id for thumbnail path
	// idHash, err := utils.HashData([]byte(fmt.Sprintf("%d", i.ID)))
	// if err != nil {
	//	return nil, fmt.Errorf("failed to hash post id: %v", err)
	// }

	// create new post frontend
	mf := &PostFrontend{
		ID:                      fmt.Sprintf("%d", i.ID),
		Title:                   i.Title,
		Description:             i.Description,
		Author:                  i.Author,
		AuthorID:                fmt.Sprintf("%d", i.AuthorID),
		CreatedAt:               i.CreatedAt,
		UpdatedAt:               i.UpdatedAt,
		RepoID:                  fmt.Sprintf("%d", i.RepoID),
		Tier:                    i.Tier,
		TierString:              i.Tier.String(),
		Awards:                  awards,
		TopReply:                &topReply,
		Coffee:                  i.Coffee,
		PostType:                i.PostType,
		PostTypeString:          i.PostType.String(),
		Views:                   i.Views,
		Attempts:                i.Attempts,
		Completions:             i.Completions,
		Languages:               i.Languages,
		LanguageStrings:         langStrings,
		Published:               i.Published,
		Visibility:              i.Visibility,
		VisibilityString:        i.Visibility.String(),
		Tags:                    tags,
		ChallengeCost:           i.ChallengeCost,
		WorkspaceConfig:         fmt.Sprintf("%d", i.WorkspaceConfig),
		WorkspaceConfigRevision: i.WorkspaceConfigRevision,
		Leads:                   i.Leads,
		Thumbnail:               fmt.Sprintf("/static/posts/t/%v", i.ID),
		ExclusiveDescription:    i.ExclusiveDescription,
	}

	return mf, nil
}

func (i *Post) ToSQLNative() ([]*SQLInsertStatement, error) {
	// marshall workspace settings to json buffer
	var buf []byte
	if i.WorkspaceSettings != nil {
		b, err := json.Marshal(i.WorkspaceSettings)
		if err != nil {
			return nil, fmt.Errorf("failed to marshall workspace settings: %v", err)
		}
		buf = b
	}

	// create slice to hold insertion statements for this post and initialize the slice with the main post insertion statement
	sqlStatements := []*SQLInsertStatement{
		{
			Statement: "insert ignore into post(_id, title, description, author, author_id, created_at, updated_at, repo_id, top_reply, tier, coffee, post_type, views, completions, attempts, published, visibility, stripe_price_id, challenge_cost, workspace_config, workspace_config_revision, workspace_settings, leads, embedded, deleted, exclusive_description, share_hash) values(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, uuid_to_bin(?));",
			Values: []interface{}{i.ID, i.Title, i.Description, i.Author, i.AuthorID, i.CreatedAt, i.UpdatedAt, i.RepoID,
				i.TopReply, i.Tier, i.Coffee, i.PostType, i.Views, i.Completions, i.Attempts, i.Published, i.Visibility,
				i.StripePriceId, i.ChallengeCost, i.WorkspaceConfig, i.WorkspaceConfigRevision, buf, i.Leads, i.Embedded,
				i.Deleted, i.ExclusiveDescription, i.ShareHash,
			},
		},
	}

	// conditionally iterate award ids formatting them for sql insertion into the post awards link table
	if len(i.Awards) > 0 {
		for _, a := range i.Awards {
			awardStatement := SQLInsertStatement{
				Statement: "insert ignore into post_awards(post_id, award_id) values(?, ?);",
				Values:    []interface{}{i.ID, a},
			}
			sqlStatements = append(sqlStatements, &awardStatement)
		}
	}

	// conditionally iterate tag ids formatting them for sql insertion into the post tags link table
	if len(i.Tags) > 0 {
		for _, a := range i.Tags {
			tagStatement := SQLInsertStatement{
				Statement: "insert ignore into post_tags(post_id, tag_id) values(?, ?);",
				Values:    []interface{}{i.ID, a},
			}
			sqlStatements = append(sqlStatements, &tagStatement)
		}
	}

	// conditionally iterate language ids formatting them for sql insertion into the post languages link table
	if len(i.Languages) > 0 {
		for _, l := range i.Languages {
			languageStatement := SQLInsertStatement{
				Statement: "insert ignore into post_langs(post_id, lang_id) values(?,?);",
				Values:    []interface{}{i.ID, l},
			}
			sqlStatements = append(sqlStatements, &languageStatement)
		}
	}

	// create insertion statement and return
	return sqlStatements, nil
}
