package events

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	"github.com/m87/ctx/util"
	"github.com/spf13/viper"
)

type EventType int

const (
	CREATE_CTX EventType = iota
	SWITCH_CTX
)

type Event struct {
	DateTime    time.Time         `json:"dateTime"`
	CtxId       string            `json:"subject"`
	Description string            `json:"description"`
	Data        map[string]string `json:"data"`
	Type        EventType         `json:"type"`
}

type EventRegistry struct {
	Events []Event `json:"events"`
}

func Load() EventRegistry {
	registryPath := filepath.Join(viper.GetString("ctxPath"), "events")
	data, err := os.ReadFile(registryPath)
	util.Check(err, "Unable to read events registry file")

	registry := EventRegistry{}
	err = json.Unmarshal(data, &registry)
	util.Check(err, "Unable to parse events registry file")

	return registry
}

func Save(registry *EventRegistry) {
	registryPath := filepath.Join(viper.GetString("ctxPath"), "events")
	data, err := json.Marshal(registry)
	if err != nil {
		panic(err)
	}
	os.WriteFile(registryPath, data, 0644)
}

func Publish(event Event, registry *EventRegistry) {
	registry.Events = append(registry.Events, event)
}
