package core

import (
	"errors"

	ctxtime "github.com/m87/ctx/time"
)

type ContextManager struct {
	TimeProvider     ctxtime.TimeProvider
	StateStore       TransactionalStore[State]
	ContextArchiver  Archiver[Context]
	MigrationManager MigrationManager
}

type Session struct {
	State        *State
	TimeProvider ctxtime.TimeProvider
}

func NewContextManager(timeProvider ctxtime.TimeProvider, stateStore TransactionalStore[State], contextArchiver Archiver[Context], migrationManager MigrationManager) *ContextManager {
	return &ContextManager{
		TimeProvider:     timeProvider,
		StateStore:       stateStore,
		ContextArchiver:  contextArchiver,
		MigrationManager: migrationManager,
	}
}

func (manager *ContextManager) WithArchiveSession(fn func(session Session) error) error {
	archiveTx, archive, archiveErr := manager.StateStore.BeginAndGet()
	if archiveErr != nil {
		return errors.Join(archiveErr)
	}

	if err := fn(Session{
		State:        archive,
		TimeProvider: manager.TimeProvider,
	}); err != nil {
		return err
	}

	archiveErr = archiveTx.Commit()

	if archiveErr != nil {
		archiveRollbackErr := archiveTx.Rollback()

		if archiveRollbackErr != nil {
			panic(errors.Join(archiveRollbackErr))
		}
	}

	return nil
}

func (manager *ContextManager) WithContextArchiver(fn func(archver Archiver[Context]) error) error {
	return fn(manager.ContextArchiver)
}

func (manager *ContextManager) WithSession(fn func(session Session) error) error {
	stateTx, state, stateErr := manager.StateStore.BeginAndGet()
	if stateErr != nil {
		return errors.Join(stateErr)
	}

	if err := fn(Session{
		State:        state,
		TimeProvider: manager.TimeProvider,
	}); err != nil {
		return err
	}

	stateErr = stateTx.Commit()

	if stateErr != nil {
		stateRollbackEer := stateTx.Rollback()

		if stateRollbackEer != nil {
			panic(errors.Join(stateRollbackEer))
		}
	}

	return nil
}
