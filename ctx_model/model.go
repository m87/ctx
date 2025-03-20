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

type EventType int

const (
	CREATE_CTX EventType = iota
	SWITCH_CTX
	START_INTERVAL
	END_INTERVAL
)

type Event struct {
	UUID        string            `json:"uuid"`
	DateTime    time.Time         `json:"dateTime"`
	CtxId       string            `json:"ctxId"`
	Description string            `json:"description"`
	Data        map[string]string `json:"data"`
	Type        EventType         `json:"type"`
}

type EventRegistry struct {
	Events []Event `json:"events"`
}

type State struct {
	Contexts  map[string]Context `json:"contexts"`
	CurrentId string             `json:"currentId"`
}

type StatePatch func(*State) error

type EventsPatch func(*EventRegistry) error

type TimeProvider interface {
	Now() time.Time
}

type ContextStore interface {
	Apply(fn StatePatch) error
	Read(fn StatePatch)
}

type EventsStore interface {
	Apply(fn EventsPatch) error
	Read(fn EventsPatch)
}

type ArchiveStore interface {
}
