package server

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

const timeZoneQueryParameter = "timeZone"

func parseRequestedDay(r *http.Request) (time.Time, error) {
	zoneName := strings.TrimSpace(r.URL.Query().Get(timeZoneQueryParameter))
	if zoneName == "" {
		zoneName = "UTC"
	}

	location, err := time.LoadLocation(zoneName)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid time zone %q", zoneName)
	}

	date, err := time.ParseInLocation("2006-01-02", r.PathValue("date"), location)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid date, expected YYYY-MM-DD")
	}
	return date, nil
}
