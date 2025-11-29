package core

import (
	"errors"
	"regexp"
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
	Comments    map[string]Comment  `json:"comments"`
	State       ContextState        `json:"state"`
	Duration    time.Duration       `json:"duration"`
	Intervals   map[string]Interval `json:"intervals"`
	Labels      []string            `json:"labels"`
}

type Comment struct {
	Id      string `json:"id"`
	Content string `json:"content"`
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
	if err := session.IsValidContext(srcId); err != nil {
		return err
	}

	if err := session.ValidateContextExists(targetId); err == nil {
		return errors.New("target context already exists")
	}

	ctx := session.State.Contexts[srcId]
	ctx.Description = description
	ctx.Id = targetId
	session.State.Contexts[targetId] = ctx
	delete(session.State.Contexts, srcId)
	return nil

}

func (session *Session) Free() error {
	if err := session.ValidateAnyActiveContext(); err != nil {
		return err
	}
	now := session.TimeProvider.Now()
	session.endInterval(session.State.CurrentId, now)
	session.State.CurrentId = ""
	return nil
}

func (session *Session) MustGetCtx(ctxId string) Context {
	return session.State.Contexts[ctxId]
}

func (session *Session) GetActiveCtx() (Context, error) {
	ctxId := session.State.CurrentId
	if err := session.ValidateAnyActiveContext(); err != nil {
		return Context{}, err
	}

	return session.State.Contexts[ctxId], nil
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

	delete(session.State.Contexts, ctxId)
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

	for _, comment := range fromCtx.Comments {
		toCtx.Comments[comment.Id] = comment
	}

	toCtx.Labels = append(toCtx.Labels, fromCtx.Labels...)
	toCtx.Duration = toCtx.Duration + fromCtx.Duration

	for _, interval := range fromCtx.Intervals {
		if _, ok := toCtx.Intervals[interval.Id]; !ok {
			toCtx.Intervals[interval.Id] = interval
		}
	}

	session.SetCtx(toCtx)
	session.deleteInternal(from)

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
		Comments:    map[string]Comment{},
	}
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
	if state.CurrentId != "" {
		session.endInterval(state.CurrentId, now)
	}

	if ctx, ok := state.Contexts[ctxId]; ok {
		state.CurrentId = ctx.Id
		intervalId := uuid.NewString()
		ctx.Intervals[intervalId] = Interval{Id: intervalId, Start: now}
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

func (session *Session) Search(regex string) ([]Context, error) {
	ctxs := []Context{}
	re := regexp.MustCompile(regex)
	for _, ctx := range session.State.Contexts {
		if re.MatchString(ctx.Description) {
			ctxs = append(ctxs, ctx)
		}
	}
	return ctxs, nil
}

func (session *Session) GetContextCountByDateMap() map[string]int {
	counts := make(map[string]int)
	for _, ctx := range session.State.Contexts {
		for k, v := range session.GetDateCounts(ctx.Id) {
			counts[k] = counts[k] + v
		}
	}
	return counts
}

func (session *Session) IsContextActive(ctxId string) bool {
	return session.State.CurrentId == ctxId
}
