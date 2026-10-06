package storage

import (
	"sort"
	"testing"
	"time"

	"github.com/m87/ctx/core"
	"github.com/stretchr/testify/assert"
)

func TestIntervalRepository(t *testing.T) {
	t.Helper()

	t.Run("Save and retrieve interval", func(t *testing.T) {
		storage, _ := CreateTestInMemoryStorage()
		repo := NewIntervalRepository(storage.DB)

		interval := &core.Interval{
			ContextId:   "context1",
			Start:       &time.Time{},
			Duration:    time.Hour,
			Status:      "active",
			WorkspaceId: "workspace1",
		}

		id, err := repo.Save(interval)
		assert.NoError(t, err)
		assert.NotEmpty(t, id)

		retrievedInterval, err := repo.GetById(id)
		assert.NoError(t, err)
		assert.NotNil(t, retrievedInterval)
		assert.Equal(t, interval.ContextId, retrievedInterval.ContextId)
		assert.Equal(t, interval.Duration, retrievedInterval.Duration)
		assert.Equal(t, interval.Status, retrievedInterval.Status)
		assert.Equal(t, interval.WorkspaceId, retrievedInterval.WorkspaceId)
	})

	t.Run("Delete interval", func(t *testing.T) {
		storage, _ := CreateTestInMemoryStorage()
		repo := NewIntervalRepository(storage.DB)

		interval := &core.Interval{
			ContextId:   "context2",
			Start:       &time.Time{},
			Duration:    time.Hour,
			Status:      "active",
			WorkspaceId: "workspace2",
		}

		id, err := repo.Save(interval)
		assert.NoError(t, err)
		assert.NotEmpty(t, id)

		err = repo.Delete(id)
		assert.NoError(t, err)

		deletedInterval, err := repo.GetById(id)
		assert.NoError(t, err)
		assert.Nil(t, deletedInterval)
	})

	t.Run("SaveAll returns transaction errors and rolls back", func(t *testing.T) {
		storage, err := CreateTestInMemoryStorage()
		assert.NoError(t, err)
		repo := NewIntervalRepository(storage.DB)
		start := time.Now()

		_, err = repo.SaveAll([]*core.Interval{
			{ContextId: "context-1", Start: &start, WorkspaceId: "workspace-1"},
			{ContextId: "context-2", Start: nil, WorkspaceId: "workspace-1"},
		})
		assert.Error(t, err)

		var count int64
		assert.NoError(t, storage.DB.Model(&IntervalEntity{}).Count(&count).Error)
		assert.Zero(t, count)
	})
}

func TestListOverlapping(t *testing.T) {
	storage, err := CreateTestInMemoryStorage()
	assert.NoError(t, err)
	repo := NewIntervalRepository(storage.DB)
	start := time.Date(2026, time.September, 27, 10, 0, 0, 0, time.UTC)
	end := start.Add(time.Hour)
	before := start.Add(-time.Hour)
	after := end.Add(time.Hour)
	zero := time.Time{}
	rows := []*core.Interval{
		{Id: "touch-start", ContextId: "c", WorkspaceId: "w", Start: &before, End: &start},
		{Id: "touch-end", ContextId: "c", WorkspaceId: "w", Start: &end, End: &after},
		{Id: "other-workspace", ContextId: "c", WorkspaceId: "other", Start: &start, End: &end},
		{Id: "open-end", ContextId: "c", WorkspaceId: "w", Start: ptr(start.Add(30 * time.Minute))},
		{Id: "zero-end", ContextId: "c", WorkspaceId: "w", Start: &before, End: &zero},
		{Id: "covers-period", ContextId: "c", WorkspaceId: "w", Start: &before, End: &after},
	}
	for _, row := range rows {
		_, err := repo.Save(row)
		assert.NoError(t, err)
	}
	got, err := repo.ListOverlapping("w", start, end)
	assert.NoError(t, err)
	assert.Equal(t, []string{"covers-period", "open-end", "zero-end"}, intervalIDs(got))
}

func ptr(value time.Time) *time.Time {
	return &value
}

func intervalIDs(intervals []*core.Interval) []string {
	ids := make([]string, len(intervals))
	for i, interval := range intervals {
		ids[i] = interval.Id
	}
	sort.Strings(ids)
	return ids
}
