package formatter

import (
	"strings"
	"testing"
	"time"

	"wppanalyticscli/internal/models"
)

func TestTemplateFormatter_FormatTemplate(t *testing.T) {
	formatter := NewTemplateFormatter()
	
	// Create test data
	response := &models.TemplateAnalyticsResponse{
		Data: []models.TemplateAnalyticsData{
			{
				Granularity: "DAILY",
				ProductType: "cloud_api",
				DataPoints: []models.TemplateDataPoint{
					{
						TemplateID: "1026573095658757",
						Start:      1750377600,
						End:        1750464000,
						Sent:       871,
						Delivered:  789,
						Read:       399,
						Clicked: []models.ClickedAction{
							{
								Type:          "quick_reply_button",
								ButtonContent: "Quero negociar",
								Count:         56,
							},
						},
						Cost: []models.CostMetric{
							{
								Type:  "amount_spent",
								Value: 6.18,
							},
						},
					},
					{
						TemplateID: "1026573095658757",
						Start:      1750464000,
						End:        1750550400,
						Sent:       0,
						Delivered:  6,
						Read:       36,
						Clicked: []models.ClickedAction{
							{
								Type:          "quick_reply_button",
								ButtonContent: "Quero negociar",
								Count:         4,
							},
						},
						Cost: []models.CostMetric{
							{
								Type:  "amount_spent",
								Value: 0.04,
							},
						},
					},
				},
			},
		},
	}

	loc, _ := time.LoadLocation("UTC")
	result := formatter.FormatTemplate(response, loc)

	// Check that the formatted output contains expected elements
	expectedStrings := []string{
		"📊 Template Analytics Report",
		"📈 Granularity: DAILY",
		"🔧 Product Type: CLOUD_API",
		"📋 Data Points: 2",
		"📤 Total Sent: 871",
		"📥 Total Delivered: 795",
		"👀 Total Read: 435",
		"👆 Total Clicked: 60",
		"💰 Total Cost: $6.22",
		"🔗 Click Details:",
		"quick_reply_button: Quero negociar: 60 clicks",
	}

	for _, expected := range expectedStrings {
		if !strings.Contains(result, expected) {
			t.Errorf("Formatted output doesn't contain expected string: %s", expected)
		}
	}

	// Check table structure
	if !strings.Contains(result, "╭──────────────┬─────────────────┬──────────┬───────────┬──────────┬──────────┬───────────┬──────────────╮") {
		t.Errorf("Formatted output doesn't contain expected table header")
	}
}

func TestTemplateFormatter_FormatTemplateEmptyData(t *testing.T) {
	formatter := NewTemplateFormatter()
	
	// Test with empty data
	response := &models.TemplateAnalyticsResponse{
		Data: []models.TemplateAnalyticsData{},
	}

	loc, _ := time.LoadLocation("UTC")
	result := formatter.FormatTemplate(response, loc)

	if !strings.Contains(result, "❌ No template analytics data found.") {
		t.Errorf("Expected 'No template analytics data found' message for empty data")
	}
}

func TestTruncateString(t *testing.T) {
	tests := []struct {
		input    string
		maxLen   int
		expected string
	}{
		{"short", 10, "short"},
		{"this is a very long string", 10, "this is..."},
		{"exactly10c", 10, "exactly10c"},
		{"", 5, ""},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := truncateString(tt.input, tt.maxLen)
			if result != tt.expected {
				t.Errorf("For input '%s' with maxLen %d, expected '%s', got '%s'", tt.input, tt.maxLen, tt.expected, result)
			}
		})
	}
}