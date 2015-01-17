package main

import (
	"net/http"
	"os"
	"strings"

	"github.com/julienschmidt/httprouter"
)

func main() {
	router := httprouter.New()

	router.HandlerFunc("POST", "/fetch", NewFetchHandler())

	broadcaster := NewBroadcastHandler()
	router.HandlerFunc("GET", "/broadcast/:channel", broadcaster)
	router.HandlerFunc("POST", "/broadcast/:channel", broadcaster)

	port := strings.TrimSpace(os.Getenv("PORT"))
	if port == "" {
		port = "3001"
	}
	http.ListenAndServe(":"+port, router)
}
