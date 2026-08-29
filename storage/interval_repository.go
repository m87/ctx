package storage

import (
	"sort"
	"time"

	"github.com/google/uuid"
	"github.com/m87/ctx/core"
	"gorm.io/gorm"
)


type IntervalRepository struct {
	db *gorm.DB
}	


func NewIntervalRepository(db *gorm.DB) *IntervalRepository {
	return &IntervalRepository{db: db}
}

func (r *IntervalRepository) GetById(id string) (*core.Interval, error) {
	var entity IntervalEntity
	if err := r.db.First(&entity, "id = ?", id).Error; err != nil {
		return nil, err
	}

	return entity.ToModel(), nil
}

func (r *IntervalRepository) Save(interval *core.Interval) (string, error) {
	return saveInterval(r.db, interval)
}

func saveInterval(db *gorm.DB, interval *core.Interval) (string, error) {
	if interval == nil {
		return "", nil
	}
	if interval.Id == "" {
		interval.Id = uuid.NewString()
	}
	entity := NewIntervalEntityFromModel(interval)
	if err := db.Save(entity).Error; err != nil {
		return "", err
	}

	return entity.Id, nil
}

func (r *IntervalRepository) Delete(id string) error {
	if err := r.db.Delete(&IntervalEntity{}, "id = ?", id).Error; err != nil {
		return err 
	}
	return nil
}


func (r *IntervalRepository) DeleteByContextId(contextId string) error {
	if contextId == "" {
		return nil
	}

	return r.db.Delete(&IntervalEntity{}, "context_id  = ?", contextId).Error
}

func (r *IntervalRepository) ListByContextId(contextId string) ([]*core.Interval, error) {
	var entities []*IntervalEntity

	if err := r.db.Find(&entities, "context_id = ?", contextId).Error; err != nil {
		return nil, err
	}

	intervals := make([]*core.Interval, len(entities))
	for i, entity := range entities {
		intervals[i] = entity.ToModel()
	}

	sort.SliceStable(intervals, func(i, j int) bool {
		if intervalTime(intervals[i].Start).Equal(intervalTime(intervals[j].Start)) {
			return intervalTime(intervals[i].End).After(intervalTime(intervals[j].End))
		}
		return intervalTime(intervals[i].Start).After(intervalTime(intervals[j].Start))
	})
	return intervals, nil
}

func (r *IntervalRepository) GetActiveIntervalByContextId(contextId string) (*core.Interval, error) {
	var entity IntervalEntity
	if err := r.db.Where("context_id = ? AND status = ?", contextId, "active").First(&entity).Error; err != nil {
		return nil, err 
	}
	
	return entity.ToModel(), nil
}

func (r *IntervalRepository) ListByDay(date time.Time, workspaceId string) ([]*core.Interval, error) {
	location := date.Location()
	if location == nil {
		location = time.UTC
	}
	localDayStart := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, location)
	dayStart := localDayStart.UTC()
	dayEnd := localDayStart.AddDate(0, 0, 1).UTC()
	
	var all []IntervalEntity
	if err := r.db.Where("workspace_id = ? AND start <= ?", workspaceId, dayEnd).Find(&all).Error; err != nil {
		return nil, err
	}

	result := make([]*core.Interval, 0, len(all))
	for _, interval := range all {
		if interval.End == nil || interval.End.IsZero() {
			result = append(result, interval.ToModel())
			continue
		}
		if !interval.End.Before(dayStart) {
			result = append(result, interval.ToModel())
		}
	}
	return result, nil
}

func intervalTime(value *time.Time) time.Time {
	if value == nil {
		return time.Time{}
	}
	return *value
}

func (r *IntervalRepository) List() ([]*core.Interval, error) {
	var entities []*IntervalEntity
	if err := r.db.Find(&entities).Error; err != nil {
		return nil, err
	}

	intervals := make([]*core.Interval, len(entities))
	for i, entity := range entities {
		intervals[i] = entity.ToModel()
	}

	return intervals, nil
}

func (r *IntervalRepository) ListToSync(limit int) ([]*core.Interval, error) {
	return r.List()
}

func (r *IntervalRepository) SaveAll(intervals []*core.Interval) ([]string, error) {
	ids := []string{}

	r.db.Transaction(func(tx *gorm.DB) error {
		for _, interval := range intervals {
			id, err := saveInterval(tx, interval)
			if err != nil {
				return err
			}
			ids = append(ids, id)
		}
		return nil
	})

	return ids, nil
}
