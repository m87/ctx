package localstore

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"path/filepath"

	"github.com/m87/ctx/ctx_model"
	"github.com/spf13/viper"
)

func LoadState() ctx_model.State {
	statePath := filepath.Join(viper.GetString("ctxPath"), "state")
	data, err := os.ReadFile(statePath)
	if err != nil {
		log.Fatal("Unable to read state file")
	}

	state := ctx_model.State{}
	err = json.Unmarshal(data, &state)
	if err != nil {
		log.Fatal("Unable to parse state file")
	}

	return state
}

func SaveState(state *ctx_model.State) {
	statePath := filepath.Join(viper.GetString("ctxPath"), "state")
	data, err := json.Marshal(state)
	if err != nil {
		panic(err)
	}
	os.WriteFile(statePath, data, 0644)
}

type LocalContextStore struct {
	path string
}

type LocalEventsStore struct {
	path string
}

type LocalArchiveStore struct {
  path string
}

func NewContextStore(path string) *LocalContextStore {
	return &LocalContextStore{
		path: path,
	}
}

func NewEventsStore(path string) *LocalEventsStore {
	return &LocalEventsStore{
		path: path,
	}
}

func NewArchiveStore(path string) *LocalArchiveStore {
  return &LocalArchiveStore {
    path: path,
  }
}

func (store *LocalArchiveStore) saveArchive(entry *ctx_model.ArchiveEntry, path string) error{
  data, err := json.Marshal(entry)
  if err != nil {
    return errors.New("unable to marshal archive for " + entry.Context.Id )
  }

  os.WriteFile(path, data, 0644)

  return nil
}

func (store *LocalArchiveStore) loadArchive(path string) (*ctx_model.ArchiveEntry, error) {
  data, err := os.ReadFile(path)

  if err != nil {
    return nil, errors.New("unable to read archive file " + path)
  }

  entry := ctx_model.ArchiveEntry{}
  err = json.Unmarshal(data, &entry)
  
  if err != nil {
    return nil, errors.New("unable to parse archive file " + path)
  }

  return &entry, nil

}

func (store *LocalArchiveStore) updateArchive(entry *ctx_model.ArchiveEntry, path string) error {
  entry2Update, err := store.loadArchive(path)

  if err != nil {
    return err
  }

  if entry2Update.Context.Id != entry.Context.Id {
    return errors.New("contexts mismatch, entry to update: " + entry2Update.Context.Id + ", entry to archive: " + entry.Context.Id)
  }

  entry2Update.Context.Duration = entry2Update.Context.Duration + entry.Context.Duration 
  entry2Update.Context.Comments = append(entry2Update.Context.Comments, entry.Context.Comments...)
  entry2Update.Context.Intervals = append(entry2Update.Context.Intervals, entry.Context.Intervals...)
  entry2Update.Context.State = entry.Context.State
  entry2Update.Events = append(entry2Update.Events, entry.Events...)


  return store.saveArchive(entry2Update, path)
}


func (store *LocalArchiveStore) UpsertArchive(entry *ctx_model.ArchiveEntry) error {
  path := filepath.Join(store.path, "archive", entry.Context.Id + ".ctx")
  if _, err := os.Stat(path); err == nil {
    return store.updateArchive(entry, path)
  } else {
    return store.saveArchive(entry, path)
  }
}

func (store *LocalArchiveStore) UpsertEventsArchive(events []ctx_model.Event) error {
  return nil
}

func (store *LocalContextStore) Apply(fn ctx_model.StatePatch) error {
	state := LoadState()
	err := fn(&state)
	if err != nil {
		return err
	} else {
		SaveState(&state)
		return nil
	}
}

func (store *LocalContextStore) Read(fn ctx_model.StatePatch) {
	state := LoadState()
	fn(&state)
}

func LoadEvents() ctx_model.EventRegistry {
	eventsPath := filepath.Join(viper.GetString("ctxPath"), "events")
	data, err := os.ReadFile(eventsPath)
	if err != nil {
		log.Fatal("Unable to read state file")
	}

	events := ctx_model.EventRegistry{}
	err = json.Unmarshal(data, &events)
	if err != nil {
		log.Fatal("Unable to parse state file")
	}

	return events
}

func SaveEvents(eventsRegistry *ctx_model.EventRegistry) {
	eventsPath := filepath.Join(viper.GetString("ctxPath"), "events")
	data, err := json.Marshal(eventsRegistry)
	if err != nil {
		panic(err)
	}
	os.WriteFile(eventsPath, data, 0644)
}

func (store *LocalEventsStore) Apply(fn ctx_model.EventsPatch) error {
	events := LoadEvents()
	err := fn(&events)
	if err != nil {
		return err
	} else {
		SaveEvents(&events)
		return nil
	}
}

func (store *LocalEventsStore) Read(fn ctx_model.EventsPatch) {
	events := LoadEvents()
	fn(&events)
}
