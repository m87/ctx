package storage

import "github.com/m87/ctx/ctx_model"

type StatePatch func(*ctx_model.State) error

type EventsPatch func(*ctx_model.EventRegistry) error

type ArchivePatch func(*ctx_model.ContextArchive) error

type ArchiveEventsPatch func(*ctx_model.EventsArchive) error

type ContextStore interface {
	Apply(fn StatePatch) error
	Read(fn StatePatch) error
}

type EventsStore interface {
	Apply(fn EventsPatch) error
	Read(fn EventsPatch) error
}

type ArchiveStore interface {
	Apply(id string, fn ArchivePatch) error
	ApplyEvents(date string, fn ArchiveEventsPatch) error
}
