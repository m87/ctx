package storage

import (
	"testing"

	"github.com/m87/ctx/core"
	"github.com/stretchr/testify/require"
)

func TestClientPropertiesRepository(t *testing.T) {
	storage, err := CreateTestInMemoryStorage()
	require.NoError(t, err)
	repository := NewClientPropertiesRepository(storage.DB)

	settings := core.NewSettings(map[string]string{
		"client.general.theme":    "dark",
		"client.general.firstDay": "Sunday",
		"client.general.timeZone": "Europe/Warsaw",
		"client.plugin.option":    "preserved",
		"database.path":           "/tmp/ignored.db",
	})

	require.NoError(t, repository.Save(settings))

	saved, err := repository.Load()
	require.NoError(t, err)
	require.Equal(t, map[string]string{
		"client.general.theme":    "dark",
		"client.general.firstDay": "Sunday",
		"client.general.timeZone": "Europe/Warsaw",
		"client.plugin.option":    "preserved",
	}, saved.Values())
}

func TestClientPropertiesRepositoryLoadsPreviousRelationalShape(t *testing.T) {
	storage, err := CreateTestInMemoryStorage()
	require.NoError(t, err)
	require.NoError(t, storage.DB.Migrator().DropTable(&ClientProperties{}))
	require.NoError(t, storage.DB.Exec(`
		CREATE TABLE client_properties (
			id TEXT PRIMARY KEY,
			theme TEXT,
			first_day TEXT,
			timezone TEXT
		)
	`).Error)
	require.NoError(t, storage.DB.Exec(`
		INSERT INTO client_properties (id, theme, first_day, timezone)
		VALUES (?, ?, ?, ?)
	`, clientPropertiesId, "dark", "Monday", "Europe/Warsaw").Error)

	require.NoError(t, Migrate(storage.DB))

	saved, err := NewClientPropertiesRepository(storage.DB).Load()
	require.NoError(t, err)
	require.Equal(t, map[string]string{
		"client.general.theme":    "dark",
		"client.general.firstDay": "Monday",
		"client.general.timeZone": "Europe/Warsaw",
	}, saved.Values())
}

func TestSettingsManagerUsesClientProperties(t *testing.T) {
	storage, err := CreateTestInMemoryStorage()
	require.NoError(t, err)
	manager := core.NewSettingsManager(NewClientPropertiesRepository(storage.DB))

	require.NoError(t, manager.InitSettingsIfNotExists())
	require.NoError(t, manager.SaveClient(map[string]string{
		"client.general.theme": "dark",
		"client.plugin.option": "preserved",
		"database.path":        "/tmp/ignored.db",
	}))

	settings, err := manager.GetClient()
	require.NoError(t, err)
	require.Equal(t, "dark", settings["client.general.theme"])
	require.Equal(t, "preserved", settings["client.plugin.option"])
	require.NotContains(t, settings, "database.path")

	var count int64
	require.NoError(t, storage.DB.Model(&ClientProperties{}).Count(&count).Error)
	require.Equal(t, int64(1), count)
	require.NoError(t, storage.DB.Model(&Properties{}).Count(&count).Error)
	require.Zero(t, count)
}
