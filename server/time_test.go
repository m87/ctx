package server

import (
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestParseRequestedDayUsesIANAZone(t *testing.T) {
	request := httptest.NewRequest("GET", "/?timeZone=Asia%2FTokyo", nil)
	request.SetPathValue("date", "2026-08-02")

	date, err := parseRequestedDay(request)

	require.NoError(t, err)
	require.Equal(t, "Asia/Tokyo", date.Location().String())
	require.Equal(t, "2026-08-01T15:00:00Z", date.UTC().Format(time.RFC3339))
}

func TestParseRequestedDayDefaultsToUTC(t *testing.T) {
	request := httptest.NewRequest("GET", "/", nil)
	request.SetPathValue("date", "2026-08-02")

	date, err := parseRequestedDay(request)

	require.NoError(t, err)
	require.Equal(t, time.UTC, date.Location())
}

func TestParseRequestedDayRejectsInvalidTimeZone(t *testing.T) {
	request := httptest.NewRequest("GET", "/?timeZone=Mars%2FOlympus_Mons", nil)
	request.SetPathValue("date", "2026-08-02")

	_, err := parseRequestedDay(request)

	require.ErrorContains(t, err, "invalid time zone")
}
