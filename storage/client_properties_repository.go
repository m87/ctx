package storage

import (
	"encoding/json"
	"strings"

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

	values := make(map[string]string)
	if strings.TrimSpace(entity.Values) != "" {
		if err := json.Unmarshal([]byte(entity.Values), &values); err != nil {
			return nil, err
		}
	}
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
	clientValues := make(map[string]string, len(values))
	for key, value := range values {
		if strings.HasPrefix(key, "client.") {
			clientValues[key] = value
		}
	}
	encoded, err := json.Marshal(clientValues)
	if err != nil {
		return err
	}

	return r.db.Save(&ClientProperties{
		Id:       clientPropertiesId,
		Theme:    clientValues["client.general.theme"],
		FirstDay: clientValues["client.general.firstDay"],
		Timezone: clientValues["client.general.timeZone"],
		Values:   string(encoded),
	}).Error
}

var _ core.SettingsRepository = (*ClientPropertiesRepository)(nil)
