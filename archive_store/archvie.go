package archive_store

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/m87/ctx/archive_model"
	"github.com/m87/ctx/events_model"
	"github.com/spf13/viper"
)

func SaveEventsArchive(eventsByDate map[string][]events_model.Event) {
	for d, e := range eventsByDate {
		path := filepath.Join(viper.GetString("ctxPath"), "archive", d+".events")
		savedEvents := loadEvents(path)
		savedEvents = append(savedEvents, e...)
		data, err := json.Marshal(savedEvents)
		if err != nil {
			panic(err)
		}

		os.WriteFile(path, data, 0644)

	}
}

func SaveArchive(entry *archive_model.ArchiveEntry) {
	entryPath := filepath.Join(viper.GetString("ctxPath"), "archive", entry.Context.Id+".ctx")
	data, err := json.Marshal(entry)
	if err != nil {
		panic(err)
	}
	os.WriteFile(entryPath, data, 0644)
}

func LoadArchive(id string) archive_model.ArchiveEntry {
	path := filepath.Join(viper.GetString("ctxPath"), "archive", id+".ctx")
	data, err := os.ReadFile(path)

	if err != nil {
		return archive_model.ArchiveEntry{}
	}

	entry := archive_model.ArchiveEntry{}
	err = json.Unmarshal(data, &entry)

	if err != nil {
		panic("Uanble to parse entry file")
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
		panic("Unable to parse state file")
	}

	return events
}
