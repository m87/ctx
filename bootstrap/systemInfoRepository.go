package bootstrap

import (
	"github.com/m87/ctx/core"
	"github.com/m87/nod"
)

type SystemInfoRepository struct {
	repository nod.TypedRepository[core.SystemInfo]
}

func NewSystemInfoRepository(repository *nod.Repository) *SystemInfoRepository {
	return &SystemInfoRepository{repository: nod.As[core.SystemInfo](repository)}
}

func (r *SystemInfoRepository) Load() (*core.SystemInfo, error) {
	return r.repository.Query().KindEquals(core.SystemInfoType).KV().First()
}

func (r *SystemInfoRepository) Save(info *core.SystemInfo) error {
	_, err := r.repository.Save(info)
	return err
}
