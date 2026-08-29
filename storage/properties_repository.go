package storage

import (
	"github.com/m87/ctx/core"
	"gorm.io/gorm"
)

type PropertiesRepository struct {
	db *gorm.DB
}

func NewPropertiesRepository(db *gorm.DB) *PropertiesRepository {
	return &PropertiesRepository{db: db}
}

func (r *PropertiesRepository) Load() (*core.SystemInfo, error) {
	var entity Properties
	result := r.db.Where("id = ?", systemPropertiesId).Limit(1).Find(&entity)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	return &core.SystemInfo{
		DatabaseVersion: entity.DatabaseVersion,
		ClientId:        entity.ClientId,
	}, nil
}

func (r *PropertiesRepository) Save(info *core.SystemInfo) error {
	if info == nil {
		return nil
	}

	return r.db.Save(&Properties{
		Id:              systemPropertiesId,
		DatabaseVersion: info.DatabaseVersion,
		ClientId:        info.ClientId,
	}).Error
}

var _ core.SystemInfoRepository = (*PropertiesRepository)(nil)
