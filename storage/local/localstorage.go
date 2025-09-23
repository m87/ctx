package localstorage

import (
	"encoding/json"
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
	data, err := Load[T](store.path)

	if err != nil {
		return nil, err
	}

	return &fileTx[T]{
		path: store.path,
		data: data}, nil
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
		&LocalStoreContextArchiver{path: filepath.Join(viper.GetString("storePath"), "archive")},
	)
}

func Load[T any](path string) (T, error) {
	var none T
	l, err := core.LockWithTimeout()
	if err != nil {
		panic(err)
	}
	defer l.Unlock()

	data, err := os.ReadFile(path)
	if err != nil {
		return none, err
	}

	var obj T
	err = json.Unmarshal(data, &obj)
	if err != nil {
		return none, err
	}

	return obj, err
}

func Save[T any](obj *T, path string) error {
	l, err := core.LockWithTimeout()
	if err != nil {
		return err
	}
	defer l.Unlock()

	data, err := json.Marshal(obj)
	if err != nil {
		return err
	}
	os.WriteFile(path, data, 0777)
	return nil
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
		log.Fatal("Unable to parse state file ", err)
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

func NewContextStore(path string) *LocalContextStore {
	return &LocalContextStore{
		path: path,
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
