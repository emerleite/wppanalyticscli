package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"syscall"
	"time"

	"golang.org/x/term"
)

type AnalyticsResponse struct {
	Analytics struct {
		PhoneNumbers []string    `json:"phone_numbers"`
		Granularity  string      `json:"granularity"`
		DataPoints   []DataPoint `json:"data_points"`
	} `json:"analytics"`
	ID string `json:"id"`
}

type DataPoint struct {
	Start     int64 `json:"start"`
	End       int64 `json:"end"`
	Sent      int   `json:"sent"`
	Delivered int   `json:"delivered"`
}

func main() {
	var wbaID = flag.String("wbaid", "", "WBA ID (required)")
	var startDate = flag.String("start", "", "Start date in ISO-8601 format: YYYY-MM-DD or YYYY-MM-DDTHH:MM:SSZ (required)")
	var endDate = flag.String("end", "", "End date in ISO-8601 format: YYYY-MM-DD or YYYY-MM-DDTHH:MM:SSZ (required)")
	var granularity = flag.String("granularity", "DAY", "Granularity: HALF_HOUR, DAY, or MONTH")
	var timezone = flag.String("timezone", "America/Sao_Paulo", "Timezone for date display (default: America/Sao_Paulo)")
	
	flag.Parse()

	if *wbaID == "" || *startDate == "" || *endDate == "" {
		fmt.Fprintf(os.Stderr, "Usage: %s -wbaid=<id> -start=<date> -end=<date> [-granularity=<HALF_HOUR|DAY|MONTH>] [-timezone=<timezone>]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Date formats: YYYY-MM-DD or YYYY-MM-DDTHH:MM:SSZ\n")
		os.Exit(1)
	}

	loc, err := time.LoadLocation(*timezone)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading timezone '%s': %v\n", *timezone, err)
		os.Exit(1)
	}

	accessToken := os.Getenv("FB_ACCESS_TOKEN")
	if accessToken == "" {
		var err error
		accessToken, err = promptForToken()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading access token: %v\n", err)
			os.Exit(1)
		}
		if accessToken == "" {
			fmt.Fprintf(os.Stderr, "Error: Access token is required\n")
			os.Exit(1)
		}
	}

	startEpoch, err := iso8601ToEpoch(*startDate)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing start date: %v\n", err)
		os.Exit(1)
	}

	endEpoch, err := iso8601ToEpoch(*endDate)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing end date: %v\n", err)
		os.Exit(1)
	}

	if !isValidGranularity(*granularity) {
		fmt.Fprintf(os.Stderr, "Error: granularity must be HALF_HOUR, DAY, or MONTH\n")
		os.Exit(1)
	}

	response, err := makeRequest(*wbaID, startEpoch, endEpoch, *granularity, accessToken, loc)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error making request: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(response)
}

func iso8601ToEpoch(dateStr string) (int64, error) {
	// Try parsing full ISO-8601 datetime first
	t, err := time.Parse(time.RFC3339, dateStr)
	if err == nil {
		return t.Unix(), nil
	}
	
	// Try parsing date-only format (YYYY-MM-DD)
	t, err = time.Parse("2006-01-02", dateStr)
	if err == nil {
		return t.Unix(), nil
	}
	
	// Try parsing date with timezone (YYYY-MM-DD+TZ)
	t, err = time.Parse("2006-01-02Z07:00", dateStr)
	if err == nil {
		return t.Unix(), nil
	}
	
	return 0, fmt.Errorf("invalid date format: expected ISO-8601 datetime (2006-01-02T15:04:05Z) or date (2006-01-02)")
}

func isValidGranularity(g string) bool {
	switch g {
	case "HALF_HOUR", "DAY", "MONTH":
		return true
	default:
		return false
	}
}

func epochToLocalTime(epoch int64, loc *time.Location) time.Time {
	return time.Unix(epoch, 0).In(loc)
}

func formatTimeRange(start, end int64, loc *time.Location, granularity string) (string, string) {
	startTime := epochToLocalTime(start, loc)
	endTime := epochToLocalTime(end, loc)
	
	switch granularity {
	case "DAY":
		date := startTime.Format("2006-01-02")
		timeRange := fmt.Sprintf("%s - %s", startTime.Format("15:04"), endTime.Format("15:04"))
		return date, timeRange
	case "MONTH":
		date := startTime.Format("2006-01")
		timeRange := fmt.Sprintf("%s - %s", startTime.Format("Jan 02"), endTime.Format("Jan 02"))
		return date, timeRange
	case "HALF_HOUR":
		date := startTime.Format("2006-01-02")
		timeRange := fmt.Sprintf("%s - %s", startTime.Format("15:04"), endTime.Format("15:04"))
		return date, timeRange
	default:
		date := startTime.Format("2006-01-02")
		timeRange := fmt.Sprintf("%s - %s", startTime.Format("15:04"), endTime.Format("15:04"))
		return date, timeRange
	}
}

func formatAnalyticsOutput(response *AnalyticsResponse, loc *time.Location) string {
	var output strings.Builder
	
	output.WriteString(fmt.Sprintf("📱 WhatsApp Business Account: %s\n", response.ID))
	output.WriteString(fmt.Sprintf("📞 Phone Numbers: %s\n", strings.Join(response.Analytics.PhoneNumbers, ", ")))
	output.WriteString(fmt.Sprintf("⏱️  Granularity: %s\n", response.Analytics.Granularity))
	output.WriteString(fmt.Sprintf("📊 Data Points: %d\n", len(response.Analytics.DataPoints)))
	output.WriteString(fmt.Sprintf("🌎 Timezone: %s\n\n", loc.String()))
	
	if len(response.Analytics.DataPoints) == 0 {
		output.WriteString("❌ No data points found.\n")
		return output.String()
	}
	
	output.WriteString("╭──────────────┬─────────────────┬─────────────┬─────────────╮\n")
	output.WriteString("│     Date     │   Time Range    │    Sent     │  Delivered  │\n")
	output.WriteString("├──────────────┼─────────────────┼─────────────┼─────────────┤\n")
	
	totalSent := 0
	totalDelivered := 0
	
	for _, dp := range response.Analytics.DataPoints {
		date, timeRange := formatTimeRange(dp.Start, dp.End, loc, response.Analytics.Granularity)
		
		output.WriteString(fmt.Sprintf("│ %-12s │ %-15s │ %11s │ %11s │\n",
			date, timeRange, 
			formatNumber(dp.Sent), 
			formatNumber(dp.Delivered)))
		
		totalSent += dp.Sent
		totalDelivered += dp.Delivered
	}
	
	output.WriteString("╰──────────────┴─────────────────┴─────────────┴─────────────╯\n")
	
	output.WriteString(fmt.Sprintf("\n📈 Summary:\n"))
	output.WriteString(fmt.Sprintf("   📤 Total Sent: %s\n", formatNumber(totalSent)))
	output.WriteString(fmt.Sprintf("   📥 Total Delivered: %s\n", formatNumber(totalDelivered)))
	output.WriteString(fmt.Sprintf("   ℹ️  Note: Delivered messages may arrive after the reporting period\n"))
	
	return output.String()
}

func formatNumber(n int) string {
	if n >= 1000000 {
		return fmt.Sprintf("%.1fM", float64(n)/1000000)
	} else if n >= 1000 {
		return fmt.Sprintf("%.1fK", float64(n)/1000)
	}
	return fmt.Sprintf("%d", n)
}

func makeRequest(wbaID string, start, end int64, granularity, accessToken string, loc *time.Location) (string, error) {
	baseURL := fmt.Sprintf("https://graph.facebook.com/v23.0/%s", wbaID)
	
	params := url.Values{}
	params.Add("fields", fmt.Sprintf("analytics.start(%d).end(%d).granularity(%s)", start, end, granularity))
	params.Add("access_token", accessToken)
	
	fullURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())
	
	resp, err := http.Get(fullURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}
	
	var response AnalyticsResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return string(body), nil
	}
	
	return formatAnalyticsOutput(&response, loc), nil
}

func promptForToken() (string, error) {
	fmt.Fprint(os.Stderr, "Enter Facebook Access Token: ")
	
	// Try to read from terminal with hidden input
	if term.IsTerminal(int(syscall.Stdin)) {
		token, err := term.ReadPassword(int(syscall.Stdin))
		fmt.Fprintln(os.Stderr) // Add newline after hidden input
		if err != nil {
			return "", err
		}
		return strings.TrimSpace(string(token)), nil
	}
	
	// Fallback to regular input if not a terminal
	reader := bufio.NewReader(os.Stdin)
	token, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	
	return strings.TrimSpace(token), nil
}