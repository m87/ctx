package storage

import "gorm.io/gorm"

type PropertiesRepository struct {
	db *gorm.DB
}

func NewPropertiesRepository(db *gorm.DB) *PropertiesRepository {
	return &PropertiesRepository{db: db}
}

func (r *PropertiesRepository) Get(id string) (*Properties, error) {
	var entity Properties
	if err := r.db.First(&entity, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	
	return &entity, nil
}

func (r *PropertiesRepository) Save(properties *Properties) error {
	if properties == nil {
		return nil
	}
	return r.db.Save(properties).Error
}

