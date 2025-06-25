package datetime

import (
	"fmt"
	"time"
)

// Parser defines the interface for date parsing
type Parser interface {
	ParseToEpoch(dateStr string) (int64, error)
}

// ISO8601Parser implements the Parser interface for ISO-8601 dates
type ISO8601Parser struct{}

// NewISO8601Parser creates a new ISO-8601 date parser
func NewISO8601Parser() *ISO8601Parser {
	return &ISO8601Parser{}
}

// ParseToEpoch converts an ISO-8601 date string to Unix epoch
func (p *ISO8601Parser) ParseToEpoch(dateStr string) (int64, error) {
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

// ConvertEpochToLocal converts Unix epoch to local time in the specified timezone
func ConvertEpochToLocal(epoch int64, loc *time.Location) time.Time {
	return time.Unix(epoch, 0).In(loc)
}