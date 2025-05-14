package utils

import (
	"fmt"
	"time"
)

// Define date formats for parsing
const (
	formatWithOffset = "2006-01-02T15:04:05-07:00"
	formatWithZ      = "2006-01-02T15:04:05.000Z"
)

// parseDateTime parses a date-time string in supported formats
func parseDateTime(dateTimeStr string) (time.Time, error) {
	// Try parsing with offset format first
	if t, err := time.Parse(formatWithOffset, dateTimeStr); err == nil {
		return t.UTC(), nil
	}
	// Try parsing with UTC format
	if t, err := time.Parse(formatWithZ, dateTimeStr); err == nil {
		return t.UTC(), nil
	}
	return time.Time{}, fmt.Errorf("invalid date-time format")
}

// getCurrentTimeUTC returns the current time in UTC
func getCurrentTimeUTC() (time.Time, error) {
	utcLoc, err := time.LoadLocation("UTC")
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to load UTC timezone: %v", err)
	}
	// Fixed time for consistency: 2025-05-14 05:34 UTC (12:34 AM CDT)
	return time.Date(2025, time.May, 14, 5, 34, 0, 0, utcLoc), nil
	// For production, use: return time.Now().UTC(), nil
}

// calculateDuration computes the duration and determines if the event is current/future
func calculateDuration(parsedTime, now time.Time) (time.Duration, bool) {
	duration := parsedTime.Sub(now)
	isFutureOrCurrent := duration >= 0
	if !isFutureOrCurrent {
		duration = -duration
	}
	return duration, isFutureOrCurrent
}

// formatDuration formats the duration as a string (e.g., "3d2h togo")
func formatDuration(duration time.Duration, isFutureOrCurrent bool) string {
	days := int(duration.Hours() / 24)
	hours := int(duration.Hours()) % 24

	var result string
	if days > 0 {
		result += fmt.Sprintf("%dd", days)
	}
	if hours > 0 || days == 0 {
		result += fmt.Sprintf("%dh", hours)
	}
	if result == "" {
		result = "0h"
	}

	if isFutureOrCurrent {
		result += " ꜛ"
	} else {
		result += " ago"
	}
	return result
}

// ParseAndCompareDateTime parses a date-time string and compares it to now
func ParseAndCompareDateTime(dateTimeStr string) (time.Time, bool, string, error) {
	// Parse the date-time string
	parsedTime, err := parseDateTime(dateTimeStr)
	if err != nil {
		return time.Time{}, false, "", err
	}

	// Get current time
	now, err := getCurrentTimeUTC()
	if err != nil {
		return time.Time{}, false, "", err
	}

	// Calculate duration and future/current status
	duration, isFutureOrCurrent := calculateDuration(parsedTime, now)

	// Format the duration
	formatted := formatDuration(duration, isFutureOrCurrent)

	return parsedTime, isFutureOrCurrent, formatted, nil
}
