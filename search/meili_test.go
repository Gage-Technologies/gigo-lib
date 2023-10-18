package search

import (
	"github.com/gage-technologies/gigo-lib/config"
	"github.com/gage-technologies/gigo-lib/db/models"
	"reflect"
	"testing"
	"time"
)

func TestCreateMeiliSearchEngine(t *testing.T) {
	cfg := config.MeiliConfig{
		Host:  "http://localhost:7700",
		Token: "MASTER_KEY",
		Indices: map[string]config.MeiliIndexConfig{
			"createTest1": {
				Name:                 "createTest1",
				PrimaryKey:           "f1",
				SearchableAttributes: []string{"f1", "f3"},
				FilterableAttributes: []string{"f2", "f4"},
				DisplayedAttributes:  []string{"f1", "f2", "f3"},
				SortableAttributes:   []string{"f1", "f2", "f3"},
			},
			"createTest2": {
				Name:                 "createTest2",
				PrimaryKey:           "f1",
				SearchableAttributes: []string{"f1"},
				FilterableAttributes: []string{"f2", "f3", "f4"},
				DisplayedAttributes:  []string{"f1", "f3"},
				SortableAttributes:   []string{"f1"},
			},
		},
	}

	meili, err := CreateMeiliSearchEngine(cfg)
	if err != nil {
		t.Fatalf("\n%s failed\n    Error: %v", t.Name(), err)
	}

	meili, err = CreateMeiliSearchEngine(cfg)
	if err != nil {
		t.Fatalf("\n%s failed\n    Error: %v", t.Name(), err)
	}

	_, _ = meili.client.DeleteIndex("index1")
	_, _ = meili.client.DeleteIndex("index2")

	t.Logf("\n%s succeeded", t.Name())
}

func TestMeiliSearchEngine_AddDocuments(t *testing.T) {
	cfg := config.MeiliConfig{
		Host:  "http://localhost:7700",
		Token: "MASTER_KEY",
		Indices: map[string]config.MeiliIndexConfig{
			"addTest": {
				Name:                 "addTest",
				PrimaryKey:           "_id",
				SearchableAttributes: []string{"title", "description", "author"},
				FilterableAttributes: []string{
					"languages",
					"attempts",
					"completions",
					"coffee",
					"views",
					"tags",
					"post_type",
					"visibility",
					"created_at",
					"updated_at",
					"published",
					"tier",
					"author_id",
				},
				SortableAttributes: []string{
					"attempts",
					"completions",
					"coffee",
					"views",
					"created_at",
					"updated_at",
					"tier",
				},
			},
		},
	}

	meili, err := CreateMeiliSearchEngine(cfg)
	if err != nil {
		t.Fatalf("\n%s failed\n    Error: %v", t.Name(), err)
	}

	defer meili.client.DeleteIndex("addTest")

	posts := []interface{}{
		models.Post{
			ID:          420,
			Title:       "test title 1",
			Description: "test description 2",
			Author:      "test author 1",
			AuthorID:    69,
			CreatedAt:   time.Now().Add(-time.Hour * 24 * 32),
			UpdatedAt:   time.Now().Add(-time.Hour * 24 * 7),
			RepoID:      42069,
			Tier:        models.Tier6,
			Awards:      []int64{1, 2, 3, 4, 5, 6},
			Coffee:      17283,
			Tags:        []int64{7, 8, 9},
			PostType:    models.CompetitiveChallenge,
			Views:       4200,
			Languages:   []models.ProgrammingLanguage{models.Go, models.JavaScript},
			Attempts:    420,
			Completions: 69,
			Published:   true,
			Visibility:  models.PublicVisibility,
		},
		models.Post{
			ID:          421,
			Title:       "test title 1",
			Description: "test description 2",
			Author:      "test author 1",
			AuthorID:    70,
			CreatedAt:   time.Now().Add(-time.Hour * 24 * 34),
			UpdatedAt:   time.Now().Add(-time.Hour * 24 * 2),
			RepoID:      42070,
			Tier:        models.Tier6,
			Awards:      []int64{1, 2},
			Coffee:      17283,
			Tags:        []int64{7, 8, 9, 10, 122},
			PostType:    models.CasualChallenge,
			Views:       421,
			Languages:   []models.ProgrammingLanguage{models.Java, models.JavaScript},
			Attempts:    450,
			Completions: 19,
			Published:   true,
			Visibility:  models.PublicVisibility,
		},
	}

	err = meili.AddDocuments("addTest", posts...)
	if err != nil {
		t.Fatalf("\n%s failed\n    Error: %v", t.Name(), err)
	}

	stats, err := meili.client.Index("addTest").GetStats()
	if err != nil {
		t.Fatalf("\n%s failed\n    Error: %v", t.Name(), err)
	}

	if stats.NumberOfDocuments != int64(len(posts)) {
		t.Fatalf("\n%s failed\n    Error: incorrect document count %d != 2", t.Name(), stats.NumberOfDocuments)
	}

	t.Logf("\n%s succeeded", t.Name())
}

func TestMeiliSearchEngine_UpdateDocuments(t *testing.T) {
	cfg := config.MeiliConfig{
		Host:  "http://localhost:7700",
		Token: "MASTER_KEY",
		Indices: map[string]config.MeiliIndexConfig{
			"updateTest": {
				Name:                 "updateTest",
				PrimaryKey:           "f1",
				SearchableAttributes: []string{"f1", "f3"},
				FilterableAttributes: []string{"f2", "f4"},
				DisplayedAttributes:  []string{"f1", "f2", "f3"},
				SortableAttributes:   []string{"f1", "f2", "f3"},
			},
		},
	}

	meili, err := CreateMeiliSearchEngine(cfg)
	if err != nil {
		t.Fatalf("\n%s failed\n    Error: %v", t.Name(), err)
	}

	defer meili.client.DeleteIndex("updateTest")

	err = meili.AddDocuments("updateTest", []interface{}{
		map[string]interface{}{
			"f1": "test1",
			"f2": "test1",
			"f3": "test1",
			"f4": "test1",
		},
		map[string]interface{}{
			"f1": "test2",
			"f2": "test2",
			"f3": "test2",
			"f4": "test2",
		},
	}...)
	if err != nil {
		t.Fatalf("\n%s failed\n    Error: %v", t.Name(), err)
	}

	err = meili.UpdateDocuments("updateTest", []interface{}{
		map[string]interface{}{
			"f1": "test1",
			"f2": "test3",
			"f3": "test3",
			"f4": "test3",
		},
		map[string]interface{}{
			"f1": "test2",
			"f2": "test4",
			"f3": "test4",
			"f4": "test4",
		},
	}...)
	if err != nil {
		t.Fatalf("\n%s failed\n    Error: %v", t.Name(), err)
	}

	var res map[string]interface{}
	err = meili.client.Index("updateTest").GetDocument("test1", nil, &res)
	if err != nil {
		t.Fatalf("\n%s failed\n    Error: %v", t.Name(), err)
	}

	if !reflect.DeepEqual(res, map[string]interface{}{
		"f1": "test1",
		"f2": "test3",
		"f3": "test3",
		"f4": "test3",
	}) {
		t.Fatalf("\n%s failed\n    Expected: %#v\n    Actual: %#v", t.Name(), map[string]interface{}{
			"f1": "test1",
			"f2": "test3",
			"f3": "test3",
			"f4": "test3",
		}, res)
	}

	res = nil
	err = meili.client.Index("updateTest").GetDocument("test2", nil, &res)
	if err != nil {
		t.Fatalf("\n%s failed\n    Error: %v", t.Name(), err)
	}

	if !reflect.DeepEqual(res, map[string]interface{}{
		"f1": "test2",
		"f2": "test4",
		"f3": "test4",
		"f4": "test4",
	}) {
		t.Fatalf("\n%s failed\n    Expected: %#v\n    Actual: %#v", t.Name(), map[string]interface{}{
			"f1": "test2",
			"f2": "test4",
			"f3": "test4",
			"f4": "test4",
		}, res)
	}

	t.Logf("\n%s succeeded", t.Name())
}

func TestMeiliSearchEngine_DeleteDocuments(t *testing.T) {
	cfg := config.MeiliConfig{
		Host:  "http://localhost:7700",
		Token: "MASTER_KEY",
		Indices: map[string]config.MeiliIndexConfig{
			"deleteTest": {
				Name:                 "deleteTest",
				PrimaryKey:           "f1",
				SearchableAttributes: []string{"f1", "f3"},
				FilterableAttributes: []string{"f2", "f4"},
				DisplayedAttributes:  []string{"f1", "f2", "f3"},
				SortableAttributes:   []string{"f1", "f2", "f3"},
			},
		},
	}

	meili, err := CreateMeiliSearchEngine(cfg)
	if err != nil {
		t.Fatalf("\n%s failed\n    Error: %v", t.Name(), err)
	}

	defer meili.client.DeleteIndex("deleteTest")

	err = meili.AddDocuments("deleteTest", []interface{}{
		map[string]interface{}{
			"f1": "test1",
			"f2": "test1",
			"f3": "test1",
			"f4": "test1",
		},
		map[string]interface{}{
			"f1": "test2",
			"f2": "test2",
			"f3": "test2",
			"f4": "test2",
		},
	}...)
	if err != nil {
		t.Fatalf("\n%s failed\n    Error: %v", t.Name(), err)
	}

	err = meili.DeleteDocuments("deleteTest", "test1", "test2")
	if err != nil {
		t.Fatalf("\n%s failed\n    Error: %v", t.Name(), err)
	}

	stats, err := meili.client.Index("deleteTest").GetStats()
	if err != nil {
		t.Fatalf("\n%s failed\n    Error: %v", t.Name(), err)
	}

	if stats.NumberOfDocuments > 0 {
		t.Fatalf("\n%s failed\n    Error: failed to delete documents %d remaining", t.Name(), stats.NumberOfDocuments)
	}

	t.Logf("\n%s succeeded", t.Name())
}

func TestMeiliSearchEngine_Search(t *testing.T) {
	cfg := config.MeiliConfig{
		Host:  "http://localhost:7700",
		Token: "MASTER_KEY",
		Indices: map[string]config.MeiliIndexConfig{
			"searchTest": {
				Name:                 "searchTest",
				PrimaryKey:           "_id",
				SearchableAttributes: []string{"title", "description", "author"},
				FilterableAttributes: []string{
					"languages",
					"attempts",
					"completions",
					"coffee",
					"views",
					"tags",
					"post_type",
					"visibility",
					"created_at",
					"updated_at",
					"published",
					"tier",
					"author_id",
				},
				SortableAttributes: []string{
					"attempts",
					"completions",
					"coffee",
					"views",
					"created_at",
					"updated_at",
					"tier",
				},
			},
		},
	}

	meili, err := CreateMeiliSearchEngine(cfg)
	if err != nil {
		t.Fatalf("\n%s failed\n    Error: %v", t.Name(), err)
	}

	defer meili.client.DeleteIndex("searchTest")

	posts := []interface{}{
		models.Post{
			ID:          420,
			Title:       "building a fullstack developer platform using Go and Javascript",
			Description: "Use Go and Javascript to build a fullstack developer platform including a backend HTTP server and frontend React app",
			Author:      "LearnGo",
			AuthorID:    69,
			CreatedAt:   time.Now().Add(-time.Hour * 24 * 32),
			UpdatedAt:   time.Now().Add(-time.Hour * 24 * 7),
			RepoID:      42069,
			Tier:        models.Tier6,
			Awards:      []int64{1, 2, 3, 4, 5, 6},
			Coffee:      17283,
			Tags:        []int64{7, 8, 9},
			PostType:    models.CompetitiveChallenge,
			Views:       4200,
			Languages:   []models.ProgrammingLanguage{models.Go, models.JavaScript},
			Attempts:    420,
			Completions: 69,
			Published:   true,
			Visibility:  models.PublicVisibility,
		},
		models.Post{
			ID:          421,
			Title:       "Advanced Machine Learning Course",
			Description: "Learn to develop custom C++ extensions for PyTorch and train custom machine learning models using your extensions",
			Author:      "MachineLearningGuy",
			AuthorID:    70,
			CreatedAt:   time.Now().Add(-time.Hour * 24 * 34),
			UpdatedAt:   time.Now().Add(-time.Hour * 24 * 2),
			RepoID:      42070,
			Tier:        models.Tier6,
			Awards:      []int64{1, 2},
			Coffee:      17283,
			Tags:        []int64{7, 8, 9, 10, 122},
			PostType:    models.CasualChallenge,
			Views:       421,
			Languages:   []models.ProgrammingLanguage{models.Python, models.Cpp},
			Attempts:    450,
			Completions: 19,
			Published:   true,
			Visibility:  models.PublicVisibility,
		},
	}

	err = meili.AddDocuments("searchTest", posts...)
	if err != nil {
		t.Fatalf("\n%s failed\n    Error: %v", t.Name(), err)
	}

	lang1 := models.Go
	lang2 := models.Java

	res, err := meili.Search("searchTest", &Request{
		Query: "fullstack",
		Filter: &FilterGroup{
			Filters: []FilterCondition{
				{
					Filters: []Filter{
						{
							Attribute: "languages",
							Operator:  OperatorIn,
							Values:    []interface{}{&lang1, &lang2},
						},
						{
							Attribute: "visibility",
							Operator:  OperatorEquals,
							Value:     models.PublicVisibility,
						},
					},
					And: true,
				},
			},
		},
	})
	if err != nil {
		t.Fatalf("\n%s failed\n    Error: %v", t.Name(), err)
	}

	if res.TotalResults != 1 {
		t.Fatalf("\n%s failed\n    Error: incorrect hit count %v", t.Name(), len(res.Hits))
	}

	ok, err := res.Next()
	if err != nil {
		t.Fatalf("\n%s failed\n    Error: %v", t.Name(), err)
	}
	if !ok {
		t.Fatalf("\n%s failed\n    Error: failed to load result into cursor first position", t.Name())
	}

	var p models.Post
	err = res.Scan(&p)
	if err != nil {
		t.Fatalf("\n%s failed\n    Error: failed to scan post from result %v", t.Name(), err)
	}

	if posts[0].(models.Post).ID != p.ID {
		t.Fatalf("\n%s failed\n    Error: incorrect post %+v", t.Name(), p)
	}

	t.Logf("\n%s succeeded", t.Name())
}
