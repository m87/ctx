package core

import (
	"strings"

	ctxtime "github.com/m87/ctx/time"
)

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
	DateTime    ctxtime.ZonedTime `json:"dateTime"`
	CtxId       string            `json:"ctxId"`
	Description string            `json:"description"`
	Data        map[string]string `json:"data"`
	Type        EventType         `json:"type"`
}

type EventRegistry struct {
	Events []Event `json:"events"`
}
