package events_model

import (
	"time"
)

type EventType int

const (
	CREATE_CTX EventType = iota
	SWITCH_CTX
)

type Event struct {
	UUID        string            `json:"uuid"`
	DateTime    time.Time         `json:"dateTime"`
	CtxId       string            `json:"subject"`
	Description string            `json:"description"`
	Data        map[string]string `json:"data"`
	Type        EventType         `json:"type"`
}

type EventRegistry struct {
	Events []Event `json:"events"`
}
