package main

import (
	"flag"
	"fmt"
	"os"
	"time"

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
	var granularity = flag.String("granularity", "DAY", "Granularity: HALF_HOUR, DAY, or MONTH")
	var timezone = flag.String("timezone", "America/Sao_Paulo", "Timezone for date display (default: America/Sao_Paulo)")
	
	flag.Parse()

	// Create configuration
	cfg := &config.Config{
		WBAID:       *wbaID,
		StartDate:   *startDate,
		EndDate:     *endDate,
		Granularity: *granularity,
		Timezone:    *timezone,
	}

	// Basic parameter validation
	if cfg.WBAID == "" || cfg.StartDate == "" || cfg.EndDate == "" {
		fmt.Fprintf(os.Stderr, "Usage: %s -wbaid=<id> -start=<date> -end=<date> [-granularity=<HALF_HOUR|DAY|MONTH>] [-timezone=<timezone>]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Date formats: YYYY-MM-DD or YYYY-MM-DDTHH:MM:SSZ\n")
		os.Exit(1)
	}

	// Load timezone
	loc, err := time.LoadLocation(cfg.Timezone)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading timezone '%s': %v\n", cfg.Timezone, err)
		os.Exit(1)
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

	// Create API client and make request
	apiClient := api.NewFacebookGraphClient()
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

