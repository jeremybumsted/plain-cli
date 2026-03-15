package util

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// ParseDateToISO8601 parses various date formats and converts them to ISO8601/RFC3339 format in UTC.
//
// Supported formats:
//   - ISO8601/RFC3339: "2026-03-15T00:00:00Z" (pass through)
//   - Relative dates: "7d", "1w", "2M", "1y" (days, weeks, months, years ago)
//   - Human-readable: "yesterday", "today", "last-week", "last-month", "last-year"
//
// Returns the date in RFC3339 format (ISO8601 compliant) in UTC.
func ParseDateToISO8601(dateStr string) (string, error) {
	if dateStr == "" {
		return "", fmt.Errorf("date string cannot be empty")
	}

	// Trim whitespace
	dateStr = strings.TrimSpace(dateStr)

	// Try parsing as ISO8601/RFC3339 first (before lowercasing)
	if parsedTime, err := time.Parse(time.RFC3339, dateStr); err == nil {
		return parsedTime.UTC().Format(time.RFC3339), nil
	}

	// Try parsing as ISO8601 date only (YYYY-MM-DD)
	if parsedTime, err := time.Parse("2006-01-02", dateStr); err == nil {
		return parsedTime.UTC().Format(time.RFC3339), nil
	}

	// Now normalize to lowercase for human-readable and relative formats
	dateStr = strings.ToLower(dateStr)

	now := time.Now().UTC()

	// Handle human-readable dates
	switch dateStr {
	case "today":
		return startOfDay(now).Format(time.RFC3339), nil
	case "yesterday":
		return startOfDay(now.AddDate(0, 0, -1)).Format(time.RFC3339), nil
	case "last-week":
		return startOfDay(now.AddDate(0, 0, -7)).Format(time.RFC3339), nil
	case "last-month":
		return startOfDay(now.AddDate(0, -1, 0)).Format(time.RFC3339), nil
	case "last-year":
		return startOfDay(now.AddDate(-1, 0, 0)).Format(time.RFC3339), nil
	}

	// Handle relative dates (e.g., "7d", "1w", "2M", "1y")
	// Note: input is already lowercased at this point, so we match lowercase 'm'
	relativePattern := regexp.MustCompile(`^(\d+)(d|w|m|y)$`)
	matches := relativePattern.FindStringSubmatch(dateStr)
	if len(matches) == 3 {
		value, err := strconv.Atoi(matches[1])
		if err != nil {
			return "", fmt.Errorf("invalid numeric value in relative date: %s", dateStr)
		}

		unit := matches[2]
		var targetTime time.Time

		switch unit {
		case "d": // days
			targetTime = now.AddDate(0, 0, -value)
		case "w": // weeks
			targetTime = now.AddDate(0, 0, -value*7)
		case "m": // months (lowercase after normalization)
			targetTime = now.AddDate(0, -value, 0)
		case "y": // years
			targetTime = now.AddDate(-value, 0, 0)
		default:
			return "", fmt.Errorf("unsupported time unit: %s", unit)
		}

		return startOfDay(targetTime).Format(time.RFC3339), nil
	}

	return "", fmt.Errorf("unsupported date format: %s (supported formats: ISO8601, relative dates like '7d/1w/2M/1y', or human-readable like 'yesterday/today/last-week/last-month/last-year')", dateStr)
}

// ValidateDateRange validates that the 'after' date is before the 'before' date.
// Both dates are parsed using ParseDateToISO8601 before comparison.
//
// Returns an error if:
//   - Either date cannot be parsed
//   - The 'after' date is greater than or equal to the 'before' date
func ValidateDateRange(after, before string) error {
	if after == "" && before == "" {
		return nil // Both empty is valid (no filtering)
	}

	var afterTime, beforeTime time.Time

	if after != "" {
		afterISO, err := ParseDateToISO8601(after)
		if err != nil {
			return fmt.Errorf("invalid 'after' date: %w", err)
		}
		afterTime, err = time.Parse(time.RFC3339, afterISO)
		if err != nil {
			return fmt.Errorf("failed to parse 'after' date: %w", err)
		}
	}

	if before != "" {
		beforeISO, err := ParseDateToISO8601(before)
		if err != nil {
			return fmt.Errorf("invalid 'before' date: %w", err)
		}
		beforeTime, err = time.Parse(time.RFC3339, beforeISO)
		if err != nil {
			return fmt.Errorf("failed to parse 'before' date: %w", err)
		}
	}

	// Only validate range if both dates are provided
	if after != "" && before != "" {
		if afterTime.After(beforeTime) || afterTime.Equal(beforeTime) {
			return fmt.Errorf("'after' date (%s) must be before 'before' date (%s)", after, before)
		}
	}

	return nil
}

// startOfDay returns the start of the day (00:00:00) for the given time in UTC
func startOfDay(t time.Time) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
}
