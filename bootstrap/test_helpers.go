package bootstrap

import (
	"time"

	"github.com/m87/ctx/core"
	ctxlog "github.com/m87/ctx/log"
	"github.com/m87/nod/sqlite"
)

func NewTestContextManager(current time.Time) *core.ContextManager {
	repository, _ := sqlite.NewRepository(":memory:", ctxlog.Logger, NewAdapterRegistry())
	return core.NewContextManager(
		NewTestTimeProvider(current),
		NewContextRepository(repository),
		NewIntervalRepository(repository),
		NewWorkspaceRepository(repository),
		NewProjectRepository(repository),
	)
}

type TestTimeProvider struct {
	current time.Time
}

func NewTestTimeProvider(current time.Time) *TestTimeProvider {
	return &TestTimeProvider{
		current: current,
	}
}

func (p *TestTimeProvider) Now() time.Time {
	return p.current
}

func (p *TestTimeProvider) Advance(d time.Duration) {
	p.current = p.current.Add(d).UTC()
}
