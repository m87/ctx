package archive

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/m87/ctx/ctx"
	"github.com/m87/ctx/ctx_model"
	"github.com/m87/ctx/events_model"
	"github.com/spf13/viper"
)

type ArchiveEntry struct {
	Context ctx_model.Context    `json:"context"`
	Events  []events_model.Event `json:"events"`
}

func Archive(id string, state *ctx_model.State, eventsRegistry *events_model.EventRegistry) {

	if id == state.CurrentId {
		log.Fatalf("context %s is active", id)
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
	entryPath := filepath.Join(viper.GetString("ctxPath"), "archive", id+".ctx")

	entry := loadArchive(entryPath)

	entry.Context.Duration = entry.Context.Duration + context.Duration
	entry.Context.Comments = append(entry.Context.Comments, context.Comments...)
	entry.Context.Intervals = append(entry.Context.Intervals, context.Intervals...)
	entry.Context.State = context.State
	entry.Events = append(entry.Events, ctxEvents...)

	data, err := json.Marshal(entry)
	if err != nil {
		panic(err)
	}
	os.WriteFile(entryPath, data, 0644)

	for d, e := range evnetsByDate {
		path := filepath.Join(viper.GetString("ctxPath"), "archive", d+".events")
		savedEvents := loadEvents(path)
		savedEvents = append(savedEvents, e...)
		data, err := json.Marshal(savedEvents)
		if err != nil {
			panic(err)
		}

		os.WriteFile(path, data, 0644)

	}

	ctx.Delete(id, state)

}

func loadArchive(path string) ArchiveEntry {
	data, err := os.ReadFile(path)

	if err != nil {
		return ArchiveEntry{}
	}

	entry := ArchiveEntry{}
	err = json.Unmarshal(data, &entry)

	if err != nil {
		log.Fatal("Uanble to parse entry file")
	}

	return entry
}

func loadEvents(path string) []events_model.Event {
	data, err := os.ReadFile(path)
	if err != nil {
		return []events_model.Event{}
	}

	events := []events_model.Event{}
	err = json.Unmarshal(data, &events)
	if err != nil {
		log.Fatal("Unable to parse state file")
	}

	return events
}
