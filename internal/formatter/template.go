package formatter

import (
	"fmt"
	"strings"
	"time"

	"wppanalyticscli/internal/datetime"
	"wppanalyticscli/internal/models"
)

// TemplateFormatter implements the OutputFormatter interface for template analytics
type TemplateFormatter struct{}

// NewTemplateFormatter creates a new template formatter
func NewTemplateFormatter() *TemplateFormatter {
	return &TemplateFormatter{}
}

// FormatTemplate formats the template analytics response as a table
func (f *TemplateFormatter) FormatTemplate(response *models.TemplateAnalyticsResponse, loc *time.Location) string {
	var output strings.Builder
	
	if len(response.Data) == 0 {
		output.WriteString("❌ No template analytics data found.\n")
		return output.String()
	}
	
	data := response.Data[0] // Usually contains one data object
	
	output.WriteString(fmt.Sprintf("📊 Template Analytics Report\n"))
	output.WriteString(fmt.Sprintf("📈 Granularity: %s\n", strings.ToUpper(data.Granularity)))
	output.WriteString(fmt.Sprintf("🔧 Product Type: %s\n", strings.ToUpper(data.ProductType)))
	output.WriteString(fmt.Sprintf("📋 Data Points: %d\n", len(data.DataPoints)))
	output.WriteString(fmt.Sprintf("🌎 Timezone: %s\n\n", loc.String()))
	
	if len(data.DataPoints) == 0 {
		output.WriteString("❌ No data points found.\n")
		return output.String()
	}
	
	// Table header
	output.WriteString("╭──────────────┬─────────────────┬──────────┬───────────┬──────────┬──────────┬───────────┬──────────────╮\n")
	output.WriteString("│     Date     │  Template ID    │   Sent   │ Delivered │   Read   │ Clicked  │   Cost    │ Click Rate % │\n")
	output.WriteString("├──────────────┼─────────────────┼──────────┼───────────┼──────────┼──────────┼───────────┼──────────────┤\n")
	
	totalSent := 0
	totalDelivered := 0
	totalRead := 0
	totalClicked := 0
	totalCost := 0.0
	
	for _, dp := range data.DataPoints {
		date := formatTemplateDate(dp.Start, loc)
		templateID := truncateString(dp.TemplateID, 15)
		
		// Calculate total clicks
		clicks := 0
		for _, clicked := range dp.Clicked {
			clicks += clicked.Count
		}
		
		// Calculate total cost (amount_spent)
		cost := 0.0
		for _, costMetric := range dp.Cost {
			if costMetric.Type == "amount_spent" {
				cost = costMetric.Value
				break
			}
		}
		
		// Calculate click rate
		clickRate := float64(0)
		if dp.Delivered > 0 {
			clickRate = (float64(clicks) / float64(dp.Delivered)) * 100
		}
		
		output.WriteString(fmt.Sprintf("│ %-12s │ %-15s │ %8s │ %9s │ %8s │ %8s │ %9s │ %11.1f%% │\n",
			date, templateID,
			formatNumber(dp.Sent),
			formatNumber(dp.Delivered),
			formatNumber(dp.Read),
			formatNumber(clicks),
			fmt.Sprintf("$%.2f", cost),
			clickRate))
		
		totalSent += dp.Sent
		totalDelivered += dp.Delivered
		totalRead += dp.Read
		totalClicked += clicks
		totalCost += cost
	}
	
	output.WriteString("╰──────────────┴─────────────────┴──────────┴───────────┴──────────┴──────────┴───────────┴──────────────╯\n")
	
	// Summary
	overallClickRate := float64(0)
	if totalDelivered > 0 {
		overallClickRate = (float64(totalClicked) / float64(totalDelivered)) * 100
	}
	
	readRate := float64(0)
	if totalDelivered > 0 {
		readRate = (float64(totalRead) / float64(totalDelivered)) * 100
	}
	
	output.WriteString(fmt.Sprintf("\n📈 Summary:\n"))
	output.WriteString(fmt.Sprintf("   📤 Total Sent: %s\n", formatNumber(totalSent)))
	output.WriteString(fmt.Sprintf("   📥 Total Delivered: %s\n", formatNumber(totalDelivered)))
	output.WriteString(fmt.Sprintf("   👀 Total Read: %s (%.1f%%)\n", formatNumber(totalRead), readRate))
	output.WriteString(fmt.Sprintf("   👆 Total Clicked: %s (%.1f%%)\n", formatNumber(totalClicked), overallClickRate))
	output.WriteString(fmt.Sprintf("   💰 Total Cost: $%.2f\n", totalCost))
	
	if totalCost > 0 && totalDelivered > 0 {
		costPerDelivered := totalCost / float64(totalDelivered)
		output.WriteString(fmt.Sprintf("   📊 Cost per Delivered: $%.4f\n", costPerDelivered))
	}
	
	// Click details if available
	if len(data.DataPoints) > 0 && len(data.DataPoints[0].Clicked) > 0 {
		output.WriteString(fmt.Sprintf("\n🔗 Click Details:\n"))
		clickSummary := make(map[string]int)
		
		for _, dp := range data.DataPoints {
			for _, clicked := range dp.Clicked {
				key := fmt.Sprintf("%s: %s", clicked.Type, clicked.ButtonContent)
				clickSummary[key] += clicked.Count
			}
		}
		
		for action, count := range clickSummary {
			output.WriteString(fmt.Sprintf("   • %s: %d clicks\n", action, count))
		}
	}
	
	return output.String()
}

// formatTemplateDate formats the date for template analytics
func formatTemplateDate(epoch int64, loc *time.Location) string {
	return datetime.ConvertEpochToLocal(epoch, loc).Format("2006-01-02")
}

// truncateString truncates a string to a maximum length
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}