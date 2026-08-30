package storage

import (
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
		assert.Error(t, err)
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
