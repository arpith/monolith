package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestDelivery(t *testing.T) {
	message := "ABCD1234"
	sourceServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, message)
	}))
	defer sourceServer.Close()
	success := make(chan bool, 1)
	destinationServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Error(err.Error())
		}
		if strings.Contains(string(body), message) {
			success <- true
		}
	}))
	defer destinationServer.Close()

	monolithServer := httptest.NewServer(http.HandlerFunc(fetchHandler))
	defer monolithServer.Close()

	url := monolithServer.URL + "?src=" + sourceServer.URL + "&dest=" + destinationServer.URL
	http.Get(url)

	timeout := make(chan bool, 1)
	go func() {
		time.Sleep(1 * time.Second)
		timeout <- true
	}()

	select {
	case <-success:

	case <-timeout:
		t.Error("No call made")
	}

}
