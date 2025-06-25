package formatter

import (
	"strings"
	"testing"
	"time"

	"wppanalyticscli/internal/models"
)

func TestTableFormatter_Format(t *testing.T) {
	formatter := NewTableFormatter()
	
	// Create test data
	response := &models.AnalyticsResponse{
		ID: "932157148829117",
		Analytics: struct {
			PhoneNumbers []string              `json:"phone_numbers"`
			Granularity  string                `json:"granularity"`
			DataPoints   []models.DataPoint    `json:"data_points"`
		}{
			PhoneNumbers: []string{"551148619349"},
			Granularity:  "DAY",
			DataPoints: []models.DataPoint{
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
	result := formatter.Format(response, loc)

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

func TestTableFormatter_FormatEmptyData(t *testing.T) {
	formatter := NewTableFormatter()
	
	// Test with empty data points
	response := &models.AnalyticsResponse{
		ID: "932157148829117",
		Analytics: struct {
			PhoneNumbers []string              `json:"phone_numbers"`
			Granularity  string                `json:"granularity"`
			DataPoints   []models.DataPoint    `json:"data_points"`
		}{
			PhoneNumbers: []string{"551148619349"},
			Granularity:  "DAY",
			DataPoints:   []models.DataPoint{},
		},
	}

	loc, _ := time.LoadLocation("UTC")
	result := formatter.Format(response, loc)

	if !strings.Contains(result, "❌ No data points found.") {
		t.Errorf("Expected 'No data points found' message for empty data")
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