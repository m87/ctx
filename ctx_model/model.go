package ctx_model

import (
	"encoding/json"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

type ContextState int

const (
	ACTIVE ContextState = iota
	FINISHED
)

type Interval struct {
	Id       string        `json:"id"`
	Start    ZonedTime     `json:"start"`
	End      ZonedTime     `json:"end"`
	Duration time.Duration `json:"duration"`
	Labels   []string      `json:"labels"`
}

type Context struct {
	Id          string        `json:"id"`
	Description string        `json:"description"`
	Comments    []string      `json:"comments"`
	State       ContextState  `json:"state"`
	Duration    time.Duration `json:"duration"`
	Intervals   []Interval    `json:"intervals"`
	Labels      []string      `json:"labels"`
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
	LABEL_CTX
	DELETE_CTX_LABEL
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
	case LABEL_CTX:
		return "LABEL_CTX"
	case DELETE_CTX_LABEL:
		return "DELETE_CTX_LABEL"
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
	case "LABEL_CTX":
		return LABEL_CTX
	case "DELETE_CTX_LABEL":
		return DELETE_CTX_LABEL
	}
	panic("undefined event type")
}

type Event struct {
	UUID        string            `json:"uuid"`
	DateTime    ZonedTime         `json:"dateTime"`
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

type ZonedTime struct {
	Time     time.Time `json:"time"`
	Timezone string    `json:"timezone"`
}

func DetectTimezoneName() string {
	switch runtime.GOOS {
	case "linux", "darwin":
		return detectUnixTimezone()
	default:
		return "UTC"
	}
}

func detectUnixTimezone() string {
	out, err := exec.Command("readlink", "-f", "/etc/localtime").Output()
	if err != nil {
		return "UTC"
	}

	path := strings.TrimSpace(string(out))
	const zoneinfoPrefix = "/usr/share/zoneinfo/"
	if strings.HasPrefix(path, zoneinfoPrefix) {
		return strings.TrimPrefix(path, zoneinfoPrefix)
	}

	return "UTC"
}

func (zt ZonedTime) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Time     string `json:"time"`
		Timezone string `json:"timezone"`
	}{
		Time:     zt.Time.Format(time.RFC3339),
		Timezone: zt.Time.Location().String(),
	})
}

func (zt *ZonedTime) UnmarshalJSON(data []byte) error {
	var tmp struct {
		Time     string `json:"time"`
		Timezone string `json:"timezone"`
	}
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	loc, err := time.LoadLocation(tmp.Timezone)
	if err != nil {
		return err
	}
	t, err := time.ParseInLocation(time.RFC3339, tmp.Time, loc)
	if err != nil {
		return err
	}
	zt.Time = t
	zt.Timezone = tmp.Timezone
	return nil
}

type StatePatch func(*State) error

type EventsPatch func(*EventRegistry) error

type ArchivePatch func(*ContextArchive) error

type ArchiveEventsPatch func(*EventsArchive) error

type TimeProvider interface {
	Now() ZonedTime
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
