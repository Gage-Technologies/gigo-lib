package search

import (
	"encoding/json"
	"fmt"
	"github.com/buger/jsonparser"
	"github.com/meilisearch/meilisearch-go"
)

// Result
//
//	Struct used to wrap the meilisearch.SearchResponse
//	This struct includes convenience functions is Scan
//	and Next. This enables an interaction with the search
//	response that is more similar to a SQL database but
//	preserves all underlying data and functionality
type Result struct {
	*meilisearch.SearchResponse
	ResultBuffer []byte
	CursorIndex  int
	HitsBuffer   []byte
	TotalResults int
}

// Next
//
//	Mimics the sql.Rows.Next function which loads the next
//	result into the first position of cursor so that the
//	result can be scanned via the Scan function.
func (r *Result) Next() (bool, error) {
	// return false if there are no more results
	if r.CursorIndex >= r.TotalResults {
		return false, nil
	}

	// load the next result into the result buffer
	buf, _, _, err := jsonparser.Get(r.HitsBuffer, fmt.Sprintf("[%d]", r.CursorIndex))
	if err != nil {
		return false, fmt.Errorf("failed to extract result from buffer: %v", err)
	}
	r.ResultBuffer = buf

	// increment the cursor index
	r.CursorIndex++

	return true, nil
}

// Scan
//
//	Mimics the sql.Rows.Scan function which scans the
//	current buffered result into the passed interface.
//	Unlike the sql.Rows.Next function, this function does
//	not scan into individual parameters but instead scans
//	into entire structs. If an arbitrary set of values need
//	to be scanned tou can use a map[string]interface{} or
//	define a struct locally for the scan operation.
//
//	Args:
//	    dst (interface{}): pointer to the struct that the result will be scanned into
func (r *Result) Scan(dst interface{}) error {
	// ensure that there is a value in the result buffer
	if r.ResultBuffer == nil {
		return fmt.Errorf("no result in bufffer - call Next()")
	}

	// unmarshall json bytes into destination
	err := json.Unmarshal(r.ResultBuffer, dst)
	if err != nil {
		return fmt.Errorf("failed to unmarshall result buffer: %v", err)
	}

	return nil
}
