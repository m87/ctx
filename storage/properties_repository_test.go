package storage

import (
	"testing"

	"github.com/m87/ctx/core"
	"github.com/stretchr/testify/require"
)

func TestPropertiesRepository(t *testing.T) {
	t.Helper()

	t.Run("Save system properties", func(t *testing.T) {
		storage, err := CreateTestInMemoryStorage()
		require.NoError(t, err)
		repo := NewPropertiesRepository(storage.DB)

		properties := &core.SystemInfo{
			DatabaseVersion: "1.2.3",
			ClientId:        "client-1",
		}

		err = repo.Save(properties)
		require.NoError(t, err)

		savedProperties, err := repo.Load()
		require.NoError(t, err)
		require.Equal(t, properties, savedProperties)
	})
}
