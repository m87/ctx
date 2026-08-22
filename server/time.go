package server

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

const timeZoneQueryParameter = "timeZone"

func pareseDateTime(r *http.Request, dateTime string) (time.Time, error) {
	zoneName := strings.TrimSpace(r.URL.Query().Get(timeZoneQueryParameter))
	location, err := parseTimeZone(zoneName)
	if err != nil {
		return time.Time{}, err
	}

	date, err := time.ParseInLocation("2006-01-02T15:04:05", dateTime, location)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid date-time, expected YYYY-MM-DDTHH:MM:SS")
	}
	return date, nil
}

func parseRequestedDay(r *http.Request) (time.Time, error) {
	zoneName := strings.TrimSpace(r.URL.Query().Get(timeZoneQueryParameter))
	location, err := parseTimeZone(zoneName)
	if err != nil {
		return time.Time{}, err
	}

	date, err := time.ParseInLocation("2006-01-02", r.PathValue("date"), location)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid date, expected YYYY-MM-DD")
	}
	return date, nil
}

func parseTimeZone(zoneName string) (*time.Location, error) {
	zoneName = strings.TrimSpace(zoneName)
	if zoneName == "" {
		zoneName = "UTC"
	}

	location, err := time.LoadLocation(zoneName)
	if err != nil {
		return nil, fmt.Errorf("invalid time zone %q", zoneName)
	}
	return location, nil
}
