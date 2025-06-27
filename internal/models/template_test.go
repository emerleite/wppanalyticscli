package models

import (
	"encoding/json"
	"testing"
)

func TestTemplateAnalyticsResponse_Unmarshal(t *testing.T) {
	// Test with real API response format
	mockResponseJSON := `{
		"data": [{
			"granularity": "DAILY",
			"product_type": "cloud_api",
			"data_points": [{
				"template_id": "1026573095658757",
				"start": 1750377600,
				"end": 1750464000,
				"sent": 871,
				"delivered": 789,
				"read": 399,
				"clicked": [{
					"type": "quick_reply_button",
					"button_content": "Quero negociar",
					"count": 56
				}],
				"cost": [{
					"type": "amount_spent",
					"value": 6.18
				}, {
					"type": "cost_per_delivered",
					"value": 0.01
				}]
			}]
		}],
		"paging": {
			"cursors": {
				"before": "MAZDZD",
				"after": "MjQZD"
			}
		}
	}`

	var response TemplateAnalyticsResponse
	err := json.Unmarshal([]byte(mockResponseJSON), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal template response: %v", err)
	}

	// Validate parsed data
	if len(response.Data) != 1 {
		t.Errorf("Expected 1 data object, got %d", len(response.Data))
	}

	data := response.Data[0]
	if data.Granularity != "DAILY" {
		t.Errorf("Expected granularity 'DAILY', got '%s'", data.Granularity)
	}

	if data.ProductType != "cloud_api" {
		t.Errorf("Expected product type 'cloud_api', got '%s'", data.ProductType)
	}

	if len(data.DataPoints) != 1 {
		t.Errorf("Expected 1 data point, got %d", len(data.DataPoints))
	}

	dp := data.DataPoints[0]
	if dp.TemplateID != "1026573095658757" {
		t.Errorf("Expected template ID '1026573095658757', got '%s'", dp.TemplateID)
	}

	if dp.Sent != 871 {
		t.Errorf("Expected sent 871, got %d", dp.Sent)
	}

	if len(dp.Clicked) != 1 {
		t.Errorf("Expected 1 clicked action, got %d", len(dp.Clicked))
	}

	clicked := dp.Clicked[0]
	if clicked.Type != "quick_reply_button" {
		t.Errorf("Expected click type 'quick_reply_button', got '%s'", clicked.Type)
	}

	if clicked.Count != 56 {
		t.Errorf("Expected click count 56, got %d", clicked.Count)
	}

	if len(dp.Cost) != 2 {
		t.Errorf("Expected 2 cost metrics, got %d", len(dp.Cost))
	}

	// Check paging
	if response.Paging == nil {
		t.Error("Expected paging object, got nil")
	} else {
		if response.Paging.Cursors.Before != "MAZDZD" {
			t.Errorf("Expected before cursor 'MAZDZD', got '%s'", response.Paging.Cursors.Before)
		}
	}
}