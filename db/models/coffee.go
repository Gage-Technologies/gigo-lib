package models

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/kisielk/sqlstruct"
	"time"
)

type Coffee struct {
	ID           int64     `json:"_id" sql:"_id"`
	UserID       int64     `json:"user_id" sql:"user_id"`
	CreatedAt    time.Time `json:"created_at" sql:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" sql:"updated_at"`
	AttemptID    *int64    `json:"attempt_id" sql:"attempt_id"`
	PostID       *int64    `json:"post_id" sql:"post_id"`
	DiscussionID *int64    `json:"discussion_id" sql:"discussion_id"`
}

type CoffeeSQL struct {
	ID           int64     `json:"_id" sql:"_id"`
	UserID       int64     `json:"user_id" sql:"user_id"`
	CreatedAt    time.Time `json:"created_at" sql:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" sql:"updated_at"`
	AttemptID    *int64    `json:"attempt_id" sql:"attempt_id"`
	PostID       *int64    `json:"post_id" sql:"post_id"`
	DiscussionID *int64    `json:"discussion_id" sql:"discussion_id"`
}

type CoffeeFrontend struct {
	ID           string    `json:"_id" sql:"_id"`
	UserID       string    `json:"user_id" sql:"user_id"`
	CreatedAt    time.Time `json:"created_at" sql:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" sql:"updated_at"`
	AttemptID    *string   `json:"attempt_id" sql:"attempt_id"`
	PostID       *string   `json:"post_id" sql:"post_id"`
	DiscussionID *string   `json:"discussion_id" sql:"discussion_id"`
}

func CreateCoffee(id int64, userID int64, createdAt time.Time, updatedAt time.Time, attemptID *int64, postID *int64,
	discussionID *int64) (*Coffee, error) {

	if attemptID == nil && postID == nil && discussionID == nil {
		return nil, errors.New("failed to create coffee, err: no post, attempt or discussion found")
	}

	return &Coffee{
		ID:           id,
		UserID:       userID,
		CreatedAt:    createdAt,
		UpdatedAt:    updatedAt,
		AttemptID:    attemptID,
		PostID:       postID,
		DiscussionID: discussionID,
	}, nil
}

func CoffeeFromSQLNative(rows *sql.Rows) (*Coffee, error) {
	// create new coffee object to load into
	coffeeSQL := new(CoffeeSQL)

	// scan row into coffee object
	err := sqlstruct.Scan(coffeeSQL, rows)
	if err != nil {
		return nil, err
	}

	// create new coffee for the output
	coffee := &Coffee{
		ID:           coffeeSQL.ID,
		CreatedAt:    coffeeSQL.CreatedAt,
		UpdatedAt:    coffeeSQL.UpdatedAt,
		UserID:       coffeeSQL.UserID,
		PostID:       coffeeSQL.PostID,
		AttemptID:    coffeeSQL.AttemptID,
		DiscussionID: coffeeSQL.DiscussionID,
	}

	return coffee, nil
}

func (i *Coffee) ToFrontend() *CoffeeFrontend {
	var attemptID string
	var postID string
	var discussionID string
	if i.AttemptID != nil {
		attemptID = fmt.Sprintf("%d", i.AttemptID)
	}

	if i.PostID != nil {
		postID = fmt.Sprintf("%d", i.PostID)
	}

	if i.DiscussionID != nil {
		discussionID = fmt.Sprintf("%d", i.DiscussionID)
	}

	// create new coffee frontend
	mf := &CoffeeFrontend{
		ID:           fmt.Sprintf("%d", i.ID),
		CreatedAt:    i.CreatedAt,
		UpdatedAt:    i.UpdatedAt,
		UserID:       fmt.Sprintf("%d", i.UserID),
		PostID:       &postID,
		AttemptID:    &attemptID,
		DiscussionID: &discussionID,
	}

	return mf
}

func (i *Coffee) ToSQLNative() *SQLInsertStatement {

	// create insertion statement and return
	return &SQLInsertStatement{
		Statement: "insert ignore into coffee(_id, created_at, updated_at, post_id, attempt_id, user_id, discussion_id) values(?, ?, ?, ?, ?, ?, ?);",
		Values:    []interface{}{i.ID, i.CreatedAt, i.UpdatedAt, i.PostID, i.AttemptID, i.UserID, i.DiscussionID},
	}
}
