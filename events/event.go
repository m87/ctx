package events

import "time"
import "os"
import "github.com/m87/ctx/util"
import "github.com/spf13/viper"
import "path/filepath"
import "encoding/json"

type EventType int

const (
	CREATE_CTX EventType = iota
	SWITCH_CTX
)

type Event struct {
	DateTime time.Time
	Data     map[string]string
	Type     EventType
}

type EventRegistry struct {
	Events []Event
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
