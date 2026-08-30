package bootstrap

import (
	"strings"

	"github.com/m87/ctx/core"
	"github.com/m87/ctx/migration"
	"github.com/m87/ctx/storage"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

func CreateManager() (*core.ContextManager, error) {
	storageDB, err := openStorage()
	if err != nil {
		return nil, err
	}
	if err := migration.NewNodDataMigrator(storageDB.DB).MigrateIfNeeded(); err != nil {
		return nil, err
	}

	manager := newContextManager(storageDB.DB)
	if err := manager.EnsureDefaultWorkspace(); err != nil {
		return nil, err
	}

	return manager, nil
}

func newContextManager(db *gorm.DB) *core.ContextManager {
	manager := core.NewContextManager(
		&core.RealTimeProvider{},
		storage.NewContextRepository(db),
		storage.NewIntervalRepository(db),
		storage.NewWorkspaceRepository(db),
		storage.NewProjectRepository(db),
	)
	manager.RunInTransaction = func(fn func(*core.ContextManager) error) error {
		return db.Transaction(func(tx *gorm.DB) error {
			return fn(newContextManager(tx))
		})
	}
	return manager
}

func CreateSettingsManager() (*core.SettingsManager, error) {
	storageDB, err := openStorage()
	if err != nil {
		return nil, err
	}
	if err := migration.NewNodDataMigrator(storageDB.DB).MigrateIfNeeded(); err != nil {
		return nil, err
	}

	manager := newSettingsManager(storageDB.DB)
	if err := manager.InitSettingsIfNotExists(); err != nil {
		return nil, err
	}

	return manager, nil
}

func openStorage() (*storage.Storage, error) {
	viper.SetDefault("database.path", "ctx.db")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()
	viper.ReadInConfig()
	return storage.NewSqliteStorage(viper.GetString("database.path"))
}

func newSettingsManager(db *gorm.DB) *core.SettingsManager {
	return core.NewSettingsManager(storage.NewClientPropertiesRepository(db))
}
