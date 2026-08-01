package core

import "time"

type TimeProvider interface {
	Now() time.Time
}

type RealTimeProvider struct{}

func (provider *RealTimeProvider) Now() time.Time {
	return time.Now().UTC()
}

func NewTimer() *RealTimeProvider {
	return &RealTimeProvider{}
}
