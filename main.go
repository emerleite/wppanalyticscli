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
	var mode = flag.String("mode", "analytics", "Mode: analytics or template")
	var metricTypes = flag.String("metrics", "", "Comma-separated metric types for templates (cost,clicked,delivered,read,sent)")
	var templateIDs = flag.String("templates", "", "Comma-separated template IDs for template analytics")
	
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
	}

	// Basic parameter validation
	if cfg.WBAID == "" || cfg.StartDate == "" || cfg.EndDate == "" {
		fmt.Fprintf(os.Stderr, "Usage: %s -wbaid=<id> -start=<date> -end=<date> [options]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\nBasic Analytics:\n")
		fmt.Fprintf(os.Stderr, "  %s -wbaid=123 -start=2025-06-20 -end=2025-06-24\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\nTemplate Analytics:\n")
		fmt.Fprintf(os.Stderr, "  %s -mode=template -wbaid=123 -start=2025-06-20 -end=2025-06-24 -templates=1026573095658757 -metrics=cost,clicked,delivered,read,sent\n", os.Args[0])
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

	// Parse dates
	parser := datetime.NewISO8601Parser()
	startEpoch, err := parser.ParseToEpoch(cfg.StartDate)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing start date: %v\n", err)
		os.Exit(1)
	}

	endEpoch, err := parser.ParseToEpoch(cfg.EndDate)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing end date: %v\n", err)
		os.Exit(1)
	}

	// Create API client
	apiClient := api.NewFacebookGraphClient()
	
	// Handle different modes
	if cfg.Mode == "template" {
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

