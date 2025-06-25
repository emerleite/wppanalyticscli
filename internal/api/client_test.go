package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"wppanalyticscli/internal/models"
)

func TestFacebookGraphClient_GetAnalytics(t *testing.T) {
	// Mock successful API response
	mockResponse := `{
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

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(mockResponse))
	}))
	defer server.Close()

	// Create client with mock server URL
	client := &FacebookGraphClient{
		httpClient: &http.Client{},
		baseURL:    server.URL,
	}

	response, err := client.GetAnalytics("932157148829117", 1750474800, 1750647600, "DAY", "test-token")
	
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
		return
	}
	
	if response.ID != "932157148829117" {
		t.Errorf("Expected ID '932157148829117', got '%s'", response.ID)
	}
	
	if len(response.Analytics.DataPoints) != 2 {
		t.Errorf("Expected 2 data points, got %d", len(response.Analytics.DataPoints))
	}
	
	if response.Analytics.Granularity != "DAY" {
		t.Errorf("Expected granularity 'DAY', got '%s'", response.Analytics.Granularity)
	}
}

func TestFacebookGraphClient_GetAnalyticsError(t *testing.T) {
	// Mock API error response
	errorResponse := `{"error": {"message": "Invalid access token"}}`
	
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(errorResponse))
	}))
	defer server.Close()

	// Create client with mock server URL
	client := &FacebookGraphClient{
		httpClient: &http.Client{},
		baseURL:    server.URL,
	}

	_, err := client.GetAnalytics("932157148829117", 1750474800, 1750647600, "DAY", "invalid-token")
	
	if err == nil {
		t.Errorf("Expected error for invalid token, but got none")
	}
}

func TestMockResponseParsing(t *testing.T) {
	// Test that our mock JSON response parses correctly
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
				}
			]
		},
		"id": "932157148829117"
	}`

	var response models.AnalyticsResponse
	err := json.Unmarshal([]byte(mockResponseJSON), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal mock response: %v", err)
	}

	if response.ID != "932157148829117" {
		t.Errorf("Expected ID '932157148829117', got '%s'", response.ID)
	}
	
	if len(response.Analytics.PhoneNumbers) != 1 || response.Analytics.PhoneNumbers[0] != "551148619349" {
		t.Errorf("Expected phone number '551148619349', got %v", response.Analytics.PhoneNumbers)
	}
}