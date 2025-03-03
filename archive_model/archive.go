package archive_model

import (
	"github.com/m87/ctx/ctx_model"
	"github.com/m87/ctx/events_model"
)

type ArchiveEntry struct {
	Context ctx_model.Context    `json:"context"`
	Events  []events_model.Event `json:"events"`
}
