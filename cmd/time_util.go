package cmd

import (
	"fmt"
	"strings"
	"time"
)

func parseDay(day string) (time.Time, error) {
	day = strings.TrimSpace(day)
	if day == "" {
		now := time.Now().UTC()
		return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC), nil
	}
	parsed, err := time.Parse("2006-01-02", day)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid day format %q, expected YYYY-MM-DD", day)
	}
	return parsed.UTC(), nil
}

func parseDateTime(value string) (time.Time, error) {
	v := strings.TrimSpace(value)
	parsed, err := time.Parse(time.RFC3339, v)
	if err == nil {
		return parsed.UTC(), nil
	}

	parsed2, err2 := time.ParseInLocation("2006-01-02 15:04:05", v, time.Local)
	if err2 == nil {
		return parsed2.UTC(), nil
	}

	return time.Time{}, fmt.Errorf("invalid datetime %q, expected RFC3339 or 'YYYY-MM-DD HH:MM:SS'", value)
}

func formatIntervalDateTime(value *time.Time) string {
	if value == nil || value.IsZero() {
		return "(ongoing)"
	}
	return value.UTC().Format(time.RFC3339)
}
