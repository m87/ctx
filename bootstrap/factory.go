package bootstrap

import (
	"strings"

	"github.com/m87/ctx/core"
	ctxlog "github.com/m87/ctx/log"
	"github.com/m87/nod"
	"github.com/m87/nod/sqlite"
	"github.com/spf13/viper"
)

func CreateManager() *core.ContextManager {
	viper.SetDefault("database.path", "ctx.db")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()
	viper.ReadInConfig()
	repository, _ := sqlite.NewRepository(viper.GetString("database.path"), ctxlog.Logger, NewMapperRegistry())
	return core.NewContextManager(
		&core.RealTimeProvider{},
		NewContextRepository(repository),
		NewIntervalRepository(repository),
		NewWorkspaceRepository(repository),
	)
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
