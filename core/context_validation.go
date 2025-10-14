package core

import (
	"errors"
	"strings"
)

func (session *Session) IsValidContext(id string) error {
	return errors.Join(
		session.ValidateActiveContext(id),
		session.ValidateContextExists(id),
	)
}

func (session *Session) ValidateActiveContext(id string) error {
	if session.State.CurrentId == id {
		return errors.New("context is active")
	}
	return nil
}

func (session *Session) ValidateAnyActiveContext() error {
	if session.State.CurrentId == "" {
		return errors.New("no active context")
	}
	return nil
}

func (session *Session) ValidateContextExists(id string) error {
	if _, ok := session.State.Contexts[id]; !ok {
		return errors.New("context does not exist")
	}
	return nil
}

func (session *Session) ValidateContextsExist(ids ...string) error {
	for _, id := range ids {
		if err := session.ValidateContextExists(id); err != nil {
			return err
		}
	}
	return nil
}

func (session *Session) ValidateContextAlreadyExists(id string) error {
	if _, ok := session.State.Contexts[id]; ok {
		return errors.New("context already exists")
	}
	return nil
}

func (session *Session) ValidateIntervalExists(ctxId, id string) error {
	ctx := session.State.Contexts[ctxId]
	if _, ok := ctx.Intervals[id]; ok {
		return nil
	}
	return errors.New("interval does not exist")
}

func (session *Session) ValidateActiveInterval(ctxId, id string) error {
	if err := session.ValidateIntervalExists(ctxId, id); err != nil {
		return err
	}
	interval := session.State.Contexts[ctxId].Intervals[id]
	if interval.IsActive() {
		return errors.New("interval is active")
	}

	return nil
}

func IsValidDescription(description string) error {
	if strings.TrimSpace(description) == "" {
		return errors.New("empty description")
	}
	return nil
}
