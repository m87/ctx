package ctx

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/m87/ctx/events"
	"github.com/m87/ctx/util"
	"github.com/spf13/viper"
)

type ContextState int

const (
	ACTIVE ContextState = iota
	FINISHED
)

type Interval struct {
	Start    time.Time
	End      time.Time
	Duration time.Duration
}

type Context struct {
	Id          string
	Description string
	State       ContextState
	Duration    time.Duration
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

func Stop(state *State) {
	now := time.Now().Local()
	if state.CurrentId != "" {
		prev := state.Contexts[state.CurrentId]
		interval := prev.Intervals[len(prev.Intervals)-1]
		interval.End = now
		interval.Duration = interval.End.Sub(interval.Start)
		state.Contexts[state.CurrentId].Intervals[len(prev.Intervals)-1] = interval
		prev.Duration = prev.Duration + interval.Duration
		state.Contexts[state.CurrentId] = prev
		state.CurrentId = ""
	}
}

func Switch(id string, state *State, eventsRegistry *events.EventRegistry) {
	if state.CurrentId == id {
		return
	}
	now := time.Now().Local()
	if state.CurrentId != "" {
		prev := state.Contexts[state.CurrentId]
		interval := prev.Intervals[len(prev.Intervals)-1]
		interval.End = now
		interval.Duration = interval.End.Sub(interval.Start)
		state.Contexts[state.CurrentId].Intervals[len(prev.Intervals)-1] = interval
		prev.Duration = prev.Duration + interval.Duration
		state.Contexts[state.CurrentId] = prev
	}

	if ctx, ok := state.Contexts[id]; ok {
		state.CurrentId = ctx.Id
		events.Publish(events.Event{
			DateTime: now,
			Type:     events.SWITCH_CTX,
			Data: map[string]string{
				"Description": ctx.Description,
			},
		}, eventsRegistry)
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
