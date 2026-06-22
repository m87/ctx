package bootstrap

import (
	"errors"
	"strings"
	"time"

	"github.com/m87/ctx/core"
	ctxlog "github.com/m87/ctx/log"
	"github.com/m87/nod"
	"github.com/m87/nod/sqlite"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

func CreateManager() (*core.ContextManager, error) {
	viper.SetDefault("database.path", "ctx.db")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()
	viper.ReadInConfig()
	repository, err := sqlite.NewRepository(viper.GetString("database.path"), ctxlog.Logger, NewMapperRegistry())
	if err != nil {
		return nil, err
	}

	if err := repository.Transaction(func(txRepository *nod.Repository) error {
		systemRepository := nod.NewRepository(txRepository.DB(), txRepository.Log(), NewSystemMapperRegistry())
		settingsManager := newSettingsManager(systemRepository)
		systemInfo, err := settingsManager.SystemInfoRepository.Load()
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		currentVersion := ""
		if systemInfo != nil {
			currentVersion = systemInfo.DatabaseVersion
		}
		needsMigration, err := core.DatabaseVersionNeedsMigration(currentVersion, core.CurrentDatabaseVersion)
		if err != nil {
			return err
		}
		if !needsMigration {
			return nil
		}

		startedAt := time.Now()
		ctxlog.Logger.Info("Starting database migration", "from_version", currentVersion, "to_version", core.CurrentDatabaseVersion)
		migrated, err := newContextManager(txRepository).EnsureDefaultWorkspaceWithResult()
		if err != nil {
			return err
		}
		if err := settingsManager.SystemInfoRepository.Save(&core.SystemInfo{DatabaseVersion: core.CurrentDatabaseVersion}); err != nil {
			return err
		}
		ctxlog.Logger.Info("Database migration completed", "database_version", core.CurrentDatabaseVersion, "records_updated", migrated, "duration", time.Since(startedAt))
		return nil
	}); err != nil {
		ctxlog.Logger.Error("Database migration failed; transaction rolled back", "error", err)
		return nil, err
	}

	return newContextManager(repository), nil
}

func newContextManager(repository *nod.Repository) *core.ContextManager {
	manager := core.NewContextManager(
		&core.RealTimeProvider{},
		NewContextRepository(repository),
		NewIntervalRepository(repository),
		NewWorkspaceRepository(repository),
	)
	manager.RunInTransaction = func(fn func(*core.ContextManager) error) error {
		return repository.Transaction(func(txRepository *nod.Repository) error {
			return fn(newContextManager(txRepository))
		})
	}
	return manager
}

func CreateSettingsManager() (*core.SettingsManager, error) {
	viper.SetDefault("database.path", "ctx.db")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()
	viper.ReadInConfig()
	repository, err := sqlite.NewRepository(viper.GetString("database.path"), ctxlog.Logger, NewSystemMapperRegistry())
	if err != nil {
		return nil, err
	}
	manager := newSettingsManager(repository)
	if err := manager.InitSettingsIfNotExists(); err != nil {
		return nil, err
	}

	return manager, nil
}

func newSettingsManager(repository *nod.Repository) *core.SettingsManager {
	return core.NewSettingsManager(
		NewSettingsRepository(repository),
		NewSystemInfoRepository(repository),
	)
}

func NewSystemMapperRegistry() *nod.MapperRegistry {
	registry := nod.NewMapperRegistry()
	nod.RegisterMapper(registry, &core.SettingsMapper{})
	nod.RegisterMapper(registry, &core.SystemInfoMapper{})
	return registry
}

func NewMapperRegistry() *nod.MapperRegistry {
	registry := nod.NewMapperRegistry()
	nod.RegisterMapper(registry, &core.IntervalMapper{})
	nod.RegisterMapper(registry, &core.ContextMapper{})
	nod.RegisterMapper(registry, &core.WorkspaceMapper{})
	return registry
}
