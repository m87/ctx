package core

func (session *Session) DeleteInterval(ctxId string, id string) error {
	if err := session.IsValidContext(ctxId); err != nil {
		return err
	}

	intervalIndex, err := session.ValidateIntervalExistsAndGet(ctxId, id)
	if err != nil {
		return err
	}

	ctx := session.State.Contexts[ctxId]
	ctx.Duration -= ctx.Intervals[intervalIndex].Duration
	ctx.Intervals = append(ctx.Intervals[:intervalIndex], ctx.Intervals[intervalIndex+1:]...)
	session.State.Contexts[ctxId] = ctx

	return nil
}
