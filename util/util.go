package util

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"log"
	"strings"

	"github.com/m87/ctx/ctx_model"
	"github.com/m87/ctx/ctx_store"
	"github.com/m87/ctx/events_model"
	"github.com/m87/ctx/events_store"
)

type statePatch func(*ctx_model.State)

type eventsPatch func(*events_model.EventRegistry)

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
