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

func (session *Session) ValidateContextExists(id string) error {
	if _, ok := session.State.Contexts[id]; !ok {
		return errors.New("context does not exist")
	}
	return nil
}

func (session *Session) ValidateIntervalExistsAndGet(ctxId, id string) (int, error) {
	ctx := session.State.Contexts[ctxId]
	for i, interval := range ctx.Intervals {
		if interval.Id == id {
			return i, nil
		}
	}
	return -1, errors.New("interval does not exist")
}
