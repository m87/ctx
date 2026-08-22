package core

import (
	"fmt"
	"time"
)


type IntervalSplitError struct {
	IntervalId string
	Reason     string
}

func (e *IntervalSplitError) Error() string {
	return fmt.Sprintf("cannot split interval %q: %s", e.IntervalId, e.Reason)
}

func NewIntervalSplitError(intervalId, reason string) *IntervalSplitError {
	return &IntervalSplitError{
		IntervalId: intervalId,
		Reason:     reason,
	}
}

func (m *ContextManager) SplitInterval(intervalId string, splitTime time.Time) error {
	if intervalId == "" {
		return NewIntervalSplitError(intervalId, "interval id is required")
	}

	if !timeIsSet(&splitTime) {
		return NewIntervalSplitError(intervalId, "split time is required")
	}

	interval, err := m.IntervalRepository.GetById(intervalId)
	if err != nil {
		return NewIntervalSplitError(intervalId, "interval not found")
	}

	if interval.Start.After(splitTime) || (timeIsSet(interval.End) && interval.End.Before(splitTime)) {
		return NewIntervalSplitError(intervalId, "split time is outside the interval range")
	}

	if interval.Status == IntervalStatusActive {
		return NewIntervalSplitError(intervalId, "cannot split an active interval")
	}

	if !timeIsSet(interval.End) {
		return NewIntervalSplitError(intervalId, "cannot split an interval with no end time")
	}

	firstInterval := &Interval{
		ContextId:   interval.ContextId,
		Start:       interval.Start,
		End:         &splitTime,
		Duration:    durationBetween(interval.Start, &splitTime),
		Status:      interval.Status,
		WorkspaceId: interval.WorkspaceId,
	}

	secondInterval := &Interval{
		ContextId:   interval.ContextId,
		Start:       &splitTime,
		End:         interval.End,
		Duration:    durationBetween(&splitTime, interval.End),
		Status:      interval.Status,
		WorkspaceId: interval.WorkspaceId,
	}

	m.RunInTransaction(func(txManager *ContextManager) error {
		if err := txManager.IntervalRepository.Delete(intervalId); err != nil {
			return NewIntervalSplitError(intervalId, "failed to delete original interval")
		}

		if _, err := txManager.IntervalRepository.Save(firstInterval); err != nil {
			return NewIntervalSplitError(intervalId, "failed to save first split interval")
		}

		if _, err := txManager.IntervalRepository.Save(secondInterval); err != nil {
			return NewIntervalSplitError(intervalId, "failed to save second split interval")
		}

		return nil
	})

	return nil
}




