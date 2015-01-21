package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/sudhirj/strobe"
)

type broker struct {
	channels map[string]*strobe.Strobe
	sync.RWMutex
}

func (b *broker) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	channelName := r.URL.Path
	b.Lock()
	channel, ok := b.channels[channelName]
	if !ok {
		channel = strobe.NewStrobe()
		b.channels[channelName] = channel
	}
	b.Unlock()
	switch r.Method {
	case "GET":
		f, ok := w.(http.Flusher)
		if !ok {
			http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
			return
		}

		closer, ok := w.(http.CloseNotifier)
		if !ok {
			http.Error(w, "Closing unsupported!", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")

		listener := channel.Listen()
		defer channel.Off(listener)
		for {
			select {
			case msg := <-listener:
				fmt.Fprintf(w, "data: %s\n\n", msg)
				f.Flush()
			case <-closer.CloseNotify():
				return
			case <-time.After(300 * time.Second):
				return
			}
		}
	case "POST":
		go channel.Pulse("PING")
	}
}

// NewBroadcastHandler creates a new handler that handles pub sub
func NewBroadcastHandler() func(w http.ResponseWriter, req *http.Request) {
	broker := &broker{channels: make(map[string]*strobe.Strobe)}
	return broker.ServeHTTP
}
