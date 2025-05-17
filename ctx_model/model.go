package ctx_model

import (
	"fmt"
	"strings"
	"time"
)

type ContextState int

const (
	ACTIVE ContextState = iota
	FINISHED
)

type Interval struct {
	Start    LocalTime     `json:"start"`
	End      LocalTime     `json:"end"`
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
	MERGE_CTX
	DELETE_CTX
	EDIT_CTX_INTERVAL
	RENAME_CTX
	DELETE_CTX_INTERVAL
)

func EventAsString(event EventType) string {
	switch event {
	case CREATE_CTX:
		return "CREATE"
	case SWITCH_CTX:
		return "SWITCH"
	case START_INTERVAL:
		return "START_INTERVAL"
	case END_INTERVAL:
		return "END_INTERVAL"
	case MERGE_CTX:
		return "MERGE_CTX"
	case DELETE_CTX:
		return "DELETE_CTX"
	case EDIT_CTX_INTERVAL:
		return "EDIT_CTX_INTERVAL"
	case RENAME_CTX:
		return "RENAME_CTX"
	case DELETE_CTX_INTERVAL:
		return "DELETE_CTX_INTERVAL"
	}
	panic("undefined event type")
}

func StringAsEvent(event string) EventType {
	switch strings.ToUpper(event) {
	case "CREATE":
		return CREATE_CTX
	case "SWITCH":
		return SWITCH_CTX
	case "START_INTERVAL":
		return START_INTERVAL
	case "END_INTERVAL":
		return END_INTERVAL
	case "MERGE_CTX":
		return MERGE_CTX
	case "DELETE_CTX":
		return DELETE_CTX
	case "EDIT_CTX_INTERVAL":
		return EDIT_CTX_INTERVAL
	case "RENAME_CTX":
		return RENAME_CTX
	case "DELETE_CTX_INTERVAL":
		return DELETE_CTX_INTERVAL
	}
	panic("undefined event type")
}

type Event struct {
	UUID        string            `json:"uuid"`
	DateTime    LocalTime         `json:"dateTime"`
	CtxId       string            `json:"ctxId"`
	Description string            `json:"description"`
	Data        map[string]string `json:"data"`
	Type        EventType         `json:"type"`
}

type EventRegistry struct {
	Events []Event `json:"events"`
}

type ContextArchive struct {
	Context Context `json:"context"`
}

type EventsArchive struct {
	Date   string  `json:"date"`
	Events []Event `json:"events"`
}

type EventsFilter struct {
	Date  string
	Types []string
	CtxId string
}

type State struct {
	Contexts  map[string]Context `json:"contexts"`
	CurrentId string             `json:"currentId"`
}

type LocalTime struct {
	time.Time
}

func (lt LocalTime) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%s"`, lt.Time.Format(time.RFC3339Nano))), nil
}

func (lt *LocalTime) UnmarshalJSON(data []byte) error {
	s := strings.Trim(string(data), `"`)
	t, err := time.ParseInLocation(time.RFC3339Nano, s, time.Local)
	if err != nil {
		return err
	}
	lt.Time = t
	return nil
}

type StatePatch func(*State) error

type EventsPatch func(*EventRegistry) error

type ArchivePatch func(*ContextArchive) error

type ArchiveEventsPatch func(*EventsArchive) error

type TimeProvider interface {
	Now() LocalTime
}

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
