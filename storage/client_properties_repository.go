package storage

import (
	"github.com/m87/ctx/core"
	"gorm.io/gorm"
)

type ClientPropertiesRepository struct {
	db *gorm.DB
}

func NewClientPropertiesRepository(db *gorm.DB) *ClientPropertiesRepository {
	return &ClientPropertiesRepository{db: db}
}

func (r *ClientPropertiesRepository) Load() (*core.Settings, error) {
	var entity ClientProperties
	if err := r.db.First(&entity, "id = ?", clientPropertiesId).Error; err != nil {
		return nil, err
	}

	values := make(map[string]string, 3)
	if entity.Theme != "" {
		values["client.general.theme"] = entity.Theme
	}
	if entity.FirstDay != "" {
		values["client.general.firstDay"] = entity.FirstDay
	}
	if entity.Timezone != "" {
		values["client.general.timeZone"] = entity.Timezone
	}

	return core.NewSettings(values), nil
}

func (r *ClientPropertiesRepository) Save(settings *core.Settings) error {
	if settings == nil {
		return nil
	}

	values := settings.Values()
	return r.db.Save(&ClientProperties{
		Id:       clientPropertiesId,
		Theme:    values["client.general.theme"],
		FirstDay: values["client.general.firstDay"],
		Timezone: values["client.general.timeZone"],
	}).Error
}

var _ core.SettingsRepository = (*ClientPropertiesRepository)(nil)
