package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/codegangsta/negroni"
)

var client = &http.Client{}

func deliver(sourceURL *url.URL, destinationURL *url.URL) {
	log.Print("Fetching ", sourceURL.String())
	sourceResponse, err := http.Get(sourceURL.String())
	if err != nil {
		log.Print(err.Error())
		return
	}
	log.Print("Fetched ", sourceURL.String())
	defer sourceResponse.Body.Close()
	log.Print("Sending ", sourceURL.String(), " to ", destinationURL.String())
	postResponse, err := http.Post(destinationURL.String(), "text/html", sourceResponse.Body)
	if err != nil {
		log.Print("Couldn't create POST to ", destinationURL.String(), err.Error())
		return
	}
	defer postResponse.Body.Close()
	log.Print("Sent ", sourceURL.String(), " to ", destinationURL.String(), " : ", postResponse.Request.ContentLength)

}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "Monolith")
	})
	mux.HandleFunc("/fetch", func(w http.ResponseWriter, req *http.Request) {
		if err := req.ParseForm(); err != nil {
			http.Error(w, "Unable to parse request", 400)
		}
		sourceURL, err := url.ParseRequestURI(req.FormValue("src"))
		if err != nil {
			http.Error(w, "Please include a valid URL as the `src` parameter - the URL that you want to fetch.", 400)
		}
		destinationURL, err := url.ParseRequestURI(req.FormValue("dest"))
		if err != nil {
			http.Error(w, "Please include a valid URL as the `dest` parameter - the URL that you want to POST the fetched page to.", 400)
		}

		go deliver(sourceURL, destinationURL)

	})

	n := negroni.Classic()
	n.UseHandler(mux)

	port := strings.TrimSpace(os.Getenv("PORT"))
	if port == "" {
		port = "3000"
	}
	n.Run(":" + port)
}
