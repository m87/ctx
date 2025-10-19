package util

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"log"
	"strings"
)

func Id(arg string, isRawId bool) (string, error) {
	id := strings.TrimSpace(arg)
	if id == "" {
		return "", errors.New("")
	}

	if !isRawId {
		id = GenerateId(id)
	}
	return id, nil
}

func Checkm(err error, msg string) {
	if err != nil {
		log.Panicln(msg)
	}
}

func Check(err error) {
	if err != nil {
		log.Panicln(err)
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

func Contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func Remove(slice []string, item string) []string {
	for i, s := range slice {
		if s == item {
			return append(slice[:i], slice[i+1:]...)
		}
	}
	return slice
}
