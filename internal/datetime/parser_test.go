package datetime

import (
	"testing"
)

func TestISO8601Parser_ParseToEpoch(t *testing.T) {
	parser := NewISO8601Parser()
	
	tests := []struct {
		name     string
		input    string
		expected int64
		hasError bool
	}{
		{
			name:     "Full ISO-8601 datetime",
			input:    "2025-06-24T15:30:00Z",
			expected: 1750779000,
			hasError: false,
		},
		{
			name:     "Date only",
			input:    "2025-06-24",
			expected: 1750723200,
			hasError: false,
		},
		{
			name:     "Invalid format",
			input:    "invalid-date",
			expected: 0,
			hasError: true,
		},
		{
			name:     "Empty string",
			input:    "",
			expected: 0,
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parser.ParseToEpoch(tt.input)
			
			if tt.hasError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}
			
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}
			
			if result != tt.expected {
				t.Errorf("Expected %d, got %d", tt.expected, result)
			}
		})
	}
}