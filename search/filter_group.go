package search

// FilterGroup
//
//			Merges multiple filters together allowing for more complex filtering
//			operations. The internal structure is a 2D slice of Filters that
//			can be used to perform the logical merge of multiple filter sets
//
//		 Usage Notes:
//		   - Conditions are interpreted in sequential order of index as if written left (index 0) to right (index > 0)
//		   - And: converts the logical join of all filters conditions in the group from a logical OR to a logical AND
//		   - Not: inverts the logical meaning of the total filter group
//
//		 Ex 1:
//		     field1 > 10
//		     	FilterGroup{
//				        Filters: []FilterCondition{
//					        {
//						        Filters: []Filter{
//							        {
//								        Value:     "10",
//								        Attribute: "field1",
//								        Operator:  OperatorGreaterThan,
//							        },
//						        },
//					        },
//		             }
//		         }
//
//			Ex 2:
//		     NOT ((field1 > 10 && field2 == 'testGreater') || NOT (field1 < 10 && field2 == 'testLess'))
//
//		     	FilterGroup{
//				        Filters: []FilterCondition{
//					        {
//						        Filters: []Filter{
//							        {
//								        Value:     "10",
//								        Attribute: "field1",
//								        Operator:  OperatorGreaterThan,
//							        },
//							        {
//								        Value:     "testGreater",
//								        Attribute: "field2",
//								        Operator:  OperatorEquals,
//							        },
//						        },
//						        And: true,
//					        },
//					        {
//						        Filters: []Filter{
//							        {
//								        Value:     "10",
//								        Attribute: "field1",
//								        Operator:  OperatorLessThan,
//							        },
//							        {
//								        Value:     "testLess",
//								        Attribute: "field2",
//								        Operator:  OperatorEquals,
//							        },
//						        },
//						        And: true,
//		                     Not: true,
//					        },
//				        },
//				        And: false,
//	                 Not: true,
//			        }
type FilterGroup struct {
	Filters []FilterCondition
	And     bool
	Not     bool
}

// String
//
//	Formats the FilterGroup into a string filter
//	supported by the meilisearch search engine
func (fg *FilterGroup) String() string {
	// create string to append to
	out := ""

	// add logical NOT to invert the value of filter conditions
	if fg.Not {
		out += "NOT "
		// add parenthesis to group filter conditions if there is more than one filter condition
		if len(fg.Filters) > 1 {
			out += "("
		}
	}

	// iterate through the filter conditions formatting them to strings and adding them to the output string
	for i, filter := range fg.Filters {
		// format filter condition
		fc := filter.String()

		// conditionally add parenthesis to the filter condition if it does not contain a logical NOT
		// since the FilterCondition.String() function does not wrap the condition in parentheses
		// unless the condition is a logical NOT
		if !filter.Not && len(filter.Filters) > 1 {
			fc = "(" + fc + ")"
		}

		// append the filter condition to the output string
		out += fc

		// conditionally add the condition's logical join
		if i < len(fg.Filters)-1 {
			if fg.And {
				out += " AND "
			} else {
				out += " OR "
			}
		}
	}

	// close parentheses for group filters if there is more than one filter condition and this is a logical NOT filter
	if fg.Not && len(fg.Filters) > 1 {
		out += ")"
	}

	return out
}
