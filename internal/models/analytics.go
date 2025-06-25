package models

// AnalyticsResponse represents the Facebook Graph API analytics response
type AnalyticsResponse struct {
	Analytics struct {
		PhoneNumbers []string    `json:"phone_numbers"`
		Granularity  string      `json:"granularity"`
		DataPoints   []DataPoint `json:"data_points"`
	} `json:"analytics"`
	ID string `json:"id"`
}

// DataPoint represents a single analytics data point
type DataPoint struct {
	Start     int64 `json:"start"`
	End       int64 `json:"end"`
	Sent      int   `json:"sent"`
	Delivered int   `json:"delivered"`
}