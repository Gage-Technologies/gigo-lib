package search

import "testing"

func TestFilterCondition_String(t *testing.T) {
	tests := []struct {
		name string
		f    FilterCondition
		want string
	}{
		{
			name: "simple1",
			f: FilterCondition{
				Filters: []Filter{
					{
						Attribute: "field1",
						Operator:  OperatorEquals,
						Value:     "value1",
					},
				},
			},
			want: "field1 = 'value1'",
		},
		{
			name: "simple2",
			f: FilterCondition{
				Filters: []Filter{
					{
						Attribute: "field1",
						Operator:  OperatorEquals,
						Value:     "value1",
					},
					{
						Attribute: "field2",
						Operator:  OperatorExists,
					},
				},
			},
			want: "field1 = 'value1' OR field2 EXISTS",
		},
		{
			name: "complex1",
			f: FilterCondition{
				Filters: []Filter{
					{
						Attribute: "field1",
						Operator:  OperatorEquals,
						Value:     "value1",
					},
					{
						Attribute: "field2",
						Operator:  OperatorExists,
					},
					{
						Attribute: "field3",
						Operator:  OperatorIn,
						// value should be ignored
						Value: "value3",
						// values should be used in place of value
						Values: []interface{}{"value4", "value5", 1, true},
					},
				},
				And: true,
				Not: true,
			},
			want: "NOT (field1 = 'value1' AND field2 EXISTS AND field3 IN ['value4', 'value5', 1, true])",
		},
		{
			name: "complex1",
			f: FilterCondition{
				Filters: []Filter{
					{
						Attribute: "field1",
						Operator:  OperatorEquals,
						Value:     "value1",
					},
					{
						Attribute: "field2",
						Operator:  OperatorExists,
					},
					{
						Attribute: "field3",
						Operator:  OperatorIn,
						// value should be ignored
						Value: "value3",
						// values should be used in place of value
						Values: []interface{}{"value4", "value5"},
					},
					{
						Attribute: "field4",
						Operator:  OperatorNotIn,
						// values should be used in place of value
						Values: []interface{}{"value6", "value7", 1, 10, true},
					},
				},
				Not: true,
			},
			want: "NOT (field1 = 'value1' OR field2 EXISTS OR field3 IN ['value4', 'value5'] OR field4 NOT IN ['value6', 'value7', 1, 10, true])",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.f.String(); got != tt.want {
				t.Fatalf("\n%s failed\n    Error: got %q want %q", tt.name, got, tt.want)
			}
			t.Logf("\n%s succeeded\n", t.Name())
		})
	}
}
