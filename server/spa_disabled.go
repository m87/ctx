//go:build !allinone

package server

import "net/http"

func registerSpaHandler() http.Handler { return nil }
