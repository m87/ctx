package util

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"log"
	"strings"
	"time"

	"github.com/m87/ctx/ctx"
	"github.com/m87/ctx/ctx_model"
	"github.com/m87/ctx/ctx_store"
	"github.com/m87/ctx/events_model"
	"github.com/m87/ctx/events_store"
)

type TimeProvider interface {
	Now() time.Time
}

type Time struct{}

func (Time) Now() time.Time { return time.Now().Local() }

type Runtime struct {
	time TimeProvider
}

type Providers struct {
	time TimeProvider
}

type Context interface {
}

type ContextImpl struct {
	state          *ctx_model.Context
	eventsRegistry *events_model.EventRegistry
	providers      Providers
}

func CreateContext() (*Context, error) {
	return nil, nil
}

type PatchContext struct {
	State          ctx_model.State
	EventsRegistry events_model.EventRegistry
	TimeProvider   TimeProvider
}

type patch func(appContext *PatchContext)

type statePatch func(*ctx_model.State)

type eventsPatch func(*events_model.EventRegistry)

func Apply(fn patch) {
	patchContext := PatchContext{
		State:          ctx_store.Load(),
		EventsRegistry: events_store.Load(),
		TimeProvider:   Time{},
	}
	fn(&patchContext)
	ctx_store.Save(&patchContext.State)
	events_store.Save(&patchContext.EventsRegistry)
}

func ApplyEventsPatch(fn eventsPatch) {
	eventsRegisty := events_store.Load()
	fn(&eventsRegisty)
	events_store.Save(&eventsRegisty)
}

func ApplyPatch(fn statePatch) {
	state := ctx_store.Load()
	fn(&state)
	ctx_store.Save(&state)
}

func ReadEvents(fn eventsPatch) {
	eventsRegistry := events_store.Load()
	fn(&eventsRegistry)
}

func Read(fn statePatch) {
	state := ctx_store.Load()
	fn(&state)
}

func Id(arg string, isRaw bool) (string, error) {
	id := strings.TrimSpace(arg)
	if id == "" {
		return "", errors.New("")
	}

	if !isRaw {
		id = GenerateId(id)
	}
	return id, nil
}

func Check(err error, msg string) {
	if err != nil {
		log.Panicln(msg)
	}
}

func Warn(err error, msg string) {
	if err != nil {
		log.Println(msg)
	}
}

func GenerateId(desc string) string {
	h := sha256.New()
	h.Write([]byte(strings.ToLower(desc)))
	bs := h.Sum(nil)
	return hex.EncodeToString(bs)
}

func CreateManager() *ctx.ContextManager {

}
