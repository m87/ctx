package util

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"log"
	"strings"
)

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
