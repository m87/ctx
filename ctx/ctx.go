package ctx

import (
	"encoding/json"
	"errors"
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
	Start    time.Time     `json:"start"`
	End      time.Time     `json:"end"`
	Duration time.Duration `json:"duration"`
}

type Context struct {
	Id          string        `json:"id"`
	Description string        `json:"description"`
	State       ContextState  `json:"state"`
	Duration    time.Duration `json:"duration"`
	Intervals   []Interval    `json:"intervals"`
}

type State struct {
	Contexts  map[string]Context `json:"contexts"`
	CurrentId string             `json:"currentId"`
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

func Pause(state *State) {
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

func Switch(id string, state *State, eventsRegistry *events.EventRegistry) error {
	if state.CurrentId == id {
		return nil
	}
	now := time.Now().Local()
	prevId := state.CurrentId
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
			DateTime:    now,
			Type:        events.SWITCH_CTX,
			CtxId:       ctx.Id,
			Description: ctx.Description,
			Data: map[string]string{
				"from": prevId,
			},
		}, eventsRegistry)
		ctx.Intervals = append(state.Contexts[id].Intervals, Interval{Start: now})
		state.Contexts[id] = ctx
		return nil
	} else {
		return errors.New("not found")
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

func Stop(id string, state *State) {
	now := time.Now().Local()
	if state.CurrentId == id {
		prev := state.Contexts[state.CurrentId]
		interval := prev.Intervals[len(prev.Intervals)-1]
		interval.End = now
		interval.Duration = interval.End.Sub(interval.Start)
		state.Contexts[state.CurrentId].Intervals[len(prev.Intervals)-1] = interval
		prev.Duration = prev.Duration + interval.Duration
		state.Contexts[state.CurrentId] = prev
		state.CurrentId = ""
	}

	//TODO create contexts history move to history

}

func Rename(id string, newDescription string, state *State) {
	newId := util.GenerateId(newDescription)
	ctx := state.Contexts[id]
	ctx.Id = newId
	ctx.Description = newDescription
	state.Contexts[newId] = ctx
	delete(state.Contexts, id)
	if state.CurrentId == id {
		state.CurrentId = newId
	}

}

func Delete(id string, state *State) {
	if state.CurrentId == id {
		state.CurrentId = ""
	}

	delete(state.Contexts, id)
}
