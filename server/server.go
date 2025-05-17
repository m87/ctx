package server

import (
	"log"
	"net/http"
)

func Serve() {
	log.Fatal(http.ListenAndServe(":8080", nil))
}
