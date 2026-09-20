package core

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func timePointer(t time.Time) *time.Time {
	return &t
}

func setupSplitIntervalTest() (*ContextManager, *Interval) {
	test := NewEmptyTestContextManager()
	interval := &Interval{
		Id:          "interval1",
		ContextId:   "context1",
		Start:       timePointer(time.Date(2024, 6, 1, 9, 0, 0, 0, time.UTC)),
		End:         timePointer(time.Date(2024, 6, 1, 17, 0, 0, 0, time.UTC)),
		Duration:    8 * time.Hour,
		Status:      "inactive",
		WorkspaceId: "workspace1",
	}
	test.Intervals.Seed(interval)
	return test.Manager, interval
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

func copyInterval(interval *Interval) *Interval {
	if interval == nil {
		return nil
	}
	copy := *interval
	return &copy
}

func TestIntervalCreationWithError(t *testing.T) {
	t.Helper()
	manager := NewEmptyTestContextManager().Manager

	interval := &Interval{
		Id:          "interval1",
		ContextId:   "context1",
		Start:       timePointer(time.Date(2024, 6, 1, 9, 0, 0, 0, time.UTC)),
		End:         timePointer(time.Date(2024, 6, 1, 17, 0, 0, 0, time.UTC)),
		Duration:    8 * time.Hour,
		Status:      "inactive",
		WorkspaceId: "workspace1",
	}

	t.Run("CreateInterval with nil interval", func(t *testing.T) {
		_, err := manager.CreateInterval(nil)
		require.Error(t, err)
		assert.IsType(t, &IntervalCreationError{}, err)
		assert.Equal(t, "cannot create interval: interval is required", err.Error())
	})

	t.Run("CreateInterval with invalid start", func(t *testing.T) {
		invalidInterval := copyInterval(interval)
		invalidInterval.Start = nil

		_, err := manager.CreateInterval(invalidInterval)
		require.Error(t, err)
		assert.IsType(t, &IntervalCreationError{}, err)
		assert.Equal(t, "cannot create interval: start time is required", err.Error())
	})

	t.Run("CreateInterval with invalid end", func(t *testing.T) {
		invalidInterval := copyInterval(interval)
		invalidInterval.End = nil

		_, err := manager.CreateInterval(invalidInterval)
		require.Error(t, err)
		assert.IsType(t, &IntervalCreationError{}, err)
		assert.Equal(t, "cannot create interval: end time is required", err.Error())
	})

	t.Run("CreateInterval with invalid context id", func(t *testing.T) {
		invalidInterval := copyInterval(interval)
		invalidInterval.ContextId = ""

		_, err := manager.CreateInterval(invalidInterval)
		require.Error(t, err)
		assert.IsType(t, &IntervalCreationError{}, err)
		assert.Equal(t, "cannot create interval: context id is required", err.Error())
	})

	t.Run("CreateInterval with end before start", func(t *testing.T) {
		invalidInterval := copyInterval(interval)
		invalidInterval.End = timePointer(time.Date(2024, 6, 1, 8, 0, 0, 0, time.UTC))

		_, err := manager.CreateInterval(invalidInterval)
		require.Error(t, err)
		assert.IsType(t, &IntervalCreationError{}, err)
		assert.Equal(t, "cannot create interval: end time must be after start time", err.Error())
	})

	t.Run("CreateInterval with end equal to start", func(t *testing.T) {
		invalidInterval := copyInterval(interval)
		invalidInterval.End = invalidInterval.Start

		_, err := manager.CreateInterval(invalidInterval)
		require.Error(t, err)
		assert.IsType(t, &IntervalCreationError{}, err)
		assert.Equal(t, "cannot create interval: end time must be after start time", err.Error())
	})
}

func TestCreateInterval(t *testing.T) {
	test := NewEmptyTestContextManager()
	test.Contexts.Seed(&Context{Id: "context1", WorkspaceId: "workspace1"})
	start := time.Date(2024, 6, 1, 9, 0, 0, 0, time.UTC)
	end := start.Add(8 * time.Hour)
	interval := &Interval{
		ContextId: "context1",
		Start:     &start,
		End:       &end,
		Status:    "completed",
	}

	id, err := test.Manager.CreateInterval(interval)

	require.NoError(t, err)
	assert.NotEmpty(t, id)
	assert.Equal(t, 8*time.Hour, interval.Duration)
	assert.Equal(t, "workspace1", interval.WorkspaceId)
}

func TestCreateActiveIntervalWithoutEnd(t *testing.T) {
	test := NewEmptyTestContextManager()
	test.Contexts.Seed(&Context{Id: "context1", WorkspaceId: "workspace1"})
	start := time.Date(2024, 6, 1, 9, 0, 0, 0, time.UTC)
	interval := &Interval{
		ContextId: "context1",
		Start:     &start,
		Status:    IntervalStatusActive,
	}

	id, err := test.Manager.CreateInterval(interval)

	require.NoError(t, err)
	assert.NotEmpty(t, id)
	assert.Nil(t, interval.End)
	assert.Zero(t, interval.Duration)
}
