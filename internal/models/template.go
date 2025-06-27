package models

// TemplateAnalyticsResponse represents the Facebook Graph API template analytics response
type TemplateAnalyticsResponse struct {
	Data   []TemplateAnalyticsData `json:"data"`
	Paging *Paging                 `json:"paging,omitempty"`
}

// TemplateAnalyticsData represents template analytics data
type TemplateAnalyticsData struct {
	Granularity  string                    `json:"granularity"`
	ProductType  string                    `json:"product_type"`
	DataPoints   []TemplateDataPoint       `json:"data_points"`
}

// TemplateDataPoint represents a single template analytics data point
type TemplateDataPoint struct {
	TemplateID string              `json:"template_id"`
	Start      int64               `json:"start"`
	End        int64               `json:"end"`
	Sent       int                 `json:"sent"`
	Delivered  int                 `json:"delivered"`
	Read       int                 `json:"read"`
	Clicked    []ClickedAction     `json:"clicked"`
	Cost       []CostMetric        `json:"cost"`
}

// ClickedAction represents a clicked action in template analytics
type ClickedAction struct {
	Type          string `json:"type"`
	ButtonContent string `json:"button_content,omitempty"`
	Count         int    `json:"count"`
}

// CostMetric represents a cost metric in template analytics
type CostMetric struct {
	Type  string  `json:"type"`
	Value float64 `json:"value,omitempty"`
}

// Paging represents pagination information
type Paging struct {
	Cursors *Cursors `json:"cursors,omitempty"`
}

// Cursors represents pagination cursors
type Cursors struct {
	Before string `json:"before,omitempty"`
	After  string `json:"after,omitempty"`
}