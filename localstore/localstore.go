package localstore

import (
	"encoding/json"
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
