package bootstrap

import (
	"github.com/m87/ctx/core"
	ctxlog "github.com/m87/ctx/log"
	"github.com/m87/nod"
	"github.com/m87/nod/sqlite"
)

func CreateManager() *core.ContextManager {
	repository := sqlite.NewRepository("ctx.db", ctxlog.Logger, NewMapperRegistry())
	return core.NewContextManager(&core.RealTimeProvider{}, repository)
}

func NewMapperRegistry() *nod.MapperRegistry {
	registry := nod.NewMapperRegistry()
	registry.Register(core.WorkspaceType, "", &core.WorkspaceMapper{})
	return registry
}
