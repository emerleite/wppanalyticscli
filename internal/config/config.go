package config

import (
	"fmt"
	"os"
)

// Config holds the application configuration
type Config struct {
	WBAID       string
	StartDate   string
	EndDate     string
	Granularity string
	Timezone    string
	AccessToken string
	// Template analytics specific fields
	Mode         string   // "analytics" or "template"
	MetricTypes  []string // For template analytics
	TemplateIDs  []string // For template analytics
}

// Validator defines the interface for configuration validation
type Validator interface {
	Validate(config *Config) error
}

// ConfigValidator implements the Validator interface
type ConfigValidator struct{}

// NewConfigValidator creates a new configuration validator
func NewConfigValidator() *ConfigValidator {
	return &ConfigValidator{}
}

// Validate validates the configuration
func (v *ConfigValidator) Validate(config *Config) error {
	if config.WBAID == "" {
		return fmt.Errorf("WBA ID is required")
	}
	
	if config.StartDate == "" {
		return fmt.Errorf("start date is required")
	}
	
	if config.EndDate == "" {
		return fmt.Errorf("end date is required")
	}
	
	if config.Mode == "template" {
		if !isValidTemplateGranularity(config.Granularity) {
			return fmt.Errorf("granularity for templates must be daily")
		}
		
		if len(config.TemplateIDs) == 0 {
			return fmt.Errorf("template IDs are required for template analytics")
		}
		
		if len(config.MetricTypes) == 0 {
			return fmt.Errorf("metric types are required for template analytics")
		}
	} else {
		if !isValidGranularity(config.Granularity) {
			return fmt.Errorf("granularity must be HALF_HOUR, DAY, or MONTH")
		}
	}
	
	if config.AccessToken == "" {
		return fmt.Errorf("access token is required")
	}
	
	return nil
}

// isValidGranularity validates the granularity value
func isValidGranularity(g string) bool {
	switch g {
	case "HALF_HOUR", "DAY", "MONTH":
		return true
	default:
		return false
	}
}

// isValidTemplateGranularity validates the granularity value for templates
func isValidTemplateGranularity(g string) bool {
	switch g {
	case "daily":
		return true
	default:
		return false
	}
}

// LoadAccessToken loads the access token from environment or prompts for it
func LoadAccessToken(promptFunc func() (string, error)) (string, error) {
	accessToken := os.Getenv("FB_ACCESS_TOKEN")
	if accessToken != "" {
		return accessToken, nil
	}
	
	return promptFunc()
}