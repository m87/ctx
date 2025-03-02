package archive

import (
	"log"
	"testing"
	"time"

	"github.com/m87/ctx/ctx"
	"github.com/m87/ctx/ctx_model"
	"github.com/m87/ctx/events"
	"github.com/m87/ctx/events_model"
	"github.com/m87/ctx/test"
)

func TestArchiveContextWithEvents(t *testing.T) {
	state := ctx_model.State{
		Contexts:  map[string]ctx_model.Context{},
		CurrentId: "",
	}
	eventsRegistry := events_model.EventRegistry{}

	ctx.Create(&state, test.TestId, test.TestDescription)
	events.Publish(test.CreateTestEvent("test", test.TestId, time.Now().Local()), &eventsRegistry)

	log.Println(eventsRegistry)
}

func TestArchiveContextWithoutEvents(t *testing.T) {

}

func TestErrorOnCurrentContext(t *testing.T) {

}

func TestArchiveAllContexts(t *testing.T) {

}

func TestSkipCurrentContext(t *testing.T) {

}
