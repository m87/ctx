package localstorage

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"path/filepath"

	"github.com/m87/ctx/core"
	ctxtime "github.com/m87/ctx/time"
	"github.com/spf13/viper"
)

type fileTx[T any] struct {
	path string
	data T
}

type FileStore[T any] struct {
	path string
}

func (store *FileStore[T]) Begin() (core.Tx[T], error) {
	return &fileTx[T]{
		path: store.path,
		data: Load[T](store.path),
	}, nil
}

func (store *FileStore[T]) BeginAndGet() (core.Tx[T], *T, error) {
	tx, err := store.Begin()
	if err != nil {
		return nil, nil, err
	}

	data, err := tx.Get()
	if err != nil {
		return nil, nil, err
	}

	return tx, data, nil
}

func (tx *fileTx[T]) Get() (*T, error) {
	return &tx.data, nil
}

func (tx *fileTx[T]) Commit() error {
	Save(&tx.data, tx.path)
	return nil
}

func (tx *fileTx[T]) Rollback() error {
	// Rollback is a no-op for file-based transactions
	return nil
}

func (store *FileStore[T]) WithTx(fn func(t *T) error) error {
	tx, err := store.Begin()
	if err != nil {
		return err
	}

	data, err := tx.Get()
	if err != nil {
		return err
	}

	if err := fn(data); err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			log.Printf("Rollback failed: %v", rollbackErr)
		}
		return err
	}

	return tx.Commit()
}

func CreateManager() *core.ContextManager {
	return core.NewContextManager(ctxtime.NewTimer(),
		&FileStore[core.State]{path: filepath.Join(viper.GetString("storePath"), "state")},
	)
}

func Load[T any](path string) T {
	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatal("Unable to read state file")
	}

	obj := new(T)
	err = json.Unmarshal(data, &obj)
	if err != nil {
		log.Fatal("Unable to parse state file")
	}

	return *obj
}

func Save[T any](obj *T, path string) {
	data, err := json.Marshal(obj)
	if err != nil {
		panic(err)
	}
	os.WriteFile(path, data, 0777)
}

func LoadState() core.State {
	statePath := filepath.Join(viper.GetString("storePath"), "state")
	data, err := os.ReadFile(statePath)
	if err != nil {
		log.Fatal("Unable to read state file")
	}

	state := core.State{}
	err = json.Unmarshal(data, &state)
	if err != nil {
		log.Fatal("Unable to parse state file")
	}

	return state
}

func SaveState(state *core.State) {
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

func (store *LocalArchiveStore) saveArchive(entry *core.ContextArchive, path string) error {
	data, err := json.Marshal(entry)
	if err != nil {
		return errors.New("unable to marshal archive for " + entry.Context.Id)
	}

	os.WriteFile(path, data, 0777)

	return nil
}

func (store *LocalArchiveStore) saveEventsArchive(entry *core.EventsArchive, path string) error {
	data, err := json.Marshal(entry)
	if err != nil {
		return errors.New("unable to marshal events archive for " + path)
	}

	os.WriteFile(path, data, 0777)

	return nil
}

func (store *LocalArchiveStore) loadArchive(id string, path string) (*core.ContextArchive, error) {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return &core.ContextArchive{
				Context: core.Context{
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

	entry := core.ContextArchive{}
	err = json.Unmarshal(data, &entry)

	if err != nil {
		return nil, errors.New("unable to parse archive file " + path)
	}

	return &entry, nil

}

func (store *LocalArchiveStore) loadEventsArchive(path string) (*core.EventsArchive, error) {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return &core.EventsArchive{
				Events: []core.Event{},
			}, nil
		} else {
			return nil, errors.New("unable to read eventsarchive file " + path)
		}
	}
	data, err := os.ReadFile(path)

	if err != nil {
		return nil, errors.New("unable to read events archive file " + path)
	}

	entry := core.EventsArchive{}
	err = json.Unmarshal(data, &entry)

	if err != nil {
		return nil, errors.New("unable to parse events archive file " + path)
	}

	return &entry, nil

}

func (store *LocalArchiveStore) Apply(id string, fn core.ArchivePatch) error {
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

func (store *LocalArchiveStore) ApplyEvents(date string, fn core.ArchiveEventsPatch) error {
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

func (store *LocalContextStore) Apply(fn core.StatePatch) error {
	state := LoadState()
	err := fn(&state)
	if err != nil {
		return err
	} else {
		SaveState(&state)
		return nil
	}
}

func (store *LocalContextStore) Read(fn core.StatePatch) error {
	state := LoadState()
	return fn(&state)
}

func LoadEvents() core.EventRegistry {
	eventsPath := filepath.Join(viper.GetString("storePath"), "events")
	data, err := os.ReadFile(eventsPath)
	if err != nil {
		log.Fatal("Unable to read state file")
	}

	events := core.EventRegistry{}
	err = json.Unmarshal(data, &events)
	if err != nil {
		log.Fatal("Unable to parse state file")
	}

	return events
}

func SaveEvents(eventsRegistry *core.EventRegistry) {
	eventsPath := filepath.Join(viper.GetString("storePath"), "events")
	data, err := json.Marshal(eventsRegistry)
	if err != nil {
		panic(err)
	}
	os.WriteFile(eventsPath, data, 0777)
}

func (store *LocalEventsStore) Apply(fn core.EventsPatch) error {
	events := LoadEvents()
	err := fn(&events)
	if err != nil {
		return err
	} else {
		SaveEvents(&events)
		return nil
	}
}

func (store *LocalEventsStore) Read(fn core.EventsPatch) error {
	events := LoadEvents()
	return fn(&events)
}
