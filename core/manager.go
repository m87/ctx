package core

import (
	"github.com/m87/ctx/ctx_model"
	ctxtime "github.com/m87/ctx/time"
)

type ContextManager struct {
	ContextStore ctx_model.ContextStore
	EventsStore  ctx_model.EventsStore
	ArchiveStore ctx_model.ArchiveStore
	TimeProvider ctxtime.TimeProvider
}
