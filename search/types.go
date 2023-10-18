package search

// Operator
//
//	Performs a comparison between two values
type Operator string

const (
	OperatorEquals              Operator = "="
	OperatorNotEquals           Operator = "!="
	OperatorGreaterThan         Operator = ">"
	OperatorLessThan            Operator = "<"
	OperatorGreaterThanOrEquals Operator = ">="
	OperatorLessThanOrEquals    Operator = "<="
	OperatorIn                  Operator = "IN"
	OperatorNotIn               Operator = "NOT IN"
	OperatorExists              Operator = "EXISTS"
	OperatorDoesNotExist        Operator = "NOT EXISTS"
)
