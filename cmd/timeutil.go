package cmd

import (
	"fmt"
	"strings"
	"time"

	"github.com/m87/ctx/core"
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

func parseDateTime(value string) (core.ZonedTime, error) {
	parsed, err := time.Parse(time.RFC3339, strings.TrimSpace(value))
	if err != nil {
		return core.ZonedTime{}, fmt.Errorf("invalid datetime %q, expected RFC3339", value)
	}
	location := parsed.Location()
	if location == nil {
		location = time.UTC
	}
	return core.ZonedTime{Time: parsed, Timezone: location.String(), IsZero: false}, nil
}
