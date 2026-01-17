package bootstrap

import (
	"github.com/m87/ctx/core"
	"github.com/m87/ctx/log"
	"github.com/m87/nod"
	"github.com/m87/nod/sqlite"
)

func CreateManager() *core.ContextManager {
	repository := sqlite.NewRepository("ctx.db", log.Logger, NewMapperRegistry())
	return core.NewContextManager(&core.RealTimeProvider{}, repository)
}

func NewMapperRegistry() *nod.MapperRegistry {
	registry := nod.NewMapperRegistry()
	return registry
}
