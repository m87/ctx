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
	v := strings.TrimSpace(value)
	// try RFC3339 first
	parsed, err := time.Parse(time.RFC3339, v)
	if err == nil {
		loc := parsed.Location()
		if loc == nil {
			loc = time.UTC
		}
		return core.ZonedTime{Time: parsed, Timezone: loc.String(), IsZero: false}, nil
	}

	// try local "YYYY-MM-DD HH:MM:SS" format
	locName := core.DetectTimezoneName()
	loc, lerr := time.LoadLocation(locName)
	if lerr != nil {
		loc = time.UTC
	}
	parsed2, err2 := time.ParseInLocation("2006-01-02 15:04:05", v, loc)
	if err2 == nil {
		return core.ZonedTime{Time: parsed2, Timezone: loc.String(), IsZero: false}, nil
	}

	return core.ZonedTime{}, fmt.Errorf("invalid datetime %q, expected RFC3339 or 'YYYY-MM-DD HH:MM:SS'", value)
}
