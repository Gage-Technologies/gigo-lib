package search

import "testing"

func TestFilterGroup_String(t *testing.T) {
	tests := []struct {
		name string
		f    FilterGroup
		want string
	}{
		{
			name: "simple1",
			f: FilterGroup{
				Filters: []FilterCondition{
					{
						Filters: []Filter{
							{
								Attribute: "field1",
								Operator:  OperatorEquals,
								Value:     "value1",
							},
						},
					},
				},
			},
			want: "field1 = 'value1'",
		},
		{
			name: "simple2",
			f: FilterGroup{
				Filters: []FilterCondition{
					{
						Filters: []Filter{
							{
								Attribute: "field2",
								Operator:  OperatorEquals,
								Value:     "value2",
							},
						},
						Not: true,
					},
				},
			},
			want: "NOT field2 = 'value2'",
		},
		{
			name: "simple3",
			f: FilterGroup{
				Filters: []FilterCondition{
					{
						Filters: []Filter{
							{
								Attribute: "field3",
								Operator:  OperatorEquals,
								Value:     "value3",
							},
						},
					},
				},
				Not: true,
			},
			want: "NOT field3 = 'value3'",
		},
		{
			name: "simple3",
			f: FilterGroup{
				Filters: []FilterCondition{
					{
						Filters: []Filter{
							{
								Attribute: "field4",
								Operator:  OperatorEquals,
								Value:     "value4",
							},
							{
								Attribute: "field5",
								Operator:  OperatorLessThanOrEquals,
								Value:     12,
							},
						},
						And: true,
					},
				},
				Not: true,
			},
			want: "NOT (field4 = 'value4' AND field5 <= 12)",
		},
		{
			name: "simple4",
			f: FilterGroup{
				Filters: []FilterCondition{
					{
						Filters: []Filter{
							{
								Attribute: "field4",
								Operator:  OperatorEquals,
								Value:     "value4",
							},
						},
					},
					{
						Filters: []Filter{
							{
								Attribute: "field5",
								Operator:  OperatorLessThanOrEquals,
								Value:     "value5",
							},
						},
					},
				},
				And: true,
				Not: true,
			},
			want: "NOT (field4 = 'value4' AND field5 <= 'value5')",
		},
		{
			name: "complex1",
			f: FilterGroup{
				Filters: []FilterCondition{
					{
						Filters: []Filter{
							{
								Attribute: "field1",
								Operator:  OperatorLessThan,
								Value:     "value1",
							},
							{
								Attribute: "field2",
								Operator:  OperatorDoesNotExist,
							},
							{
								Attribute: "field3",
								Operator:  OperatorNotIn,
								Values:    []interface{}{1, 2, 3, 4, true, "val", "test"},
							},
							{
								Attribute: "field4",
								Operator:  OperatorEquals,
								Value:     "value4",
							},
							{
								Attribute: "field5",
								Operator:  OperatorExists,
							},
						},
					},
					{
						Filters: []Filter{
							{
								Attribute: "field6",
								Operator:  OperatorGreaterThanOrEquals,
								Value:     "value6",
							},
							{
								Attribute: "field7",
								Operator:  OperatorExists,
							},
							{
								Attribute: "field8",
								Operator:  OperatorIn,
								Values:    []interface{}{5, 6, 7, false, "v", "test2"},
							},
							{
								Attribute: "field9",
								Operator:  OperatorNotEquals,
								Value:     "value9",
							},
							{
								Attribute: "field10",
								Operator:  OperatorDoesNotExist,
							},
						},
						And: true,
						Not: true,
					},
				},
				And: true,
				Not: true,
			},
			want: "NOT ((field1 < 'value1' OR field2 NOT EXISTS OR field3 NOT IN [1, 2, 3, 4, true, 'val', 'test'] OR field4 = 'value4' OR field5 EXISTS) AND NOT (field6 >= 'value6' AND field7 EXISTS AND field8 IN [5, 6, 7, false, 'v', 'test2'] AND field9 != 'value9' AND field10 NOT EXISTS))",
		},
		{
			name: "complex2",
			f: FilterGroup{
				Filters: []FilterCondition{
					{
						Filters: []Filter{
							{
								Attribute: "field1",
								Operator:  OperatorLessThan,
								Value:     "value1",
							},
							{
								Attribute: "field2",
								Operator:  OperatorDoesNotExist,
							},
							{
								Attribute: "field3",
								Operator:  OperatorNotIn,
								Values:    []interface{}{1, 2, 3, 4, true, "val", "test"},
							},
							{
								Attribute: "field4",
								Operator:  OperatorEquals,
								Value:     "value4",
							},
							{
								Attribute: "field5",
								Operator:  OperatorExists,
							},
						},
					},
					{
						Filters: []Filter{
							{
								Attribute: "field6",
								Operator:  OperatorGreaterThanOrEquals,
								Value:     "value6",
							},
							{
								Attribute: "field7",
								Operator:  OperatorExists,
							},
							{
								Attribute: "field8",
								Operator:  OperatorIn,
								Values:    []interface{}{5, 6, 7, false, "v", "test2"},
							},
							{
								Attribute: "field9",
								Operator:  OperatorNotEquals,
								Value:     "value9",
							},
							{
								Attribute: "field10",
								Operator:  OperatorDoesNotExist,
							},
						},
						And: true,
						Not: true,
					},
					{
						Filters: []Filter{
							{
								Attribute: "field1",
								Operator:  OperatorEquals,
								Value:     "value1",
							},
						},
					},
					{
						Filters: []Filter{
							{
								Attribute: "field4",
								Operator:  OperatorEquals,
								Value:     "value4",
							},
							{
								Attribute: "field5",
								Operator:  OperatorLessThanOrEquals,
								Value:     "value5",
							},
						},
						And: true,
					},
				},
			},
			want: "(field1 < 'value1' OR field2 NOT EXISTS OR field3 NOT IN [1, 2, 3, 4, true, 'val', 'test'] OR field4 = 'value4' OR field5 EXISTS) OR NOT (field6 >= 'value6' AND field7 EXISTS AND field8 IN [5, 6, 7, false, 'v', 'test2'] AND field9 != 'value9' AND field10 NOT EXISTS) OR field1 = 'value1' OR (field4 = 'value4' AND field5 <= 'value5')",
		},
		{
			name: "empty",
			f: FilterGroup{
				And: true,
			},
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.f.String(); got != tt.want {
				t.Fatalf("\n%s failed\n    Error: got %q wanted: %q", t.Name(), got, tt.want)
			}
			t.Logf("\n%s succeeded", t.Name())
		})
	}
}
