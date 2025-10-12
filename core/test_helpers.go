package core

import (
	"time"

	ctxtime "github.com/m87/ctx/time"
)

var TEST_ID = "test-context"
var TEST_ID_2 = "test-context-2"
var TEST_INTERVAL_ID = "test-interval-1"
var TEST_INTERVAL_2_ID = "test-interval-2"
var TEST_INTERVAL_3_ID = "test-interval-3"

type TestTimeProvider struct {
	currentTime time.Time
}

func (t *TestTimeProvider) Now() ctxtime.ZonedTime {
	return ctxtime.ZonedTime{
		Time:     t.currentTime,
		Timezone: time.UTC.String(),
	}
}

func (t *TestTimeProvider) SetCurrentTimeFromString(timeStr string) error {
	parsedTime, err := time.Parse(time.DateTime, timeStr)
	if err != nil {
		return err
	}
	t.currentTime = parsedTime
	return nil
}

func CreateTestSession() *Session {
	dt := time.Date(2025, 2, 2, 12, 12, 12, 0, time.UTC)
	return &Session{
		TimeProvider: &TestTimeProvider{currentTime: dt},
		State: &State{
			Contexts: map[string]Context{
				"test-context": {
					Id:          TEST_ID,
					Description: "Test Context",
					Duration:    2 * time.Hour,
					Labels:      []string{"test1-2", "test1-1"},
					Intervals: map[string]Interval{
						"test-interval-1": {
							Id:       "test-interval",
							Start:    ctxtime.ZonedTime{Time: dt, Timezone: "UTC"},
							End:      ctxtime.ZonedTime{Time: dt.Add(1 * time.Hour), Timezone: "UTC"},
							Duration: 1 * time.Hour,
						},
						"test-interval-2": {
							Id:       "test-interval-2",
							Start:    ctxtime.ZonedTime{Time: dt.Add(1 * time.Hour), Timezone: "UTC"},
							End:      ctxtime.ZonedTime{Time: dt.Add(3 * time.Hour), Timezone: "UTC"},
							Duration: 1 * time.Hour,
						},
						"test-interval-3": {
							Id:       "test-interval-2",
							Start:    ctxtime.ZonedTime{Time: dt.Add(1 * time.Hour), Timezone: "UTC"},
							End:      ctxtime.ZonedTime{Time: time.Time{}, Timezone: "UTC"},
							Duration: 1 * time.Hour,
						},
					},
				},
				"test-context-2": {
					Id:          TEST_ID_2,
					Description: "Test2 Context",
					Duration:    2 * time.Hour,
					Labels:      []string{"test2-2", "test2-1"},
					Intervals: map[string]Interval{
						"test-interval-2-1": {
							Id:       "test-interval-2-1",
							Start:    ctxtime.ZonedTime{Time: dt, Timezone: "UTC"},
							End:      ctxtime.ZonedTime{Time: dt.Add(1 * time.Hour), Timezone: "UTC"},
							Duration: 1 * time.Hour,
						},
						"test-interval-2-2": {
							Id:       "test-interval-2-2",
							Start:    ctxtime.ZonedTime{Time: dt.Add(1 * time.Hour), Timezone: "UTC"},
							End:      ctxtime.ZonedTime{Time: dt.Add(3 * time.Hour), Timezone: "UTC"},
							Duration: 1 * time.Hour,
						},
					},
				},
			},
		},
	}
}
