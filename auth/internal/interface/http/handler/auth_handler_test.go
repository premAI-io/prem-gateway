package httphandler

import "testing"

func TestExtractService(t *testing.T) {
	tests := []struct {
		host     string
		uri      string
		expected string
	}{
		{
			host:     "service.prem.com",
			uri:      "https://service.prem.com/notrelevant/path",
			expected: "service",
		},
		{
			host:     "1.1.1.1",
			uri:      "http://1.1.1.1/service/v1/chat",
			expected: "service",
		},
		{
			host:     "service.sub.prem.com",
			uri:      "https://service.sub.prem.com/notrelevant/path",
			expected: "service",
		},
	}

	for _, tt := range tests {
		result := extractService(tt.host, tt.uri)
		if result != tt.expected {
			t.Errorf("For host=%s and uri=%s, expected %s but got %s", tt.host, tt.uri, tt.expected, result)
		}
	}
}
