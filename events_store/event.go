package events_store

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"

	"github.com/m87/ctx/events_model"
	"github.com/spf13/viper"
)

func Load() events_model.EventRegistry {
	registryPath := filepath.Join(viper.GetString("ctxPath"), "events")
	data, err := os.ReadFile(registryPath)
	if err != nil {
		log.Fatal("Unable to read evnets file")
	}

	registry := events_model.EventRegistry{}
	err = json.Unmarshal(data, &registry)
	if err != nil {
		log.Fatal("Unable to parse events file")
	}

	return registry
}

func Save(registry *events_model.EventRegistry) {
	registryPath := filepath.Join(viper.GetString("ctxPath"), "events")
	data, err := json.Marshal(registry)
	if err != nil {
		panic(err)
	}
	os.WriteFile(registryPath, data, 0644)
}
