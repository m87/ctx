package localstorage

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"path/filepath"

	"github.com/m87/ctx/ctx_model"
	"github.com/m87/ctx/storage"
	"github.com/spf13/viper"
)

func LoadState() ctx_model.State {
	statePath := filepath.Join(viper.GetString("storePath"), "state")
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
	statePath := filepath.Join(viper.GetString("storePath"), "state")
	data, err := json.Marshal(state)
	if err != nil {
		panic(err)
	}
	os.WriteFile(statePath, data, 0777)
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
	return &LocalArchiveStore{
		path: path,
	}
}

func (store *LocalArchiveStore) saveArchive(entry *ctx_model.ContextArchive, path string) error {
	data, err := json.Marshal(entry)
	if err != nil {
		return errors.New("unable to marshal archive for " + entry.Context.Id)
	}

	os.WriteFile(path, data, 0777)

	return nil
}

func (store *LocalArchiveStore) saveEventsArchive(entry *ctx_model.EventsArchive, path string) error {
	data, err := json.Marshal(entry)
	if err != nil {
		return errors.New("unable to marshal events archive for " + path)
	}

	os.WriteFile(path, data, 0777)

	return nil
}

func (store *LocalArchiveStore) loadArchive(id string, path string) (*ctx_model.ContextArchive, error) {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return &ctx_model.ContextArchive{
				Context: ctx_model.Context{
					Id: id,
				},
			}, nil
		} else {
			return nil, errors.New("unable to read archive file " + path)
		}
	}

	data, err := os.ReadFile(path)

	if err != nil {
		return nil, errors.New("unable to read archive file " + path)
	}

	entry := ctx_model.ContextArchive{}
	err = json.Unmarshal(data, &entry)

	if err != nil {
		return nil, errors.New("unable to parse archive file " + path)
	}

	return &entry, nil

}

func (store *LocalArchiveStore) loadEventsArchive(path string) (*ctx_model.EventsArchive, error) {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return &ctx_model.EventsArchive{
				Events: []ctx_model.Event{},
			}, nil
		} else {
			return nil, errors.New("unable to read eventsarchive file " + path)
		}
	}
	data, err := os.ReadFile(path)

	if err != nil {
		return nil, errors.New("unable to read events archive file " + path)
	}

	entry := ctx_model.EventsArchive{}
	err = json.Unmarshal(data, &entry)

	if err != nil {
		return nil, errors.New("unable to parse events archive file " + path)
	}

	return &entry, nil

}

func (store *LocalArchiveStore) Apply(id string, fn storage.ArchivePatch) error {
	path := filepath.Join(store.path, "archive", id+".ctx")
	entry, err := store.loadArchive(id, path)

	if err != nil {
		return err
	}

	if err := fn(entry); err != nil {
		return err
	}

	return store.saveArchive(entry, path)
}

func (store *LocalArchiveStore) ApplyEvents(date string, fn storage.ArchiveEventsPatch) error {
	path := filepath.Join(store.path, "archive", date+".events")
	events, err := store.loadEventsArchive(path)

	if err != nil {
		return err
	}

	if err := fn(events); err != nil {
		return err
	} else {
		return store.saveEventsArchive(events, path)
	}
}

func (store *LocalContextStore) Apply(fn storage.StatePatch) error {
	state := LoadState()
	err := fn(&state)
	if err != nil {
		return err
	} else {
		SaveState(&state)
		return nil
	}
}

func (store *LocalContextStore) Read(fn storage.StatePatch) error {
	state := LoadState()
	return fn(&state)
}

func LoadEvents() ctx_model.EventRegistry {
	eventsPath := filepath.Join(viper.GetString("storePath"), "events")
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
	eventsPath := filepath.Join(viper.GetString("storePath"), "events")
	data, err := json.Marshal(eventsRegistry)
	if err != nil {
		panic(err)
	}
	os.WriteFile(eventsPath, data, 0777)
}

func (store *LocalEventsStore) Apply(fn storage.EventsPatch) error {
	events := LoadEvents()
	err := fn(&events)
	if err != nil {
		return err
	} else {
		SaveEvents(&events)
		return nil
	}
}

func (store *LocalEventsStore) Read(fn storage.EventsPatch) error {
	events := LoadEvents()
	return fn(&events)
}
