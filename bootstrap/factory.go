package bootstrap

import (
	"strings"

	"github.com/m87/ctx/core"
	ctxlog "github.com/m87/ctx/log"
	"github.com/m87/nod"
	"github.com/m87/nod/sqlite"
	"github.com/spf13/viper"
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

	newManager := func(repository *nod.Repository) *core.ContextManager {
		return core.NewContextManager(
			&core.RealTimeProvider{},
			NewContextRepository(repository),
			NewIntervalRepository(repository),
			NewWorkspaceRepository(repository),
		)
	}

	if err := repository.Transaction(func(txRepository *nod.Repository) error {
		return newManager(txRepository).EnsureDefaultWorkspace()
	}); err != nil {
		return nil, err
	}

	return newManager(repository), nil
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
	manager := core.NewSettingsManager(NewSettingsRepository(repository))
	if err := manager.InitSettingsIfNotExists(); err != nil {
		return nil, err
	}

	return manager, nil
}

func NewSystemMapperRegistry() *nod.MapperRegistry {
	registry := nod.NewMapperRegistry()
	nod.RegisterMapper(registry, &core.SettingsMapper{})
	return registry
}

func NewMapperRegistry() *nod.MapperRegistry {
	registry := nod.NewMapperRegistry()
	nod.RegisterMapper(registry, &core.IntervalMapper{})
	nod.RegisterMapper(registry, &core.ContextMapper{})
	nod.RegisterMapper(registry, &core.WorkspaceMapper{})
	return registry
}
