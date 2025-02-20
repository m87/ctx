package ctx

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/m87/ctx/ctx_model"
	"github.com/m87/ctx/events"
	"github.com/m87/ctx/util"
)

func Pause(state *ctx_model.State) {
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

func Switch(id string, state *ctx_model.State, eventsRegistry *events.EventRegistry) error {
	if state.CurrentId == id {
		return errors.New("already active")
	}

	if _, ok := state.Contexts[id]; !ok {
		return errors.New("not found")

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
			UUID:        uuid.NewString(),
			DateTime:    now,
			Type:        events.SWITCH_CTX,
			CtxId:       ctx.Id,
			Description: ctx.Description,
			Data: map[string]string{
				"from": prevId,
			},
		}, eventsRegistry)
		ctx.Intervals = append(state.Contexts[id].Intervals, ctx_model.Interval{Start: now})
		state.Contexts[id] = ctx
		return nil
	} else {
		return errors.New("not found")
	}
}

func Comment(id string, comment string, state *ctx_model.State) {
	ctx := state.Contexts[id]
	ctx.Comments = append(ctx.Comments, comment)
	state.Contexts[id] = ctx
}

func Stop(id string, state *ctx_model.State) {
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

func Rename(id string, newDescription string, state *ctx_model.State) {
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

func Delete(id string, state *ctx_model.State) {
	if state.CurrentId == id {
		state.CurrentId = ""
	}

	delete(state.Contexts, id)
}
