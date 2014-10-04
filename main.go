package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "Monolith")
	})
	mux.HandleFunc("/fetch", fetchHandler)

	port := strings.TrimSpace(os.Getenv("PORT"))
	if port == "" {
		port = "3000"
	}
	http.ListenAndServe(":"+port, mux)
}
