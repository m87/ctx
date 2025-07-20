package core

import ctxtime "github.com/m87/ctx/time"

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

func (session *Session) GetActiveIntervals(ctxId string) ([]string, error) {
  intervals := []string{}
	if err := session.ValidateContextExists(ctxId); err != nil {
		return intervals, err
	}
	if err := session.ValidateActiveContext(ctxId); err != nil {
    return intervals, err
	}

	for _, interval := range session.State.Contexts[ctxId].Intervals {
		if interval.End.Time.IsZero() {
			intervals = append(intervals, interval.Id)
		}
	}

	return intervals, nil
}

func (session *Session) endInterval(ctxId string, now ctxtime.ZonedTime) error {
  if err := session.ValidateContextExists(ctxId); err != nil {
		return err
	}
  intervals, err := session.GetActiveIntervals(ctxId);
	if err != nil {
		return err
	}

	state := session.State

	for _, intervalId := range intervals {
		ctx := state.Contexts[ctxId]
		interval := ctx.Intervals[intervalId]
		interval.End = now
		interval.Duration = interval.End.Time.Sub(interval.Start.Time)
		ctx.Intervals[interval.Id] = interval
		ctx.Duration += interval.Duration
	//	manager.PublishContextEvent(state.Contexts[id], now, END_INTERVAL, map[string]string{
	//		"duration": interval.Duration.String(),
	//	})
	}
	return nil
}
