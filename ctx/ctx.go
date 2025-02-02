package ctx

import (
	"encoding/json"
	"log"
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
	start    int64
	end      int64
	duration int64
}

type Context struct {
	Id          string
	Description string
	State       ContextState
	Duration    int64
	Intervals   []Interval
}

type State struct {
	Contexts  []Context
	CurrentId string
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

func Switch(id string, state *State) {
	for _, v := range state.Contexts {
		if id == v.Id {
			prev := state.Current
			now := time.Now().Local().UnixMilli()
			if state.Current.Id != "" {
				interval := prev.Intervals[len(prev.Intervals)-1]
				interval.end = now
				interval.duration = interval.end - interval.start
				prev.Duration = prev.Duration + interval.duration
			}
			state.Current = v
			state.Current.Intervals = append(state.Current.Intervals, Interval{start: now})

			return
		}
	}
	log.Printf("context: %s, not found\n", id)
}

func Save(state *State) {
	statePath := filepath.Join(viper.GetString("ctxPath"), "state")
	data, err := json.Marshal(state)
	if err != nil {
		panic(err)
	}
	os.WriteFile(statePath, data, 0644)
}
