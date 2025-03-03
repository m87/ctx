package archive

import (
	"testing"
	"time"

	"github.com/m87/ctx/archive_model"
	"github.com/m87/ctx/assert"
	"github.com/m87/ctx/ctx"
	"github.com/m87/ctx/ctx_model"
	"github.com/m87/ctx/events"
	"github.com/m87/ctx/events_model"
	"github.com/m87/ctx/test"
)

func TestArchiveContextWithEvents(t *testing.T) {
	state := ctx_model.State{
		Contexts:  map[string]ctx_model.Context{},
		CurrentId: "",
	}
	eventsRegistry := events_model.EventRegistry{}

	ctx.Create(&state, test.TestId, test.TestDescription)
	dateTime, _ := time.Parse(time.DateTime, "2025-02-02 12:25:23")
	events.Publish(test.CreateTestEvent("test", test.TestId, dateTime), &eventsRegistry)

	entry := archive_model.ArchiveEntry{}

	eventsByDate, err := Archive(test.TestId, &state, &eventsRegistry, &entry)

	assert.NoErr(t, err)
	assert.Equal(t, entry.Context.Id, test.TestId)
	assert.Size(t, entry.Events, 1)
	assert.Equal(t, entry.Events[0].Description, "test")
	assert.IsNotNil(t, eventsByDate["2025-02-02"])
	assert.Size(t, eventsByDate["2025-02-02"], 1)
	assert.Equal(t, eventsByDate["2025-02-02"][0].Description, "test")
}

func TestArchiveContextWithoutEvents(t *testing.T) {
	state := ctx_model.State{
		Contexts:  map[string]ctx_model.Context{},
		CurrentId: "",
	}
	eventsRegistry := events_model.EventRegistry{}

	ctx.Create(&state, test.TestId, test.TestDescription)

	entry := archive_model.ArchiveEntry{}

	eventsByDate, err := Archive(test.TestId, &state, &eventsRegistry, &entry)

	assert.NoErr(t, err)
	assert.Equal(t, entry.Context.Id, test.TestId)
	assert.Size(t, entry.Events, 0)
	assert.Equal(t, len(eventsByDate), 0)
}

func TestShouldAppenToExistingContextAndEvents(t *testing.T) {
}

func TestErrorOnCurrentContext(t *testing.T) {

}

func TestArchiveAllContexts(t *testing.T) {

}

func TestSkipCurrentContext(t *testing.T) {

}
