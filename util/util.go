package util

import (
	"crypto/sha256"
	"encoding/hex"
	"log"
	"strings"
)

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
