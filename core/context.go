package core

import (
	"errors"
	"sort"
	"time"

	"github.com/google/uuid"
)

type ContextState int

const (
	ACTIVE ContextState = iota
	FINISHED
)

type Context struct {
	Id          string              `json:"id"`
	Description string              `json:"description"`
	Comments    []string            `json:"comments"`
	State       ContextState        `json:"state"`
	Duration    time.Duration       `json:"duration"`
	Intervals   map[string]Interval `json:"intervals"`
	Labels      []string            `json:"labels"`
}

type ContextArchive struct {
	Context Context `json:"context"`
}

type EventsArchive struct {
	Date   string  `json:"date"`
	Events []Event `json:"events"`
}

type EventsFilter struct {
	Date  string
	Types []string
	CtxId string
}

type State struct {
	Contexts  map[string]Context `json:"contexts"`
	CurrentId string             `json:"currentId"`
}

func (session *Session) GetSortedContextIds() []string {
	ids := []string{}
	for k := range session.State.Contexts {
		ids = append(ids, k)
	}
	sort.Strings(ids)
	return ids
}

func (session *Session) RenameContext(srcId string, targetId string, description string) error {
	ctx := session.State.Contexts[srcId]
	ctx.Description = description
	session.State.Contexts[targetId] = ctx
	delete(session.State.Contexts, srcId)
	//		manager.PublishContextEvent(ctx, manager.TimeProvider.Now(), RENAME_CTX, map[string]string{
	//			"src.id":             ctx.Id,
	//			"src.description":    ctx.Description,
	//			"target.id":          targetId,
	//			"target:description": name,
	//		})
	return nil

}

func (session *Session) Free() error {
	if err := session.ValidateAnyActiveContext(); err != nil {
		return err
	}
	//now := session.TimeProvider.Now()
	//manager.endInterval(session.State, session.State.CurrentId, now)
	session.State.CurrentId = ""
	return nil
}

func (session *Session) MustGetCtx(ctxId string) Context {
	return session.State.Contexts[ctxId]
}

func (session *Session) GetCtx(ctxId string) (Context, error) {
	if err := session.IsValidContext(ctxId); err != nil {
		return Context{}, err
	}

	return session.State.Contexts[ctxId], nil
}

func (session *Session) SetCtx(ctx Context) {
	session.State.Contexts[ctx.Id] = ctx
}

func (session *Session) deleteInternal(ctxId string) error {
	if err := session.IsValidContext(ctxId); err != nil {
		return err
	}

	// ctx := session.State.Contexts[ctxId]
	delete(session.State.Contexts, ctxId)
	// manager.PublishContextEvent(context, manager.TimeProvider.Now(), DELETE_CTX, nil)
	return nil
}

func (session *Session) Delete(ctxId string) error {
	return session.deleteInternal(ctxId)
}

func (session *Session) MergeContext(from string, to string) error {
	if from == to {
		return errors.New("contexts are the same")
	}

	if err := session.ValidateContextsExist(from, to); err != nil {
		return err
	}

	if err := session.ValidateActiveContext(from); err != nil {
		return err
	}

	fromCtx := session.MustGetCtx(from)
	toCtx := session.MustGetCtx(to)

	toCtx.Comments = append(toCtx.Comments, fromCtx.Comments...)
	toCtx.Labels = append(toCtx.Labels, fromCtx.Labels...)
	toCtx.Duration = toCtx.Duration + fromCtx.Duration

	for _, interval := range fromCtx.Intervals {
		if _, ok := toCtx.Intervals[interval.Id]; !ok {
			toCtx.Intervals[interval.Id] = interval
		}
	}

	session.SetCtx(toCtx)
	session.deleteInternal(from)

	// manager.PublishContextEvent(state.Contexts[to], manager.TimeProvider.Now(), MERGE_CTX, map[string]string{
	// 	"from": from,
	// 	"to":   to,
	// })

	return nil
}

func (session *Session) createContetxtInternal(id string, description string) error {
	if err := IsValidDescription(description); err != nil {
		return err
	}

	if err := session.ValidateContextAlreadyExists(id); err != nil {
		return err
	}

	session.State.Contexts[id] = Context{
		Id:          id,
		Description: description,
		State:       ACTIVE,
		Intervals:   map[string]Interval{},
		Labels:      []string{},
		Comments:    []string{},
	}
	//manager.PublishContextEvent(state.Contexts[id], manager.TimeProvider.Now(), CREATE_CTX, nil)
	return nil
}

func (session *Session) CreateContext(ctxId string, description string) error {
	return session.createContetxtInternal(ctxId, description)
}

func (session *Session) switchInternal(ctxId string) error {
	if err := session.IsValidContext(ctxId); err != nil {
		return nil
	}

	state := session.State
	now := session.TimeProvider.Now()
	//prevId := state.CurrentId
	if state.CurrentId != "" {
		session.endInterval(state.CurrentId, now)
	}

	if ctx, ok := state.Contexts[ctxId]; ok {
		state.CurrentId = ctx.Id
		//manager.PublishContextEvent(state.Contexts[id], now, SWITCH_CTX, map[string]string{
		//	"from": prevId,
		//})
		intervalId := uuid.NewString()
		ctx.Intervals[intervalId] = Interval{Id: uuid.NewString(), Start: now}
		//manager.PublishContextEvent(state.Contexts[ctxId], now, START_INTERVAL, nil)
		state.Contexts[ctxId] = ctx
	}
	return nil
}

func (session *Session) Switch(ctxId string) error {
	if err := session.IsValidContext(ctxId); err != nil {
		return err
	}
	return session.switchInternal(ctxId)
}

func (session *Session) contextExists(ctxId string) bool {
	_, ok := session.State.Contexts[ctxId]
	return ok
}

func (session *Session) contextNotExists(ctxId string) bool {
	return !session.contextExists(ctxId)
}

func (session *Session) CreateIfNotExistsAndSwitch(ctxId string, description string) error {
	if session.contextNotExists(ctxId) {
		err := session.createContetxtInternal(ctxId, description)
		if err != nil {
			return err
		}
	}
	return session.switchInternal(ctxId)
}
