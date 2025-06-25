package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"wppanalyticscli/internal/models"
)

// Client defines the interface for API operations
type Client interface {
	GetAnalytics(wbaID string, start, end int64, granularity, accessToken string) (*models.AnalyticsResponse, error)
}

// FacebookGraphClient implements the Client interface for Facebook Graph API
type FacebookGraphClient struct {
	httpClient *http.Client
	baseURL    string
}

// NewFacebookGraphClient creates a new Facebook Graph API client
func NewFacebookGraphClient() *FacebookGraphClient {
	return &FacebookGraphClient{
		httpClient: &http.Client{},
		baseURL:    "https://graph.facebook.com/v23.0",
	}
}

// GetAnalytics fetches analytics data from Facebook Graph API
func (c *FacebookGraphClient) GetAnalytics(wbaID string, start, end int64, granularity, accessToken string) (*models.AnalyticsResponse, error) {
	requestURL := fmt.Sprintf("%s/%s", c.baseURL, wbaID)
	
	params := url.Values{}
	params.Add("fields", fmt.Sprintf("analytics.start(%d).end(%d).granularity(%s)", start, end, granularity))
	params.Add("access_token", accessToken)
	
	fullURL := fmt.Sprintf("%s?%s", requestURL, params.Encode())
	
	resp, err := c.httpClient.Get(fullURL)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}
	
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}
	
	var response models.AnalyticsResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}
	
	return &response, nil
}