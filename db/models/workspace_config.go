package models

import (
	"context"
	"database/sql"
	"fmt"
	ti "github.com/gage-technologies/gigo-lib/db"
	"github.com/kisielk/sqlstruct"
	"go.opentelemetry.io/otel"
)

type WorkspaceConfig struct {
	ID          int64                 `json:"_id" sql:"_id"`
	Title       string                `json:"title" sql:"title"`
	Description string                `json:"description" sql:"description"`
	Content     string                `json:"content" sql:"content"`
	AuthorID    int64                 `json:"author_id" sql:"author_id"`
	Revision    int                   `json:"revision" sql:"revision"`
	Official    bool                  `json:"official" sql:"official"`
	Tags        []int64               `json:"tags" sql:"tags"`
	Languages   []ProgrammingLanguage `json:"languages" sql:"languages"`
	Uses        int                   `json:"uses" sql:"uses"`
	Completions int                   `json:"completions" sql:"completions"`
}

type WorkspaceConfigSQL struct {
	ID          int64  `json:"_id" sql:"_id"`
	Title       string `json:"title" sql:"title"`
	Description string `json:"description" sql:"description"`
	Content     string `json:"content" sql:"content"`
	AuthorID    int64  `json:"author_id" sql:"author_id"`
	Revision    int    `json:"revision" sql:"revision"`
	Official    bool   `json:"official" sql:"official"`
	Uses        int    `json:"uses" sql:"uses"`
	Completions int    `json:"completions" sql:"completions"`
}

type WorkspaceConfigFrontend struct {
	ID              string                `json:"_id"`
	Title           string                `json:"title" sql:"title"`
	Description     string                `json:"description"`
	Content         string                `json:"content"`
	AuthorID        string                `json:"author_id"`
	Author          string                `json:"author"`
	Revision        int                   `json:"revision"`
	Official        bool                  `json:"official"`
	Tags            []string              `json:"tags"`
	Languages       []ProgrammingLanguage `json:"languages"`
	LanguageStrings []string              `json:"languages_strings"`
	Uses            int                   `json:"uses" sql:"uses"`
	Completions     int                   `json:"completions" sql:"completions"`
}

func CreateWorkspaceConfig(_id int64, title string, description string, content string, authorID int64, revision int,
	tags []int64, languages []ProgrammingLanguage, uses int) *WorkspaceConfig {
	return &WorkspaceConfig{
		ID:          _id,
		Title:       title,
		Description: description,
		Content:     content,
		AuthorID:    authorID,
		Revision:    revision,
		Official:    false,
		Tags:        tags,
		Languages:   languages,
		Uses:        uses,
		Completions: 0,
	}
}

func WorkspaceConfigFromSQLNative(db *ti.Database, rows *sql.Rows) (*WorkspaceConfig, error) {
	// create new config object to load into
	var configSQL WorkspaceConfigSQL

	// scan row into config object
	err := sqlstruct.Scan(&configSQL, rows)
	if err != nil {
		return nil, err
	}

	ctx, span := otel.Tracer("gigo-core").Start(context.Background(), "gigo-lib")
	callerName := "WorkspaceConfigFromSQLNative"
	// query tag link table to get tab ids
	tagRows, err := db.QueryContext(ctx, &span, &callerName, "select tag_id from workspace_config_tags where cfg_id = ? and revision = ?", configSQL.ID, configSQL.Revision)
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
			return nil, fmt.Errorf("failed to scan tag id from link table cursor: %v", err)
		}
		tags = append(tags, tag)
	}

	ctx, span = otel.Tracer("gigo-core").Start(context.Background(), "gigo-lib")
	// query lang link table to get lang ids
	langRows, err := db.QueryContext(ctx, &span, &callerName, "select lang_id from workspace_config_langs where cfg_id = ? and revision = ?", configSQL.ID, configSQL.Revision)
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

	return &WorkspaceConfig{
		ID:          configSQL.ID,
		Title:       configSQL.Title,
		Description: configSQL.Description,
		Content:     configSQL.Content,
		AuthorID:    configSQL.AuthorID,
		Revision:    configSQL.Revision,
		Official:    configSQL.Official,
		Tags:        tags,
		Languages:   languages,
		Uses:        configSQL.Uses,
		Completions: configSQL.Completions,
	}, nil
}

func (c *WorkspaceConfig) ToFrontend() *WorkspaceConfigFrontend {
	// create slice to hold tag ids in string form
	tags := make([]string, 0)

	// iterate tag ids formatting them to string format and saving them to the above slice
	for b := range c.Tags {
		tags = append(tags, fmt.Sprintf("%d", b))
	}

	// create slice to hold language strings
	langStrings := make([]string, 0)

	// iterate language ids formatting them to string format and saving them to the above slice
	for _, l := range c.Languages {
		langStrings = append(langStrings, l.String())
	}

	return &WorkspaceConfigFrontend{
		ID:              fmt.Sprintf("%d", c.ID),
		Title:           c.Title,
		Description:     c.Description,
		Content:         c.Content,
		AuthorID:        fmt.Sprintf("%d", c.AuthorID),
		Revision:        c.Revision,
		Official:        c.Official,
		Tags:            tags,
		Languages:       c.Languages,
		LanguageStrings: langStrings,
		Uses:            c.Uses,
		Completions:     c.Completions,
	}
}

func (c *WorkspaceConfig) ToSQLNative() ([]*SQLInsertStatement, error) {
	// create slice to hold insertion statements for this workspace config and initialize the slice with the main insertion statement
	sqlStatements := []*SQLInsertStatement{
		{
			Statement: "insert ignore into workspace_config(_id, title, description, content, author_id, revision, official, uses, completions) values(?, ?, ?, ?, ?, ?, ?, ?, ?);",
			Values: []interface{}{
				c.ID,
				c.Title,
				c.Description,
				c.Content,
				c.AuthorID,
				c.Revision,
				c.Official,
				c.Uses,
				c.Completions,
			},
		},
	}

	// conditionally iterate tag ids formatting them for sql insertion into the post tags link table
	if len(c.Tags) > 0 {
		for _, a := range c.Tags {
			tagStatement := SQLInsertStatement{
				Statement: "insert ignore into workspace_config_tags(cfg_id, tag_id, revision) values(?, ?, ?);",
				Values:    []interface{}{c.ID, a, c.Revision},
			}
			sqlStatements = append(sqlStatements, &tagStatement)
		}
	}

	// conditionally iterate language ids formatting them for sql insertion into the post languages link table
	if len(c.Languages) > 0 {
		for _, l := range c.Languages {
			languageStatement := SQLInsertStatement{
				Statement: "insert ignore into workspace_config_langs(cfg_id, lang_id, revision) values(?, ?, ?);",
				Values:    []interface{}{c.ID, l, c.Revision},
			}
			sqlStatements = append(sqlStatements, &languageStatement)
		}
	}

	return sqlStatements, nil
}
