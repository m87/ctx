package core

import (
	"sort"
	"time"
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
