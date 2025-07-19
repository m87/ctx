package core

import "errors"

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
	if session.State.CurrentId != "" {
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

func (session *Session) ValidateIntervalExists(ctxId, id string) error {
	ctx := session.State.Contexts[ctxId]
	if _, ok := ctx.Intervals[id]; ok {
		return nil
	}
	return errors.New("interval does not exist")
}
