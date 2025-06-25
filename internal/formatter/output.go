package formatter

import (
	"fmt"
	"strings"
	"time"

	"wppanalyticscli/internal/datetime"
	"wppanalyticscli/internal/models"
)

// OutputFormatter defines the interface for formatting output
type OutputFormatter interface {
	Format(response *models.AnalyticsResponse, loc *time.Location) string
}

// TableFormatter implements the OutputFormatter interface for table output
type TableFormatter struct{}

// NewTableFormatter creates a new table formatter
func NewTableFormatter() *TableFormatter {
	return &TableFormatter{}
}

// Format formats the analytics response as a table
func (f *TableFormatter) Format(response *models.AnalyticsResponse, loc *time.Location) string {
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

// formatTimeRange formats the time range based on granularity
func formatTimeRange(start, end int64, loc *time.Location, granularity string) (string, string) {
	startTime := datetime.ConvertEpochToLocal(start, loc)
	endTime := datetime.ConvertEpochToLocal(end, loc)
	
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

// formatNumber formats numbers with K/M suffixes
func formatNumber(n int) string {
	if n >= 1000000 {
		return fmt.Sprintf("%.1fM", float64(n)/1000000)
	} else if n >= 1000 {
		return fmt.Sprintf("%.1fK", float64(n)/1000)
	}
	return fmt.Sprintf("%d", n)
}