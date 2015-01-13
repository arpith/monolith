package main

import (
	"bufio"
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestPulsar(t *testing.T) {
	pulsar := NewPulsar()
	waiter := &sync.WaitGroup{}
	for i := 0; i < 100; i++ {
		waiter.Add(1)
		listener := pulsar.Listen()
		go func() {
			t.Log("adding listener")
			message := <-listener
			if message == "PULSE" {
				waiter.Done()
			}
		}()
	}
	forgottenListener := pulsar.Listen()
	pulsar.Forget(forgottenListener)
	go func() {
		<-forgottenListener
		t.Error("should not have sent on this channel")
	}()

	success := make(chan bool, 1)
	go func() {
		waiter.Wait()
		success <- true
	}()

	pulsar.Pulse("PULSE")

	select {
	case <-success:
	case <-time.After(2 * time.Second):
		t.Error("No pulse received")
	}
}

func TestPubsub(t *testing.T) {
	monolithServer := httptest.NewServer(http.HandlerFunc(NewBroadcastHandler()))
	url := monolithServer.URL + "/channel"
	success := make(chan bool, 1)
	group := &sync.WaitGroup{}
	for i := 1; i < 100; i++ {
		group.Add(1)

		go func(t *testing.T, group *sync.WaitGroup) {
			resp, err := http.Get(url)
			if err != nil {
				t.Error("should have been able to make the connection")
			}
			defer resp.Body.Close()
			reader := bufio.NewReader(resp.Body)
			data := ""
			for {
				line, err := reader.ReadBytes('\n')
				line = bytes.TrimSpace(line)
				t.Log("Received SSE message message ", string(line))
				if err != nil {
					break
				}
				if strings.Contains(string(line), "PING") {
					group.Done()
				}
			}
			if err != nil {
				t.Error("Shouldn't be an error", data)
			}
		}(t, group)
	}

	go func(group *sync.WaitGroup) {
		group.Wait()
		success <- true
	}(group)

	go func() {
		time.Sleep(time.Second / 100)
		http.Post(url, "text", nil)
	}()

	select {
	case <-success:
	case <-time.After(3 * time.Second):
		t.Error("No message received")
	}
}
