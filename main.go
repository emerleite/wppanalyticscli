package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
	_ "time/tzdata" // Embed timezone data for Windows compatibility

	"wppanalyticscli/internal/api"
	"wppanalyticscli/internal/config"
	"wppanalyticscli/internal/datetime"
	"wppanalyticscli/internal/formatter"
	"wppanalyticscli/internal/input"
)

func main() {
	var wbaID = flag.String("wbaid", "", "WBA ID (required)")
	var startDate = flag.String("start", "", "Start date in ISO-8601 format: YYYY-MM-DD or YYYY-MM-DDTHH:MM:SSZ (required)")
	var endDate = flag.String("end", "", "End date in ISO-8601 format: YYYY-MM-DD or YYYY-MM-DDTHH:MM:SSZ (required)")
	var granularity = flag.String("granularity", "DAY", "Granularity: HALF_HOUR, DAY, or MONTH (for analytics) / daily (for templates)")
	var timezone = flag.String("timezone", "America/Sao_Paulo", "Timezone for date display (default: America/Sao_Paulo)")
	
	// Template analytics specific flags
	var mode = flag.String("mode", "analytics", "Mode: analytics, template, or list-templates")
	var metricTypes = flag.String("metrics", "", "Comma-separated metric types for templates (cost,clicked,delivered,read,sent)")
	var templateIDs = flag.String("templates", "", "Comma-separated template IDs for template analytics")
	
	// Template listing specific flags
	var limit = flag.Int("limit", 25, "Number of templates to retrieve (default: 25)")
	var after = flag.String("after", "", "Pagination cursor for next page")
	
	flag.Parse()

	// Parse template-specific parameters
	var metricTypesList []string
	var templateIDsList []string
	
	if *metricTypes != "" {
		metricTypesList = strings.Split(*metricTypes, ",")
		// Trim whitespace
		for i, metric := range metricTypesList {
			metricTypesList[i] = strings.TrimSpace(metric)
		}
	}
	
	if *templateIDs != "" {
		templateIDsList = strings.Split(*templateIDs, ",")
		// Trim whitespace
		for i, template := range templateIDsList {
			templateIDsList[i] = strings.TrimSpace(template)
		}
	}

	// Create configuration
	cfg := &config.Config{
		WBAID:       *wbaID,
		StartDate:   *startDate,
		EndDate:     *endDate,
		Granularity: *granularity,
		Timezone:    *timezone,
		Mode:        *mode,
		MetricTypes: metricTypesList,
		TemplateIDs: templateIDsList,
		Limit:       *limit,
		After:       *after,
	}

	// Basic parameter validation
	if cfg.WBAID == "" {
		fmt.Fprintf(os.Stderr, "Usage: %s -wbaid=<id> [options]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\nBasic Analytics:\n")
		fmt.Fprintf(os.Stderr, "  %s -wbaid=123 -start=2025-06-20 -end=2025-06-24\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\nTemplate Analytics:\n")
		fmt.Fprintf(os.Stderr, "  %s -mode=template -wbaid=123 -start=2025-06-20 -end=2025-06-24 -templates=1026573095658757 -metrics=cost,clicked,delivered,read,sent\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\nList Templates:\n")
		fmt.Fprintf(os.Stderr, "  %s -mode=list-templates -wbaid=123\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\nDate formats: YYYY-MM-DD or YYYY-MM-DDTHH:MM:SSZ\n")
		os.Exit(1)
	}

	// Load timezone with fallback
	loc, err := time.LoadLocation(cfg.Timezone)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: Could not load timezone '%s': %v\n", cfg.Timezone, err)
		fmt.Fprintf(os.Stderr, "Falling back to UTC timezone\n")
		loc = time.UTC
	}

	// Load access token
	prompter := input.NewSecurePrompter()
	accessToken, err := config.LoadAccessToken(prompter.PromptForToken)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading access token: %v\n", err)
		os.Exit(1)
	}
	cfg.AccessToken = accessToken

	// Validate configuration
	validator := config.NewConfigValidator()
	if err := validator.Validate(cfg); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Parse dates (only for analytics and template modes)
	var startEpoch, endEpoch int64
	if cfg.Mode != "list-templates" {
		parser := datetime.NewISO8601Parser()
		var err error
		startEpoch, err = parser.ParseToEpoch(cfg.StartDate)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing start date: %v\n", err)
			os.Exit(1)
		}

		endEpoch, err = parser.ParseToEpoch(cfg.EndDate)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing end date: %v\n", err)
			os.Exit(1)
		}
	}

	// Create API client
	apiClient := api.NewFacebookGraphClient()
	
	// Handle different modes
	if cfg.Mode == "list-templates" {
		// Make template list request
		listResponse, err := apiClient.ListTemplates(cfg.WBAID, cfg.AccessToken, cfg.Limit, cfg.After)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error listing templates: %v\n", err)
			os.Exit(1)
		}

		// Format and display list output
		listFormatter := formatter.NewListFormatter()
		result := listFormatter.FormatList(listResponse)
		fmt.Print(result)
	} else if cfg.Mode == "template" {
		// Make template analytics request
		templateResponse, err := apiClient.GetTemplateAnalytics(cfg.WBAID, startEpoch, endEpoch, cfg.Granularity, cfg.MetricTypes, cfg.TemplateIDs, cfg.AccessToken)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error making template request: %v\n", err)
			os.Exit(1)
		}

		// Format and display template output
		templateFormatter := formatter.NewTemplateFormatter()
		result := templateFormatter.FormatTemplate(templateResponse, loc)
		fmt.Print(result)
	} else {
		// Make regular analytics request
		response, err := apiClient.GetAnalytics(cfg.WBAID, startEpoch, endEpoch, cfg.Granularity, cfg.AccessToken)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error making request: %v\n", err)
			os.Exit(1)
		}

		// Format and display output
		outputFormatter := formatter.NewTableFormatter()
		result := outputFormatter.Format(response, loc)
		fmt.Print(result)
	}
}

