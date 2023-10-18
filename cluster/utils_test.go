package cluster

import "testing"

func TestExtractNodeIDFromKey(t *testing.T) {
	tests := []struct {
		name string
		key  string
		id   int64
		base string
	}{
		{
			name: "empty key",
			key:  "",
			id:   -1,
			base: "",
		},
		{
			name: "single digit key",
			key:  "/test/state-data/key/1",
			id:   1,
			base: "/test/state-data/key",
		},
		{
			name: "two digit key",
			key:  "/test/state-data/key/12",
			id:   12,
			base: "/test/state-data/key",
		},
		{
			name: "three digit key",
			key:  "/test/state-data/key/123",
			id:   123,
			base: "/test/state-data/key",
		},
		{
			name: "bad key",
			key:  "/test/state-data/key/badkey",
			id:   -1,
			base: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id, base, _ := extractNodeIDFromKey(tt.key)
			if id != tt.id {
				t.Errorf("ExtractNodeIDFromKey() = %v, want %v", id, tt.id)
			}
			if base != tt.base {
				t.Errorf("ExtractNodeIDFromKey() = %v, want %v", base, tt.base)
			}
		})
	}
}
