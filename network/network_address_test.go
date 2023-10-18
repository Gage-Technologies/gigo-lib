package network

import (
	"net/http"
	"testing"
)

func TestGetRequestIP(t *testing.T) {
	// create test cases
	testCases := []struct {
		name     string
		request  *http.Request
		expected string
	}{
		{
			name:     "valid ip",
			request:  &http.Request{RemoteAddr: "192.168.1.1"},
			expected: "192.168.1.1",
		},
		{
			name:     "forwarded for",
			request:  &http.Request{RemoteAddr: "192.168.1.1", Header: http.Header{"X-Forwarded-For": []string{"192.168.1.2"}}},
			expected: "192.168.1.2",
		},
		{
			name:     "original forwarded for",
			request:  &http.Request{RemoteAddr: "192.168.1.1", Header: http.Header{"X-Forwarded-For": []string{"192.168.1.2"}, "X-Original-Forwarded-For": []string{"192.168.1.3"}}},
			expected: "192.168.1.3",
		},
		{
			name:     "original forwarded for public",
			request:  &http.Request{RemoteAddr: "192.168.1.1", Header: http.Header{"X-Forwarded-For": []string{"192.168.1.2"}, "X-Original-Forwarded-For": []string{"192.168.1.3", "172.64.32.1"}}},
			expected: "172.64.32.1",
		},
	}

	// run test cases
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			// retrieve ip address
			ip := GetRequestIP(testCase.request)

			// compare ip address
			if ip != testCase.expected {
				t.Errorf("expected %s, got %s", testCase.expected, ip)
			}
		})
	}
}
