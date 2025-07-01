package formatter

import (
	"fmt"
	"strings"

	"wppanalyticscli/internal/models"
)

// ListFormatter implements the formatter for template lists
type ListFormatter struct{}

// NewListFormatter creates a new list formatter
func NewListFormatter() *ListFormatter {
	return &ListFormatter{}
}

// FormatList formats the template list response as a table
func (f *ListFormatter) FormatList(response *models.TemplateListResponse) string {
	var output strings.Builder
	
	output.WriteString(fmt.Sprintf("📋 WhatsApp Message Templates\n"))
	output.WriteString(fmt.Sprintf("📊 Total Templates: %d\n\n", len(response.Data)))
	
	if len(response.Data) == 0 {
		output.WriteString("❌ No templates found.\n")
		return output.String()
	}
	
	// Table header
	output.WriteString("╭──────────────────────┬─────────────────────────────────┬──────────┬─────────────┬─────────────╮\n")
	output.WriteString("│         ID           │              Name               │ Language │   Status    │  Category   │\n")
	output.WriteString("├──────────────────────┼─────────────────────────────────┼──────────┼─────────────┼─────────────┤\n")
	
	for _, template := range response.Data {
		templateID := truncateString(template.ID, 20)
		templateName := truncateString(template.Name, 31)
		language := truncateString(template.Language, 8)
		status := formatStatusSimple(template.Status)
		category := truncateString(template.Category, 11)
		
		output.WriteString(fmt.Sprintf("│ %-20s │ %-31s │ %-8s │ %-11s │ %-11s │\n",
			templateID, templateName, language, status, category))
	}
	
	output.WriteString("╰──────────────────────┴─────────────────────────────────┴──────────┴─────────────┴─────────────╯\n")
	
	// Status summary
	statusCounts := make(map[string]int)
	categoryCounts := make(map[string]int)
	languageCounts := make(map[string]int)
	
	for _, template := range response.Data {
		statusCounts[template.Status]++
		categoryCounts[template.Category]++
		languageCounts[template.Language]++
	}
	
	output.WriteString(fmt.Sprintf("\n📈 Summary:\n"))
	
	// Status breakdown
	output.WriteString(fmt.Sprintf("   📊 Status Breakdown:\n"))
	for status, count := range statusCounts {
		emoji := getStatusEmoji(status)
		output.WriteString(fmt.Sprintf("      %s %s: %d\n", emoji, strings.ToUpper(status), count))
	}
	
	// Category breakdown
	if len(categoryCounts) > 0 {
		output.WriteString(fmt.Sprintf("   🏷️  Category Breakdown:\n"))
		for category, count := range categoryCounts {
			output.WriteString(fmt.Sprintf("      • %s: %d\n", strings.ToUpper(category), count))
		}
	}
	
	// Language breakdown
	if len(languageCounts) > 0 {
		output.WriteString(fmt.Sprintf("   🌐 Language Breakdown:\n"))
		for language, count := range languageCounts {
			output.WriteString(fmt.Sprintf("      • %s: %d\n", strings.ToUpper(language), count))
		}
	}
	
	// Pagination info
	if response.Paging != nil && response.Paging.Cursors != nil {
		output.WriteString(fmt.Sprintf("\n📄 Pagination:\n"))
		if response.Paging.Cursors.After != "" {
			// Show only first 20 characters of cursor for readability
			cursorDisplay := response.Paging.Cursors.After
			if len(cursorDisplay) > 20 {
				cursorDisplay = cursorDisplay[:20] + "..."
			}
			output.WriteString(fmt.Sprintf("   Next Page Available: %s\n", cursorDisplay))
			output.WriteString(fmt.Sprintf("   Use: -after=\"%s\" to get next page\n", response.Paging.Cursors.After))
		}
	}
	
	return output.String()
}

// formatStatus formats the status with appropriate styling
func formatStatus(status string) string {
	switch strings.ToLower(status) {
	case "approved":
		return "✅ APPROVED"
	case "pending":
		return "⏳ PENDING"
	case "rejected":
		return "❌ REJECTED"
	case "disabled":
		return "🚫 DISABLED"
	case "pending_deletion":
		return "🗑️  DELETING"
	default:
		return strings.ToUpper(status)
	}
}

// formatStatusSimple formats the status without emojis for table alignment
func formatStatusSimple(status string) string {
	return strings.ToUpper(status)
}

// getStatusEmoji returns appropriate emoji for status
func getStatusEmoji(status string) string {
	switch strings.ToLower(status) {
	case "approved":
		return "✅"
	case "pending":
		return "⏳"
	case "rejected":
		return "❌"
	case "disabled":
		return "🚫"
	case "pending_deletion":
		return "🗑️"
	default:
		return "📄"
	}
}