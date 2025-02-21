package util

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"log"
	"strings"

	"github.com/m87/ctx/ctx_model"
	"github.com/m87/ctx/ctx_store"
)

type stateConsumer func(*ctx_model.State)

func ApplyPatch(fn stateConsumer) {
	state := ctx_store.Load()
	fn(&state)
	ctx_store.Save(&state)
}

func Read(fn stateConsumer) {
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
		log.Fatal(msg)
	}
}

func GenerateId(desc string) string {
	h := sha256.New()
	h.Write([]byte(strings.ToLower(desc)))
	bs := h.Sum(nil)
	return hex.EncodeToString(bs)
}
