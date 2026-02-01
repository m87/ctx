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
	return core.NewContextManager(&core.RealTimeProvider{}, repository)
}

func NewMapperRegistry() *nod.MapperRegistry {
	registry := nod.NewMapperRegistry()
	registry.Register(core.WorkspaceType, "", &core.WorkspaceMapper{})
	return registry
}
