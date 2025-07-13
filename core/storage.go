package core

type StatePatch func(*State) error

type EventsPatch func(*EventRegistry) error

type ArchivePatch func(*ContextArchive) error

type ArchiveEventsPatch func(*EventsArchive) error

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
