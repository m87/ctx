package bootstrap

import (
	"time"

	"github.com/m87/ctx/core"
	ctxlog "github.com/m87/ctx/log"
	"github.com/m87/nod/sqlite"
)

func NewTextContextManager(current core.ZonedTime) *core.ContextManager {
	repository, _ := sqlite.NewRepository(":memory:", ctxlog.Logger, NewMapperRegistry())
	return core.NewContextManager(
		NewTestTimeProvider(current),
		NewContextRepository(repository),
		NewIntervalRepository(repository),
		NewWorkspaceRepository(repository),
	)
}

type TestTimeProvider struct {
	current core.ZonedTime
}

func NewTestTimeProvider(current core.ZonedTime) *TestTimeProvider {
	return &TestTimeProvider{
		current: current,
	}
}

func (p *TestTimeProvider) Now() core.ZonedTime {
	return p.current
}

func (p *TestTimeProvider) Advance(d time.Duration) {
	p.current = core.ZonedTime{
		Time: p.current.Time.Add(d),
	}
}
