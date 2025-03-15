package ctx_model

import "time"

type ContextState int

const (
	ACTIVE ContextState = iota
	FINISHED
)

type Interval struct {
	Start    time.Time     `json:"start"`
	End      time.Time     `json:"end"`
	Duration time.Duration `json:"duration"`
}

type Context struct {
	Id          string        `json:"id"`
	Description string        `json:"description"`
	Comments    []string      `json:"comments"`
	State       ContextState  `json:"state"`
	Duration    time.Duration `json:"duration"`
	Intervals   []Interval    `json:"intervals"`
}

type State struct {
	Contexts  map[string]Context `json:"contexts"`
	CurrentId string             `json:"currentId"`
}

type StatePatch func(*State)

type TimeProvider interface {
	Now() time.Time
}

type ContextStore interface {
	Apply(fn StatePatch)
}

type EventsRegistryStore interface {
}

type ArchiveStore interface {
}

type ContextManager struct {
	ContextStore ContextStore
	TimeProvider TimeProvider
}
