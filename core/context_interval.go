package core

func (session *Session) DeleteInterval(ctxId string, id string) error {
	if err := session.IsValidContext(ctxId); err != nil {
		return err
	}

	err := session.ValidateIntervalExists(ctxId, id)
	if err != nil {
		return err
	}

	ctx := session.State.Contexts[ctxId]
	ctx.Duration -= ctx.Intervals[id].Duration
	delete(ctx.Intervals, id)
	session.State.Contexts[ctxId] = ctx

	return nil
}
