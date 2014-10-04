package main

import (
	"log"
	"net/http"
	"net/url"
)

func deliver(sourceURL *url.URL, destinationURL *url.URL) {
	log.Print("Fetching ", sourceURL.String())
	sourceResponse, err := http.Get(sourceURL.String())
	if err != nil {
		log.Print("Error GETting source: ", err.Error())
		return
	}

	values := destinationURL.Query()
	values.Add("x_monolith_final_url", sourceResponse.Request.URL.String())
	destinationURL.RawQuery = values.Encode()

	postResponse, err := http.Post(destinationURL.String(), "text/html", sourceResponse.Body)
	if err != nil {
		log.Print("Couldn't create POST to ", destinationURL.String(), err.Error())
		return
	}
	log.Print("Deilvered ", sourceURL.String(), " to ", destinationURL.String(), " : ", postResponse.Request.ContentLength, " bytes")
	postResponse.Body.Close()
	sourceResponse.Body.Close()
}

var fetchHandler = func(w http.ResponseWriter, req *http.Request) {
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

}
