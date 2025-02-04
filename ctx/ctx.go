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
	Start    int64
	End      int64
	Duration int64
}

type Context struct {
	Id          string
	Description string
	State       ContextState
	Duration    int64
	Intervals   []Interval
}

type State struct {
	Contexts  map[string]Context
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
	now := time.Now().Local().UnixMilli()
	if state.CurrentId != "" {
		prev := state.Contexts[state.CurrentId]
		interval := prev.Intervals[len(prev.Intervals)-1]
		interval.End = now
		interval.Duration = interval.End - interval.Start
		prev.Duration = prev.Duration + interval.Duration
	}

	if ctx, ok := state.Contexts[id]; ok {
		state.CurrentId = ctx.Id
		ctx.Intervals = append(state.Contexts[id].Intervals, Interval{Start: now})
		state.Contexts[id] = ctx
	} else {
		log.Printf("context: %s, not found\n", id)
	}
}

func Save(state *State) {
	statePath := filepath.Join(viper.GetString("ctxPath"), "state")
	data, err := json.Marshal(state)
	if err != nil {
		panic(err)
	}
	os.WriteFile(statePath, data, 0644)
}
