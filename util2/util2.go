package util2

import (
	"time"

	"github.com/m87/ctx/ctx"
	"github.com/m87/ctx/ctx_store"
	"github.com/spf13/viper"
)

type RealTimeProvider struct{}

func (provider *RealTimeProvider) Now() time.Time {
	return time.Now().Local()
}

func NewTimer() *RealTimeProvider {
	return &RealTimeProvider{}
}

func CreateManager() *ctx.ContextManager {
	return ctx.New(ctx_store.New(viper.GetString("path")), NewTimer())
}
