package bootstrap

import (
	"time"

	"github.com/m87/ctx/core"
	"github.com/m87/ctx/storage"
)

func NewTestContextManager(current time.Time) *core.ContextManager {
	s, _ := storage.NewSqliteStorage(":memory:")
	return core.NewContextManager(
		NewTestTimeProvider(current),
		storage.NewContextRepository(s.DB),
		storage.NewIntervalRepository(s.DB),
		storage.NewWorkspaceRepository(s.DB),
		storage.NewProjectRepository(s.DB),
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
