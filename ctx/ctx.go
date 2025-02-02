package ctx

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	"github.com/m87/ctx/util"
	"github.com/spf13/viper"
)

type ContextState int

const (
	ACTIVE ContextState = iota
	FINISHED
)

type Interval struct {
	start    time.Time
	end      time.Time
	duration int32
}

type Context struct {
	Id          string
	Description string
	State       ContextState
	Duration    int32
	Intervals   []Interval
}

type State struct {
	Contexts []Context
	Current  Context
}

func Load() State {
	statePath := filepath.Join(viper.GetString("ctxPath"), "state")
	data, err := os.ReadFile(statePath)
	util.Check(err, "Unable to read state file")

	state := State{}
	err = json.Unmarshal(data, &state)
	util.Check(err, "Unable to parse state file")

	return state
}

func Save(state State) {
	statePath := filepath.Join(viper.GetString("ctxPath"), "state")
	data, err := json.Marshal(state)
	if err != nil {
		panic(err)
	}
	os.WriteFile(statePath, data, 0644)
}
