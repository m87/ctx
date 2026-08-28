package bootstrap

import (
	"time"

	"github.com/m87/ctx/core"
	ctxlog "github.com/m87/ctx/log"
	"github.com/m87/ctx/storage"
	"github.com/m87/nod/sqlite"
)

func NewTestContextManager(current time.Time) *core.ContextManager {
	repository, _ := sqlite.NewRepository(":memory:", ctxlog.Logger, NewAdapterRegistry())
	s, _ := storage.NewSqliteStorage(":memory:")
	return core.NewContextManager(
		NewTestTimeProvider(current),
		NewContextRepository(repository),
		NewIntervalRepository(repository),
		storage.NewWorkspaceRepository(s.DB),
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
