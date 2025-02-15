package ctx_store

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"

	"github.com/m87/ctx/ctx_model"
	"github.com/spf13/viper"
)

func Load() ctx_model.State {
	statePath := filepath.Join(viper.GetString("ctxPath"), "state")
	data, err := os.ReadFile(statePath)
	if err != nil {
		log.Fatal("Unable to read state file")
	}

	state := ctx_model.State{}
	err = json.Unmarshal(data, &state)
	if err != nil {
		log.Fatal("Unable to parse state file")
	}

	return state
}

func Save(state *ctx_model.State) {
	statePath := filepath.Join(viper.GetString("ctxPath"), "state")
	data, err := json.Marshal(state)
	if err != nil {
		panic(err)
	}
	os.WriteFile(statePath, data, 0644)
}
