package core

import (
	"fmt"
	"time"
)

type IntervalSplitError struct {
	IntervalId string
	Reason     string
}

type IntervalSplitResult struct {
	Origin      *Interval    `json:"origin"`
	SplitResult [2]*Interval `json:"splitResult"`
}

type IntervalSplitUndoError struct {
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

func (e *IntervalSplitUndoError) Error() string {
	return fmt.Sprintf("cannot undo split of interval %q: %s", e.IntervalId, e.Reason)
}

func NewIntervalSplitUndoError(intervalId, reason string) *IntervalSplitUndoError {
	return &IntervalSplitUndoError{
		IntervalId: intervalId,
		Reason:     reason,
	}
}

func (m *ContextManager) SplitInterval(intervalId string, splitTime time.Time) (*IntervalSplitResult, error) {
	if intervalId == "" {
		return nil, NewIntervalSplitError(intervalId, "interval id is required")
	}

	if !timeIsSet(&splitTime) {
		return nil, NewIntervalSplitError(intervalId, "split time is required")
	}

	interval, err := m.IntervalRepository.GetById(intervalId)
	if err != nil || interval == nil {
		return nil, NewIntervalSplitError(intervalId, "interval not found")
	}

	if interval.Start.After(splitTime) || (timeIsSet(interval.End) && interval.End.Before(splitTime)) {
		return nil, NewIntervalSplitError(intervalId, "split time is outside the interval range")
	}

	if interval.Status == IntervalStatusActive {
		return nil, NewIntervalSplitError(intervalId, "cannot split an active interval")
	}

	if !timeIsSet(interval.End) {
		return nil, NewIntervalSplitError(intervalId, "cannot split an interval with no end time")
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

	err = m.RunInTransaction(func(txManager *ContextManager) error {
		if err := txManager.IntervalRepository.Delete(intervalId); err != nil {
			return NewIntervalSplitError(intervalId, "failed to delete original interval")
		}

		firstIntervalId, err := txManager.IntervalRepository.Save(firstInterval)
		if err != nil {
			return NewIntervalSplitError(intervalId, "failed to save first split interval")
		}
		firstInterval.Id = firstIntervalId

		secondIntervalId, err := txManager.IntervalRepository.Save(secondInterval)
		if err != nil {
			return NewIntervalSplitError(intervalId, "failed to save second split interval")
		}
		secondInterval.Id = secondIntervalId

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &IntervalSplitResult{
		Origin:      interval,
		SplitResult: [2]*Interval{firstInterval, secondInterval},
	}, nil
}

func (m *ContextManager) UndoSplitInterval(intervalId string, result *IntervalSplitResult) (*Interval, error) {
	if intervalId == "" {
		return nil, NewIntervalSplitUndoError(intervalId, "interval id is required")
	}
	if result == nil || result.Origin == nil {
		return nil, NewIntervalSplitUndoError(intervalId, "origin interval is required")
	}
	if result.Origin.Id != intervalId {
		return nil, NewIntervalSplitUndoError(intervalId, "origin interval id does not match")
	}

	firstInterval := result.SplitResult[0]
	secondInterval := result.SplitResult[1]
	if firstInterval == nil || firstInterval.Id == "" || secondInterval == nil || secondInterval.Id == "" {
		return nil, NewIntervalSplitUndoError(intervalId, "split interval ids are required")
	}
	if firstInterval.Id == secondInterval.Id || firstInterval.Id == intervalId || secondInterval.Id == intervalId {
		return nil, NewIntervalSplitUndoError(intervalId, "interval ids must be distinct")
	}

	origin := result.Origin
	origin.Duration = durationBetween(origin.Start, origin.End)

	err := m.RunInTransaction(func(txManager *ContextManager) error {
		if err := txManager.IntervalRepository.Delete(firstInterval.Id); err != nil {
			return NewIntervalSplitUndoError(intervalId, "failed to delete first split interval")
		}
		if err := txManager.IntervalRepository.Delete(secondInterval.Id); err != nil {
			return NewIntervalSplitUndoError(intervalId, "failed to delete second split interval")
		}

		restoredId, err := txManager.SaveInterval(origin)
		if err != nil {
			return NewIntervalSplitUndoError(intervalId, "failed to restore origin interval")
		}
		origin.Id = restoredId
		return nil
	})
	if err != nil {
		return nil, err
	}

	return origin, nil
}
