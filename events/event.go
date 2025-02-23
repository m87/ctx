package events

import "github.com/m87/ctx/events_model"

func Publish(event events_model.Event, registry *events_model.EventRegistry) {
	registry.Events = append(registry.Events, event)
}
