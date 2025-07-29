package localstorage

import (
	"encoding/json"
	"log"
	"os"

	"github.com/m87/ctx/core"
)

type multifileTx[T any] struct {
	path string
	data T
}

type MultiFileStore[T any] struct {
	path string
}

func (store *MultiFileStore[T]) Begin() (core.Tx[T], error) {
	return &multifileTx[T]{
		path: store.path,
		data: LoadMultifile[T](store.path),
	}, nil
}

func (store *MultiFileStore[T]) BeginAndGet() (core.Tx[T], *T, error) {
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

func (tx *multifileTx[T]) Get() (*T, error) {
	return &tx.data, nil
}

func (tx *multifileTx[T]) Commit() error {
	SaveMultifile(&tx.data, tx.path)
	return nil
}

func (tx *multifileTx[T]) Rollback() error {
	// Rollback is a no-op for file-based transactions
	return nil
}

func (store *MultiFileStore[T]) WithTx(fn func(t *T) error) error {
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

func LoadMultifile[T any](path string) T {
	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatal("Unable to read state file")
	}

	obj := new(T)
	err = json.Unmarshal(data, &obj)
	if err != nil {
		log.Fatal("Unable to parse state file", err)
	}

	return *obj
}

func SaveMultifile[T any](obj *T, path string) {
	data, err := json.Marshal(obj)
	if err != nil {
		panic(err)
	}
	os.WriteFile(path, data, 0777)
}
