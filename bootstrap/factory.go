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
	repository := sqlite.NewRepository(viper.GetString("database.path"), ctxlog.Logger, NewMapperRegistry())
	return core.NewContextManager(&core.RealTimeProvider{}, nod.As[core.Context](repository), nod.As[core.Workspace](repository), nod.As[core.Project](repository), nod.As[core.Interval](repository))
}

func NewMapperRegistry() *nod.MapperRegistry {
	registry := nod.NewMapperRegistry()
	nod.RegisterMapper(registry, &core.IntervalMapper{})
	nod.RegisterMapper(registry, &core.ContextMapper{})
	nod.RegisterMapper(registry, &core.WorkspaceMapper{})
	nod.RegisterMapper(registry, &core.ProjectMapper{})
	return registry
}
