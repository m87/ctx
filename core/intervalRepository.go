package core

type IntervalRepository interface {
	GetById(id string) (*Interval, error)
	Save(interval *Interval) (string, error)
	Delete(id string) error
}
