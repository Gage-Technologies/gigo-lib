package search

import (
	"fmt"
)

// Filter
//
//		 Structures the filtering logic and API format in a programmatic
//		 interface so that filtering can be performed by developers
//		 unfamiliar with the meilisearch API.
//
//				 Usage Notes:
//	              - Attribute: attribute of the documents to filter on
//	              - Operator:  comparison operator used to compare the Attribute and Value
//	              - Value:     value that will be compared to the Attribute via the Operator (not used for OperatorIn and OperatorNotIn)
//	              - Values:    values that will be compared to the Attribute when using the Operators OperatorIn and OperatorNotIn (ignored otherwise)
type Filter struct {
	Attribute string
	Operator  Operator
	Value     interface{}
	Values    []interface{}
}

// String
//
//	Formats the Filter into a string filter
//	supported by the meilisearch search engine
func (f *Filter) String() string {
	// handle exists based operators for their special valueless syntax
	if f.Operator == OperatorExists || f.Operator == OperatorDoesNotExist {
		return fmt.Sprintf("%s %s", f.Attribute, f.Operator)
	}

	// create local variable to hold value that will be formatted into the output
	val := ""

	// conditionally format values slice into an array string if we are performing an in operation
	if f.Operator == OperatorIn || f.Operator == OperatorNotIn {
		// iterate values for `in` based operator formatting them into an array inside a string
		val = "["
		for i, v := range f.Values {
			if i > 0 {
				val += ", "
			}

			// conditionally wrap string types in single quotes
			val += formatInterfaceToString(v, true)
		}
		val += "]"
	} else {
		// conditionally wrap string types in single quotes
		val += formatInterfaceToString(f.Value, true)
	}

	// format normal filters using default `ATTRIBUTE OPERATOR VALUE` format
	return fmt.Sprintf("%s %s %s", f.Attribute, f.Operator, val)
}
