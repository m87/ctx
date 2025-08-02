package core

import "sync"

var Mutex sync.Mutex

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
