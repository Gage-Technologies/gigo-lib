package search

import (
	"github.com/gage-technologies/gigo-lib/db/models"
	"testing"
	"time"
)

func TestFilter_String(t *testing.T) {
	tests := []struct {
		name   string
		filter Filter
		want   string
	}{
		{
			name: "simple1",
			filter: Filter{
				Attribute: "field1",
				Operator:  OperatorEquals,
				Value:     "value1",
			},
			want: "field1 = 'value1'",
		},
		{
			name: "simple2",
			filter: Filter{
				Attribute: "field1",
				Operator:  OperatorGreaterThanOrEquals,
				Value:     "value1",
			},
			want: "field1 >= 'value1'",
		},
		{
			name: "notIn1",
			filter: Filter{
				Attribute: "field1",
				Operator:  OperatorNotIn,
				// value should be ignored
				Value: "value1",
				// Values should be used in place of Value
				Values: []interface{}{"value1", "value2", 1, true},
			},
			want: "field1 NOT IN ['value1', 'value2', 1, true]",
		},
		{
			name: "exists1",
			filter: Filter{
				Attribute: "field1",
				Operator:  OperatorExists,
				// Value and Values should be ignored
				Value:  "value1",
				Values: []interface{}{"value1", "value2"},
			},
			want: "field1 EXISTS",
		},
		{
			name: "customType1",
			filter: Filter{
				Attribute: "lang",
				Operator:  OperatorEquals,
				Value:     models.Go,
			},
			want: "lang = 6",
		},
		{
			name: "time1",
			filter: Filter{
				Attribute: "since",
				Operator:  OperatorGreaterThanOrEquals,
				Value:     time.Unix(300, 0),
			},
			want: "since >= 300",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.filter.String(); got != tt.want {
				t.Fatalf("\n%s failed\n    Error: got %q wanted %q", t.Name(), got, tt.want)
			}
			t.Logf("\n%s succeeded", t.Name())
		})
	}
}
