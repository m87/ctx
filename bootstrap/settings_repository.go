package bootstrap

import (
	"github.com/m87/ctx/core"
	"github.com/m87/nod"
)

type SettingsRepository struct {
	scope *nod.NodeScope[core.Settings]
}

func NewSettingsRepository(repository *nod.Repository) *SettingsRepository {
	return &SettingsRepository{scope: nod.Nodes[core.Settings](repository)}
}

func (r *SettingsRepository) Save(settings *core.Settings) error {
	_, err := r.scope.SaveNode(settings)
	return err
}

func (r *SettingsRepository) Load() (*core.Settings, error) {
	return r.scope.Query().
		Where(nod.NodeFields.Kind.Equals(core.SettingsType)).
		WithKV().
		FindFirst()
}
