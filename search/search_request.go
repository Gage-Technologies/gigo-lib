package search

// Request
//
//		Struct used to format the string based filter, sort and facet
//		functionality of meilisearch into a programmatic API that does
//		not require understanding of the meilisearch API
//
//	 Usage Notes:
//	     - Query:   query that will be search
//	     - Filter:  filter that will be applied to the query results
//	     - Sort:    sort that will be applied to the query+filter results
//	     - Facet:   facets are the fields that you would like a distribution breakdown for the results for
//	                the search (see https://docs.meilisearch.com/learn/advanced/filtering_and_faceted_search.html#filtering-with-georadius)
//	     - Offset:  amount of documents to skip before returning any results
//	     - Limit:   max amount of results that can be returned after the offset
type Request struct {
	Query  string
	Filter *FilterGroup
	Sort   *SortGroup
	Facets []string
	Offset int
	Limit  int
}
