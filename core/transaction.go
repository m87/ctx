package core

import (
	"context"
	"path"
	"sync"
	"time"

	"github.com/gofrs/flock"
	"github.com/spf13/viper"
)

var Mutex sync.Mutex

func LockWithTimeout() (*flock.Flock, error){
	l := flock.New(path.Join(viper.GetString("storePath"), "ctx.lock"))
	ctx, cancel := context.WithTimeout(context.Background(), 30 * time.Second)
	defer cancel()


	if _, err := l.TryLockContext(ctx, 100*time.Millisecond); err != nil {
		return nil, err
	}

	return l, nil
}


type TransactionalStore[T any] interface {
	Begin() (Tx[T], error)
	BeginAndGet() (Tx[T], *T, error)
	WithTx(fn func(t *T) error) error
}

type Tx[T any] interface {
	Get() (*T, error)
	Commit() error
	Rollback() error
}

type TxEntry interface {
	Commit() error
	Rollback() error
}

type TransactionManager struct {
	txs []TxEntry
}

func NewTransactionManager() *TransactionManager {
	return &TransactionManager{txs: make([]TxEntry, 0)}
}

func (tm *TransactionManager) Add(tx TxEntry) {
	tm.txs = append(tm.txs, tx)
}

func (tm *TransactionManager) Commit() error {
	for _, tx := range tm.txs {
		if err := tx.Commit(); err != nil {
			return err
		}
	}
	tm.txs = nil
	return nil
}

func (tm *TransactionManager) Rollback() error {
	for i := len(tm.txs) - 1; i >= 0; i-- {
		if err := tm.txs[i].Rollback(); err != nil {
			return err
		}
	}
	tm.txs = nil
	return nil
}
