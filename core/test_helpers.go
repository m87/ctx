package core

import (
	"time"

	ctxtime "github.com/m87/ctx/time"
)

var TEST_ID = "test-context"
var TEST_INTERVAL_ID = "test-interval"

func CreateTestSession() *Session {
	dt := time.Date(2025, 2, 2, 12, 12, 12, 0, time.UTC)
	return &Session{
		State: &State{
			Contexts: map[string]Context{
				"test-context": {
					Id:          TEST_ID,
					Description: "Test Context",
					Duration:    2 * time.Hour,
					Intervals: []Interval{
						{
							Id:       "test-interval",
							Start:    ctxtime.ZonedTime{Time: dt, Timezone: "UTC"},
							End:      ctxtime.ZonedTime{Time: dt.Add(1 * time.Hour), Timezone: "UTC"},
							Duration: 1 * time.Hour,
						},
						{
							Id:       "test-interval-2",
							Start:    ctxtime.ZonedTime{Time: dt.Add(1 * time.Hour), Timezone: "UTC"},
							End:      ctxtime.ZonedTime{Time: dt.Add(3 * time.Hour), Timezone: "UTC"},
							Duration: 1 * time.Hour,
						},
					},
				},
			},
		},
		EventsRegistry: &EventRegistry{
			Events: []Event{},
		},
	}
}
