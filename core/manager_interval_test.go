package core

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func timePointer(t time.Time) *time.Time {
	return &t
}

func setupSplitIntervalTest() (*ContextManager, *Interval) {
	interval := &Interval{
		Id:          "interval1",
		ContextId:   "context1",
		Start:       timePointer(time.Date(2024, 6, 1, 9, 0, 0, 0, time.UTC)),
		End:         timePointer(time.Date(2024, 6, 1, 17, 0, 0, 0, time.UTC)),
		Duration:    8 * time.Hour,
		Status:      "inactive",
		WorkspaceId: "workspace1",
	}
	manager := NewContextManager(nil, nil, NewIntervalRepositoryMock([]*Interval{interval}), nil, nil)
	return manager, interval
}

func TestIntervalSplit(t *testing.T) {
	t.Helper()
	splitTime := time.Date(2024, 6, 1, 13, 0, 0, 0, time.UTC)

	t.Run("SplitInterval with valid input", func(t *testing.T) {
		manager, interval := setupSplitIntervalTest()

		result, err := manager.SplitInterval(interval.Id, splitTime)
		assert.NoError(t, err)
		assert.Equal(t, interval, result.Origin)
		assert.NotEmpty(t, result.SplitResult[0].Id)
		assert.NotEmpty(t, result.SplitResult[1].Id)

		intervals, err := manager.IntervalRepository.ListByContextId(interval.ContextId)
		assert.NoError(t, err)
		assert.Len(t, intervals, 2)

		firstInterval := intervals[0]
		secondInterval := intervals[1]

		assert.Equal(t, interval.ContextId, firstInterval.ContextId)
		assert.Equal(t, interval.ContextId, secondInterval.ContextId)

		assert.Equal(t, interval.Start, firstInterval.Start)
		assert.Equal(t, &splitTime, firstInterval.End)
		assert.Equal(t, splitTime.Sub(*interval.Start), firstInterval.Duration)

		assert.Equal(t, &splitTime, secondInterval.Start)
		assert.Equal(t, interval.End, secondInterval.End)
		assert.Equal(t, interval.End.Sub(splitTime), secondInterval.Duration)
	})

	t.Run("SplitInterval with invalid id", func(t *testing.T) {
		manager, _ := setupSplitIntervalTest()

		_, err := manager.SplitInterval("invalid_id", splitTime)
		assert.Error(t, err)
		assert.IsType(t, &IntervalSplitError{}, err)
		assert.Equal(t, "cannot split interval \"invalid_id\": interval not found", err.Error())
	})

	t.Run("SplitInterval with invalid split time", func(t *testing.T) {
		manager, interval := setupSplitIntervalTest()

		_, err := manager.SplitInterval(interval.Id, time.Time{})
		assert.Error(t, err)
		assert.IsType(t, &IntervalSplitError{}, err)
		assert.Equal(t, "cannot split interval \"interval1\": split time is required", err.Error())
	})

	t.Run("SplitInterval with split time outside interval range", func(t *testing.T) {
		splitTime := time.Date(2024, 6, 1, 18, 0, 0, 0, time.UTC)
		manager, interval := setupSplitIntervalTest()

		_, err := manager.SplitInterval(interval.Id, splitTime)
		assert.Error(t, err)
		assert.IsType(t, &IntervalSplitError{}, err)
		assert.Equal(t, "cannot split interval \"interval1\": split time is outside the interval range", err.Error())

		splitTime = time.Date(2024, 6, 1, 8, 0, 0, 0, time.UTC)
		_, err = manager.SplitInterval(interval.Id, splitTime)
		assert.Error(t, err)
		assert.IsType(t, &IntervalSplitError{}, err)
		assert.Equal(t, "cannot split interval \"interval1\": split time is outside the interval range", err.Error())
	})

	t.Run("SplitInterval with active interval", func(t *testing.T) {
		manager, interval := setupSplitIntervalTest()
		interval.Status = IntervalStatusActive
		manager.IntervalRepository.Save(interval)

		_, err := manager.SplitInterval(interval.Id, splitTime)
		assert.Error(t, err)
		assert.IsType(t, &IntervalSplitError{}, err)
		assert.Equal(t, "cannot split interval \"interval1\": cannot split an active interval", err.Error())

	})

	t.Run("SplitInterval with no end time", func(t *testing.T) {
		manager, interval := setupSplitIntervalTest()
		interval.End = nil
		manager.IntervalRepository.Save(interval)

		_, err := manager.SplitInterval(interval.Id, splitTime)
		assert.Error(t, err)
		assert.IsType(t, &IntervalSplitError{}, err)
		assert.Equal(t, "cannot split interval \"interval1\": cannot split an interval with no end time", err.Error())
	})
}
