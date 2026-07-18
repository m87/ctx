package bootstrap

import (
	"github.com/m87/ctx/core"
	"github.com/m87/nod"
)

type SystemInfoRepository struct {
	scope *nod.NodeScope[core.SystemInfo]
}

func NewSystemInfoRepository(repository *nod.Repository) *SystemInfoRepository {
	return &SystemInfoRepository{scope: nod.Nodes[core.SystemInfo](repository)}
}

func (r *SystemInfoRepository) Load() (*core.SystemInfo, error) {
	return r.scope.Query().
		Where(nod.NodeFields.Kind.Equals(core.SystemInfoType)).
		WithKV().
		FindFirst()
}

func (r *SystemInfoRepository) Save(info *core.SystemInfo) error {
	_, err := r.scope.SaveNode(info)
	return err
}
