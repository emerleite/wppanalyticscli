package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestIso8601ToEpoch(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int64
		hasError bool
	}{
		{
			name:     "Full ISO-8601 datetime",
			input:    "2025-06-24T15:30:00Z",
			expected: 1750779000, // Corrected epoch value
			hasError: false,
		},
		{
			name:     "Date only",
			input:    "2025-06-24",
			expected: 1750723200, // Corrected epoch value
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
			result, err := iso8601ToEpoch(tt.input)
			
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

func TestIsValidGranularity(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"DAY", true},
		{"HALF_HOUR", true},
		{"MONTH", true},
		{"INVALID", false},
		{"day", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := isValidGranularity(tt.input)
			if result != tt.expected {
				t.Errorf("For input '%s', expected %v, got %v", tt.input, tt.expected, result)
			}
		})
	}
}

func TestFormatNumber(t *testing.T) {
	tests := []struct {
		input    int
		expected string
	}{
		{500, "500"},
		{1500, "1.5K"},
		{999, "999"},
		{1000, "1.0K"},
		{1500000, "1.5M"},
		{999999, "1000.0K"},
		{1000000, "1.0M"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			result := formatNumber(tt.input)
			if result != tt.expected {
				t.Errorf("For input %d, expected '%s', got '%s'", tt.input, tt.expected, result)
			}
		})
	}
}

func TestFormatTimeRange(t *testing.T) {
	// Create a test timezone (UTC for simplicity)
	loc, _ := time.LoadLocation("UTC")
	
	// Test epoch timestamps
	start := int64(1750779000) // 2025-06-24T15:30:00Z
	end := int64(1750865400)   // 2025-06-25T15:30:00Z

	tests := []struct {
		granularity   string
		expectedDate  string
		expectedRange string
	}{
		{
			granularity:   "DAY",
			expectedDate:  "2025-06-24",
			expectedRange: "15:30 - 15:30",
		},
		{
			granularity:   "HALF_HOUR",
			expectedDate:  "2025-06-24",
			expectedRange: "15:30 - 15:30",
		},
	}

	for _, tt := range tests {
		t.Run(tt.granularity, func(t *testing.T) {
			date, timeRange := formatTimeRange(start, end, loc, tt.granularity)
			
			if date != tt.expectedDate {
				t.Errorf("Expected date '%s', got '%s'", tt.expectedDate, date)
			}
			
			if timeRange != tt.expectedRange {
				t.Errorf("Expected time range '%s', got '%s'", tt.expectedRange, timeRange)
			}
		})
	}
}

// Mock HTTP server for testing API requests
func createMockServer(response string, statusCode int) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(statusCode)
		w.Write([]byte(response))
	}))
}

func TestMakeRequestWithMockResponse(t *testing.T) {
	// Test with a mock JSON response directly
	mockResponseJSON := `{
		"analytics": {
			"phone_numbers": ["551148619349"],
			"granularity": "DAY",
			"data_points": [
				{
					"start": 1750474800,
					"end": 1750561200,
					"sent": 523,
					"delivered": 539
				},
				{
					"start": 1750561200,
					"end": 1750647600,
					"sent": 92,
					"delivered": 100
				}
			]
		},
		"id": "932157148829117"
	}`

	// Parse the JSON response
	var response AnalyticsResponse
	err := json.Unmarshal([]byte(mockResponseJSON), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal mock response: %v", err)
	}

	// Test the formatting function directly
	loc, _ := time.LoadLocation("UTC")
	result := formatAnalyticsOutput(&response, loc)

	// Check that the result contains expected content
	if !strings.Contains(result, "WhatsApp Business Account: 932157148829117") {
		t.Errorf("Result doesn't contain expected WBA ID")
	}
	
	if !strings.Contains(result, "Phone Numbers: 551148619349") {
		t.Errorf("Result doesn't contain expected phone number")
	}
	
	if !strings.Contains(result, "Total Sent: 615") {
		t.Errorf("Result doesn't contain expected total sent count")
	}
	
	if !strings.Contains(result, "Total Delivered: 639") {
		t.Errorf("Result doesn't contain expected total delivered count")
	}
}

func TestHTTPMockServer(t *testing.T) {
	// Test that our mock server works correctly
	mockResponse := `{
		"analytics": {
			"phone_numbers": ["551148619349"],
			"granularity": "DAY",
			"data_points": []
		},
		"id": "932157148829117"
	}`

	server := createMockServer(mockResponse, http.StatusOK)
	defer server.Close()

	// Make a request to the mock server
	resp, err := http.Get(server.URL)
	if err != nil {
		t.Fatalf("Failed to make request to mock server: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	// Test error response
	errorServer := createMockServer(`{"error": "Invalid token"}`, http.StatusUnauthorized)
	defer errorServer.Close()

	errorResp, err := http.Get(errorServer.URL)
	if err != nil {
		t.Fatalf("Failed to make request to error server: %v", err)
	}
	defer errorResp.Body.Close()

	if errorResp.StatusCode != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", errorResp.StatusCode)
	}
}

func TestFormatAnalyticsOutput(t *testing.T) {
	// Create test data
	response := &AnalyticsResponse{
		ID: "932157148829117",
		Analytics: struct {
			PhoneNumbers []string    `json:"phone_numbers"`
			Granularity  string      `json:"granularity"`
			DataPoints   []DataPoint `json:"data_points"`
		}{
			PhoneNumbers: []string{"551148619349"},
			Granularity:  "DAY",
			DataPoints: []DataPoint{
				{
					Start:     1750474800,
					End:       1750561200,
					Sent:      523,
					Delivered: 539,
				},
				{
					Start:     1750561200,
					End:       1750647600,
					Sent:      92,
					Delivered: 100,
				},
			},
		},
	}

	loc, _ := time.LoadLocation("UTC")
	result := formatAnalyticsOutput(response, loc)

	// Check that the formatted output contains expected elements
	expectedStrings := []string{
		"📱 WhatsApp Business Account: 932157148829117",
		"📞 Phone Numbers: 551148619349",
		"⏱️  Granularity: DAY",
		"📊 Data Points: 2",
		"📤 Total Sent: 615",
		"📥 Total Delivered: 639",
		"ℹ️  Note: Delivered messages may arrive after the reporting period",
	}

	for _, expected := range expectedStrings {
		if !strings.Contains(result, expected) {
			t.Errorf("Formatted output doesn't contain expected string: %s", expected)
		}
	}

	// Check table structure
	if !strings.Contains(result, "╭──────────────┬─────────────────┬─────────────┬─────────────╮") {
		t.Errorf("Formatted output doesn't contain expected table header")
	}
}

func TestFormatAnalyticsOutputEmptyData(t *testing.T) {
	// Test with empty data points
	response := &AnalyticsResponse{
		ID: "932157148829117",
		Analytics: struct {
			PhoneNumbers []string    `json:"phone_numbers"`
			Granularity  string      `json:"granularity"`
			DataPoints   []DataPoint `json:"data_points"`
		}{
			PhoneNumbers: []string{"551148619349"},
			Granularity:  "DAY",
			DataPoints:   []DataPoint{},
		},
	}

	loc, _ := time.LoadLocation("UTC")
	result := formatAnalyticsOutput(response, loc)

	if !strings.Contains(result, "❌ No data points found.") {
		t.Errorf("Expected 'No data points found' message for empty data")
	}
}