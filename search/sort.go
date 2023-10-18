package search

import "fmt"

// Sort
//
//	Structures the sorting logic and API format in a programmatic
//	interface so that sorting can be performed by developers
//	unfamiliar with the meilisearch API.
type Sort struct {
	Attribute string
	Desc      bool
}

// String
//
//	Formats the Sort into a string sort
//	supported by the meilisearch search engine
func (s *Sort) String() string {
	// select order of sort default to ascending
	order := "asc"
	if s.Desc {
		order = "desc"
	}
	return fmt.Sprintf("%s:%s", s.Attribute, order)
}

// SortGroup
//
//				Merges multiple sorts together allowing for more complex sorting
//				operations. Sort ordering is prioritized in order of their position
//	         in the Sorts slice. Lower indices correlate to a higher priority.
type SortGroup struct {
	Sorts []Sort
}

// Format
//
//	Formats the sort group into a slice of strings compatible
//	with the meilisearch search engine API
func (sg *SortGroup) Format() []string {
	// create slice to hold output
	var result []string

	// iterate sorts formatting them into their string format and appending to the result slice
	for _, sort := range sg.Sorts {
		result = append(result, sort.String())
	}
	return result
}
