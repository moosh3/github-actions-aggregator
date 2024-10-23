package api

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

func parseTimeParameter(param string, defaultTime time.Time) (time.Time, error) {
	if param == "" {
		return defaultTime, nil
	}

	// Handle special keywords
	if param == "now" {
		return time.Now(), nil
	}

	// Check for relative time expressions
	if strings.HasSuffix(param, "_ago") {
		return parseRelativeTimeAgo(param)
	} else if strings.Contains(param, "_") {
		return parseRelativeTime(param)
	} else {
		// Try to parse as an explicit date-time
		t, err := time.Parse(time.RFC3339, param)
		if err != nil {
			return time.Time{}, fmt.Errorf("invalid time format: %v", param)
		}
		return t, nil
	}
}

func parseRelativeTime(param string) (time.Time, error) {
	parts := strings.Split(param, "_")
	if len(parts) != 2 {
		return time.Time{}, fmt.Errorf("invalid relative time format: %v", param)
	}

	value, err := strconv.Atoi(parts[0])
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid number in relative time: %v", param)
	}

	unit := parts[1]

	return calculateRelativeTime(value, unit)
}

func parseRelativeTimeAgo(param string) (time.Time, error) {
	parts := strings.Split(param, "_")
	if len(parts) != 3 || parts[2] != "ago" {
		return time.Time{}, fmt.Errorf("invalid relative time format: %v", param)
	}

	value, err := strconv.Atoi(parts[0])
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid number in relative time: %v", param)
	}

	unit := parts[1]

	return calculateRelativeTime(value, unit)
}

func calculateRelativeTime(value int, unit string) (time.Time, error) {
	switch unit {
	case "minute", "minutes":
		return time.Now().Add(-time.Duration(value) * time.Minute), nil
	case "hour", "hours":
		return time.Now().Add(-time.Duration(value) * time.Hour), nil
	case "day", "days":
		return time.Now().AddDate(0, 0, -value), nil
	case "week", "weeks":
		return time.Now().AddDate(0, 0, -value*7), nil
	case "month", "months":
		return time.Now().AddDate(0, -value, 0), nil
	case "year", "years":
		return time.Now().AddDate(-value, 0, 0), nil
	default:
		return time.Time{}, fmt.Errorf("invalid time unit in relative time: %v", unit)
	}
}
