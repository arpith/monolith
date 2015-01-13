package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

// Deliver POSTs the content of the source URL to the destination URL.
func Deliver(sourceURL *url.URL, destinationURL *url.URL) {
	log.Print("Fetching ", sourceURL.String())
	sourceResponse, err := http.Get(sourceURL.String())
	if err != nil {
		log.Print("Error GETting source: ", err.Error())
		return
	}

	data, err := ioutil.ReadAll(sourceResponse.Body)
	if err != nil {
		log.Print("Unable to parse response")
		return
	}
	dataReader := bytes.NewReader(data)

	postResponse, err := http.Post(destinationURL.String(), "text/html", dataReader)
	if err != nil {
		log.Print("Couldn't create POST to ", destinationURL.String(), err.Error())
		return
	}
	log.Print("Deilvered ", sourceURL.String(), " to ", destinationURL.String(), " : ", postResponse.Request.ContentLength, " bytes")
	postResponse.Body.Close()
	sourceResponse.Body.Close()
}

// NewFetchHandler creates a fetch handler. The handler take a `src` and `dest` parameter and posts the response from `src` into `dest`.
func NewFetchHandler() func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		log.Println("Starting FETCH request.")
		if err := req.ParseForm(); err != nil {
			http.Error(w, "Unable to parse request", 400)
			return
		}
		sourceURL, err := url.ParseRequestURI(req.FormValue("src"))
		if err != nil {
			http.Error(w, "Please include a valid URL as the `src` parameter - the URL that you want to fetch.", 400)
			return
		}
		destinationURL, err := url.ParseRequestURI(req.FormValue("dest"))
		if err != nil {
			http.Error(w, "Please include a valid URL as the `dest` parameter - the URL that you want to POST the fetched page to.", 400)
			return
		}
		go Deliver(sourceURL, destinationURL)
	}
}
