package search

import (
	"encoding/json"
	"fmt"
	"github.com/buger/jsonparser"
	"github.com/gage-technologies/gigo-lib/config"
	"github.com/meilisearch/meilisearch-go"
	"strings"
	"time"
)

// MeiliSearchEngine
//
//	Struct for interacting with the Meilisearch Search Engine
//	for the GIGO codebase
type MeiliSearchEngine struct {
	client *meilisearch.Client
}

// CreateMeiliSearchEngine
//
//	Create a new MeiliSearchEngine and initializes the configured indexes
func CreateMeiliSearchEngine(config config.MeiliConfig) (*MeiliSearchEngine, error) {
	// create a new client
	client := meilisearch.NewClient(meilisearch.ClientConfig{
		Host:    config.Host,
		APIKey:  config.Token,
		Timeout: time.Second * 3,
	})

	// iterate the indices config initializing the indices and configuring them
	for _, indexConfig := range config.Indices {
		// create bools to track whether we are creating the index and whether we are configuring it
		createIndex := false
		configureIndex := indexConfig.UpdateConfig

		// create index if it doesn't exist
		_, err := client.GetIndex(indexConfig.Name)
		if err != nil {
			// return error if this is an unexpected error
			if !strings.Contains(err.Error(), "index_not_found") {
				return nil, fmt.Errorf("failed to retrieve index %q: %v", indexConfig.Name, err)
			}

			// mark index for creation and configuration since it doesn't exist
			createIndex = true
			configureIndex = true
		}

		// conditionally create index
		if createIndex {
			// execute index creation task
			indexCreationTask, err := client.CreateIndex(&meilisearch.IndexConfig{
				Uid:        indexConfig.Name,
				PrimaryKey: indexConfig.PrimaryKey,
			})
			if err != nil {
				return nil, fmt.Errorf("failed to create index creation task for %q: %v", indexConfig.Name, err)
			}

			// wait for index creation to complete
			task, err := client.WaitForTask(indexCreationTask.TaskUID)
			if err != nil {
				return nil, fmt.Errorf("failed to wait for index creation task for %q: %v", indexConfig.Name, err)
			}

			// check if there was an error creating the index
			if task.Error.Code != "" {
				return nil, fmt.Errorf("failed to create index %q: %+v", indexConfig.Name, task.Error)
			}
		}

		// conditionally configure index
		if configureIndex {
			// set distinct attribute
			t, err := client.Index(indexConfig.Name).UpdateDistinctAttribute(indexConfig.PrimaryKey)
			if err != nil {
				return nil, fmt.Errorf("failed to update distinct attribute for index %q: %v", indexConfig.Name, err)
			}

			// wait for configuration of distinct attribute to complete
			task, err := client.WaitForTask(t.TaskUID)
			if err != nil {
				return nil, fmt.Errorf("failed to wait for distinct attribute configuration task for %q: %v", indexConfig.Name, err)
			}

			// check if there was an error configuring the distinct attribute
			if task.Error.Code != "" {
				return nil, fmt.Errorf("failed to configure distinct attribute %q: %+v", indexConfig.Name, task.Error)
			}

			// set searchable attributes
			t, err = client.Index(indexConfig.Name).UpdateSearchableAttributes(&indexConfig.SearchableAttributes)
			if err != nil {
				return nil, fmt.Errorf("failed to update searchable attributes for index %q: %v", indexConfig.Name, err)
			}

			// wait for configuration of searchable attributes to complete
			task, err = client.WaitForTask(t.TaskUID)
			if err != nil {
				return nil, fmt.Errorf("failed to wait for searchable attributes configuration task for %q: %v", indexConfig.Name, err)
			}

			// check if there was an error configuring the searchable attributes
			if task.Error.Code != "" {
				return nil, fmt.Errorf("failed to configure searchable attributes %q: %+v", indexConfig.Name, task.Error)
			}

			// set filterable attributes
			t, err = client.Index(indexConfig.Name).UpdateFilterableAttributes(&indexConfig.FilterableAttributes)
			if err != nil {
				return nil, fmt.Errorf("failed to update filterable attributes for index %q: %v", indexConfig.Name, err)
			}

			// wait for configuration of filterable attributes to complete
			task, err = client.WaitForTask(t.TaskUID)
			if err != nil {
				return nil, fmt.Errorf("failed to wait for filterable attributes configuration task for %q: %v", indexConfig.Name, err)
			}

			// check if there was an error configuring the filterable attributes
			if task.Error.Code != "" {
				return nil, fmt.Errorf("failed to configure filterable attributes %q: %+v", indexConfig.Name, task.Error)
			}

			// set displayed attributes
			t, err = client.Index(indexConfig.Name).UpdateDisplayedAttributes(&indexConfig.DisplayedAttributes)
			if err != nil {
				return nil, fmt.Errorf("failed to update displayed attributes for index %q: %v", indexConfig.Name, err)
			}

			// wait for configuration of displayed attributes to complete
			task, err = client.WaitForTask(t.TaskUID)
			if err != nil {
				return nil, fmt.Errorf("failed to wait for displayed attributes configuration task for %q: %v", indexConfig.Name, err)
			}

			// check if there was an error configuring the displayed attributes
			if task.Error.Code != "" {
				return nil, fmt.Errorf("failed to configure displayed attributes %q: %+v", indexConfig.Name, task.Error)
			}

			// set sortable attributes
			t, err = client.Index(indexConfig.Name).UpdateSortableAttributes(&indexConfig.SortableAttributes)
			if err != nil {
				return nil, fmt.Errorf("failed to update sortable attributes for index %q: %v", indexConfig.Name, err)
			}

			// wait for configuration of sortable attributes to complete
			task, err = client.WaitForTask(t.TaskUID)
			if err != nil {
				return nil, fmt.Errorf("failed to wait for sortable attributes configuration task for %q: %v", indexConfig.Name, err)
			}

			// check if there was an error configuring the sortable attributes
			if task.Error.Code != "" {
				return nil, fmt.Errorf("failed to configure sortable attributes %q: %+v", indexConfig.Name, task.Error)
			}
		}
	}

	return &MeiliSearchEngine{
		client: client,
	}, nil
}

// AddDocuments
//
//	 Adds a set of documents to the configured index
//
//	 Args:
//		    - index (string): Index that documents will be added to
//		    - documents (variable interface{}): Documents that will be inserted into the index - all structs passed to
//	                                         this function should have json serializer decorators
func (e *MeiliSearchEngine) AddDocuments(index string, documents ...interface{}) error {
	// return for empty documents
	if len(documents) == 0 {
		return nil
	}

	// execute document insertion
	_, err := e.client.Index(index).AddDocuments(documents)
	if err != nil {
		return fmt.Errorf("failed to add documents to index %q: %v", index, err)
	}

	//// wait for task to complete
	//task, err := e.client.WaitForTask(res.TaskUID)
	//if err != nil {
	//	return fmt.Errorf("failed to wait for task %q: %v", res.TaskUID, err)
	//}
	//
	//// check if there was an error inserting into the index
	//if task.Error.Code != "" {
	//	return fmt.Errorf("failed to add documents to index %q: %+v", index, task.Error)
	//}

	return nil
}

// UpdateDocuments
//
//	Updates a set of documents in the configured index
//
//	Args:
//	      - index (string): Index that documents will be updated in
//	      - documents (variable interface{}): Documents that will be updated in the index - all structs passed to
//	                                          this function should have json serializer decorators
func (e *MeiliSearchEngine) UpdateDocuments(index string, documents ...interface{}) error {
	// return for empty documents
	if len(documents) == 0 {
		return nil
	}

	// execute document update
	_, err := e.client.Index(index).UpdateDocuments(documents)
	if err != nil {
		return fmt.Errorf("failed to update documents in index %q: %v", index, err)
	}

	//// wait for task to complete
	//task, err := e.client.WaitForTask(res.TaskUID)
	//if err != nil {
	//	return fmt.Errorf("failed to wait for task %q: %v", res.TaskUID, err)
	//}
	//
	//// check if there was an error updating the index
	//if task.Error.Code != "" {
	//	return fmt.Errorf("failed to update documents in index %q: %+v", index, task.Error)
	//}

	return nil
}

// DeleteDocuments
//
//	Deletes a set of documents from the configured index
//
//	Args:
//	      - index (string): Index that documents will be deleted from
//	      - ids (interface{}): Ids of documents that will be deleted from the index
func (e *MeiliSearchEngine) DeleteDocuments(index string, ids ...interface{}) error {
	// return for empty ids
	if len(ids) == 0 {
		return nil
	}

	// create slice to hold strings
	idStrings := make([]string, 0)

	// iterate over ids formatting them to strings
	for _, id := range ids {
		idStrings = append(idStrings, formatInterfaceToString(id, false))
	}

	// execute document deletion
	_, err := e.client.Index(index).DeleteDocuments(idStrings)
	if err != nil {
		return fmt.Errorf("failed to delete documents from index %q: %v", index, err)
	}

	//// wait for task to complete
	//task, err := e.client.WaitForTask(res.TaskUID)
	//if err != nil {
	//	return fmt.Errorf("failed to wait for task %q: %v", res.TaskUID, err)
	//}
	//
	//// check if there was an error deleting the index
	//if task.Error.Code != "" {
	//	return fmt.Errorf("failed to delete documents from index %q: %+v", index, task.Error)
	//}

	return nil
}

// Search
//
//		Searches the configured index in meilisearch using the configured
//	 parameters from Request
//
//	 Args:
//	     index (string): The name of the index to search
//	     req (*Request): Configuration for the search operation
//
//	 Returns:
//	     (*Result): Result of the search operation
func (e *MeiliSearchEngine) Search(index string, req *Request) (*Result, error) {
	// initialize an empty search request using system defaults
	searchRequest := &meilisearch.SearchRequest{
		Offset: 0,
		Limit:  25,
	}

	// conditionally add offset to search request
	if req.Offset != 0 {
		searchRequest.Offset = int64(req.Offset)
	}

	// conditionally add limit to search request
	if req.Limit != 0 {
		searchRequest.Limit = int64(req.Limit)
	}

	// conditionally add filter to search request
	if req.Filter != nil {
		searchRequest.Filter = req.Filter.String()
	}

	// conditionally add sort to search request
	if req.Sort != nil {
		searchRequest.Sort = req.Sort.Format()
	}

	// conditionally add facets to search request
	if req.Facets != nil {
		searchRequest.Facets = req.Facets
	}

	// execute search
	res, err := e.client.Index(index).SearchRaw(req.Query, searchRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to search index %q: %v", index, err)
	}

	// ensure that the response is not nil
	if res == nil {
		return nil, fmt.Errorf("failed to search index %q: response is nil", index)
	}

	// extract hits from search response because go-meilisearch is broken and truncates int64 on unmarshall
	hitsBytes, _, _, err := jsonparser.Get(*res, "hits")
	if err != nil {
		return nil, fmt.Errorf("failed to extract hits from search response: %v", err)
	}

	// marshall the rest of the body into a search response
	var searchResponse meilisearch.SearchResponse
	err = json.Unmarshal(*res, &searchResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal search response: %v", err)
	}

	// save total results length
	totalResults := len(searchResponse.Hits)

	// set hits in search response to nil to release memory since we can't trust that data anyway
	searchResponse.Hits = nil

	// wrap response and return
	return &Result{
		SearchResponse: &searchResponse,
		HitsBuffer:     hitsBytes,
		TotalResults:   totalResults,
	}, nil
}
