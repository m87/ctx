package archive

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/m87/ctx/ctx"
	"github.com/m87/ctx/ctx_model"
	"github.com/m87/ctx/events"
	"github.com/spf13/viper"
)

type ArchiveEntry struct {
	Context ctx_model.Context `json:"context"`
	Events  []events.Event    `json:"events"`
}

func Archive(id string, state *ctx_model.State, eventsRegistry *events.EventRegistry) {

	if id == state.CurrentId {
		log.Fatalf("context %s is active", id)
	}

	context := state.Contexts[id]
	var ctxEvents []events.Event
	evnetsByDate := map[string][]events.Event{}

	originalEvents := append([]events.Event{}, eventsRegistry.Events...)

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
	data, err := json.Marshal(
		ArchiveEntry{
			Context: context,
			Events:  ctxEvents,
		},
	)
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

func loadEvents(path string) []events.Event {
	data, err := os.ReadFile(path)
	if err != nil {
		return []events.Event{}
	}

	events := []events.Event{}
	err = json.Unmarshal(data, &events)
	if err != nil {
		log.Fatal("Unable to parse state file")
	}

	return events
}
