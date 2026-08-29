package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"
)



func TestPropertiesRepository(t *testing.T) {
	t.Helper()

	t.Run("Save properties", func(t *testing.T) {
		storage, _ := CreateTestInMemoryStorage()
		repo := NewPropertiesRepository(storage.DB)

		properties := &Properties{
			Id:       "test",
			Theme:    "dark",
			FirstDay: "Monday",
			Timezone: "UTC",
		}

		err := repo.Save(properties)
		assert.NoError(t, err)

		savedProperties, err := repo.Get("test")
		assert.NoError(t, err)
		assert.NotNil(t, savedProperties)
		assert.Equal(t, "dark", savedProperties.Theme)
		assert.Equal(t, "Monday", savedProperties.FirstDay)
		assert.Equal(t, "UTC", savedProperties.Timezone)
	})
}

