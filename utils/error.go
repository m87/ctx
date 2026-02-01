package utils

import ctxlog "github.com/m87/ctx/log"

func Check(err error) {
	if err != nil {
		ctxlog.Logger.Error("Fatal error", "error", err)
		panic(err)
	}
}

func CheckM(err error, msg string) {
	if err != nil {
		ctxlog.Logger.Error("Fatal error: "+msg, "error", err)
		panic(err)
	}
}
