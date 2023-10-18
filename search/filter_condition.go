package search

// FilterCondition
//
//	 Groups multiple filters together with a logical AND/OR
//	 FilterConditions default to OR unless the attribute `And` is set true
//
//	 Usage Notes:
//	   - Filters are interpreted in sequential order of index as if written left (index 0) to right (index > 0)
//	   - And: converts the logical join of all filters in the condition from a logical OR to a logical AND
//	   - Not: inverts the logical meaning of the total filter condition
//
//	 Ex 1:
//	     field1 > 10
//	         FilterCondition{
//					Filters: []Filter{
//						{
//							Value:     "10",
//							Attribute: "field1",
//							Operator:  OperatorGreaterThan,
//						},
//					},
//				}
//
//	 Ex 2:
//	     NOT (field1 > 10 && field2 == 'testGreater')
//	         FilterCondition{
//	               Filters: []Filter{
//	                   {
//	                       Value:     "10",
//	                       Attribute: "field1",
//	                       Operator:  OperatorGreaterThan,
//	                   },
//	                   {
//	                       Value:     "testGreater",
//	                       Attribute: "field2",
//	                       Operator:  OperatorEquals,
//	                   },
//	               },
//	               And: true,
//	               Not: true,
//	         }
type FilterCondition struct {
	Filters []Filter
	And     bool
	Not     bool
}

// String
//
//	Formats the FilterCondition into a string filter
//	supported by the meilisearch search engine
func (fc *FilterCondition) String() string {
	// create string to append to
	out := ""

	// add logical NOT to invert the value of filters
	if fc.Not {
		out += "NOT "
		// add parenthesis to group filters if there is more than one filter
		if len(fc.Filters) > 1 {
			out += "("
		}
	}

	// add each filter
	for i, f := range fc.Filters {
		// format filter and append to the output string
		out += f.String()

		// conditionally add the filter's logical join
		if i < len(fc.Filters)-1 {
			if fc.And {
				out += " AND "
			} else {
				out += " OR "
			}
		}
	}

	// close parentheses for group filters if there is more than one filter and this is a logical NOT filter
	if fc.Not && len(fc.Filters) > 1 {
		out += ")"
	}

	return out
}
