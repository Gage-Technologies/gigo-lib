package models

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/gage-technologies/gigo-lib/db"
	"github.com/kisielk/sqlstruct"
	"go.opentelemetry.io/otel"
	"time"
)

type SearchRec struct {
	ID               int64     `sql:"_id" json:"_id"`
	UserID           int64     `sql:"user_id" json:"user_id"`
	Query            string    `sql:"query" json:"query"`
	PostIDs          []int64   `sql:"post_ids" json:"post_ids"`
	SelectedPostID   *int64    `sql:"selected_post_id" json:"selected_post_id"`
	SelectedPostName *string   `sql:"selected_post_name" json:"selected_post_name"`
	CreatedAt        time.Time `sql:"created_at" json:"created_at"`
}

type SearchRecSQL struct {
	ID               int64     `sql:"_id" json:"_id" json:"_id"`
	UserID           int64     `sql:"user_id" json:"user_id"`
	Query            string    `sql:"query" json:"query"`
	SelectedPostID   *int64    `sql:"selected_post_id" json:"selected_post_id"`
	SelectedPostName *string   `sql:"selected_post_name" json:"selected_post_name"`
	CreatedAt        time.Time `sql:"created_at" json:"created_at"`
}

type SearchRecFrontend struct {
	ID               string    `sql:"_id" json:"_id"`
	UserID           string    `sql:"user_id" json:"user_id"`
	Query            string    `sql:"query" json:"query"`
	PostIDs          []string  `sql:"post_ids" json:"post_ids"`
	SelectedPostID   *string   `sql:"selected_post_id" json:"selected_post_id"`
	SelectedPostName *string   `sql:"selected_post_name" json:"selected_post_name"`
	CreatedAt        time.Time `sql:"created_at" json:"created_at"`
}

func CreateSearchRec(id int64, userId int64, postIds []int64, query string, selectedPostId *int64, selectedPostName *string, createdAt time.Time) *SearchRec {
	return &SearchRec{
		ID:               id,
		UserID:           userId,
		PostIDs:          postIds,
		Query:            query,
		SelectedPostID:   selectedPostId,
		SelectedPostName: selectedPostName,
		CreatedAt:        createdAt,
	}
}

func SearchRecFromSQLNative(db *ti.Database, rows *sql.Rows) (*SearchRec, error) {
	// create new Search post object to load into
	recSql := new(SearchRecSQL)

	// scan row into Search post object
	err := sqlstruct.Scan(recSql, rows)
	if err != nil {
		return nil, err
	}

	ctx, span := otel.Tracer("gigo-core").Start(context.Background(), "gigo-lib")
	callerName := "SearchRecFromSQLNative"
	res, err := db.QueryContext(ctx, &span, &callerName, "select post_id from search_rec_posts where user_id = ? and search_id", recSql.UserID, recSql.ID)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed to query search_rec_posts: %v", err))
	}

	postIds := make([]int64, 0)

	for res.Next() {
		var postId int64
		err = res.Scan(&postId)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("failed to scan search_rec_posts: %v", err))
		}

		postIds = append(postIds, postId)
	}

	return &SearchRec{
		ID:               recSql.ID,
		UserID:           recSql.UserID,
		PostIDs:          postIds,
		Query:            recSql.Query,
		SelectedPostID:   recSql.SelectedPostID,
		SelectedPostName: recSql.SelectedPostName,
		CreatedAt:        recSql.CreatedAt,
	}, nil
}

func (i *SearchRec) ToFrontend() *SearchRecFrontend {
	postIds := make([]string, 0)
	for _, p := range i.PostIDs {
		postIds = append(postIds, fmt.Sprintf("%d", p))
	}

	// create new Search post frontend
	mf := &SearchRecFrontend{
		ID:        fmt.Sprintf("%d", i.ID),
		UserID:    fmt.Sprintf("%d", i.UserID),
		PostIDs:   postIds,
		Query:     i.Query,
		CreatedAt: i.CreatedAt,
	}

	return mf
}

func (i *SearchRec) ToSQLNative() []*SQLInsertStatement {

	statements := make([]*SQLInsertStatement, 0)

	statements = append(statements, &SQLInsertStatement{
		Statement: "insert into search_rec(_id, user_id, query, selected_post_id, selected_post_name, created_at) values (?, ?, ?, ?, ?, ?);",
		Values:    []interface{}{i.ID, i.UserID, i.Query, i.SelectedPostID, i.SelectedPostName, i.CreatedAt},
	})

	for _, p := range i.PostIDs {
		statements = append(statements, &SQLInsertStatement{
			Statement: "insert into search_rec_posts(search_id, post_id) values (?, ?);",
			Values:    []interface{}{i.ID, p},
		})
	}

	// create insertion statement and return
	return statements
}
