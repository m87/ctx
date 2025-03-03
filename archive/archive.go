package archive

import (
	"errors"
	"time"

	"github.com/m87/ctx/archive_model"
	"github.com/m87/ctx/ctx"
	"github.com/m87/ctx/ctx_model"
	"github.com/m87/ctx/events_model"
)

func Archive(id string, state *ctx_model.State, eventsRegistry *events_model.EventRegistry, entry *archive_model.ArchiveEntry) (map[string][]events_model.Event, error) {
	if id == state.CurrentId {
		return nil, errors.New("context is active")
	}

	context := state.Contexts[id]
	var ctxEvents []events_model.Event
	evnetsByDate := map[string][]events_model.Event{}

	originalEvents := append([]events_model.Event{}, eventsRegistry.Events...)

	for _, event := range eventsRegistry.Events {
		if event.CtxId == id {
			ctxEvents = append(ctxEvents, event)

			date := event.DateTime.Local().Format(time.DateOnly)
			evnetsByDate[date] = append(evnetsByDate[date], event)

			for i, ev := range originalEvents {
				if ev.UUID == event.UUID {
					if len(originalEvents) == i {
						originalEvents = originalEvents[:i]
					} else {
						originalEvents = append(originalEvents[:i], originalEvents[i+1:]...)
					}
					break
				}
			}
		}
	}

	eventsRegistry.Events = originalEvents
	entry.Context.Id = context.Id
	entry.Context.Duration = entry.Context.Duration + context.Duration
	entry.Context.Comments = append(entry.Context.Comments, context.Comments...)
	entry.Context.Intervals = append(entry.Context.Intervals, context.Intervals...)
	entry.Context.State = context.State
	entry.Events = append(entry.Events, ctxEvents...)

	ctx.Delete(id, state)
	return evnetsByDate, nil
}
