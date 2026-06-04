package bootstrap

import (
	"github.com/m87/ctx/core"
	"github.com/m87/nod"
)

type SettingsRepository struct {
	repository nod.TypedRepository[core.Settings]
}

func NewSettingsRepository(repository *nod.Repository) *SettingsRepository {
	return &SettingsRepository{
		repository: nod.As[core.Settings](repository),
	}
}

func (r *SettingsRepository) Save(settings *core.Settings) error {
	_, err := r.repository.Save(settings)
	return err
}

func (r *SettingsRepository) Load() (*core.Settings, error) {
	return r.repository.Query().KindEquals(core.SettingsType).KV().First()
}
