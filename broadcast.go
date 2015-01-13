package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

//Pulsar is an emitter that allows broadcasting and listening to messages via channels
type Pulsar struct {
	listeners map[chan string]bool
	views     map[<-chan string]chan string
	lock      sync.Locker
}

//Listen creates a new receiver channel which acts as a subscription. In order to prevent leaks, always return a channel after use via `Forget`
func (p *Pulsar) Listen() <-chan string {
	listener := make(chan string)
	p.listeners[listener] = true
	p.views[listener] = listener
	return listener
}

//Pulse sends a message to all listening channels
func (p *Pulsar) Pulse(message string) {
	p.lock.Lock()
	for c := range p.listeners {
		c <- message
	}
	p.lock.Unlock()
}

//Forget removes a channel from the list of receivers
func (p *Pulsar) Forget(view <-chan string) {
	p.lock.Lock()
	delete(p.listeners, p.views[view])
	delete(p.views, view)
	p.lock.Unlock()
}

//NewPulsar creates a new Pulsar that can be used for PubSub
func NewPulsar() *Pulsar {
	return &Pulsar{listeners: make(map[chan string]bool), views: make(map[<-chan string]chan string), lock: &sync.Mutex{}}
}

type broker struct {
	channels           map[string]*Pulsar
	pulsarCreationLock sync.Locker
}

func (b *broker) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	channelName := r.URL.Path
	b.pulsarCreationLock.Lock()
	channel, ok := b.channels[channelName]
	if !ok {
		channel = NewPulsar()
		b.channels[channelName] = channel
	}
	b.pulsarCreationLock.Unlock()
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
		defer channel.Forget(listener)
		for {
			select {
			case msg := <-listener:
				fmt.Fprintf(w, "data: %s\n\n", msg)
				f.Flush()
			case <-closer.CloseNotify():
				return
			case <-time.After(60 * time.Second):
				return
			}
		}
	case "POST":
		channel.Pulse("PING")
	}
}

// NewBroadcastHandler creates a new handler that handles pub sub
func NewBroadcastHandler() func(w http.ResponseWriter, req *http.Request) {
	broker := &broker{channels: make(map[string]*Pulsar), pulsarCreationLock: &sync.Mutex{}}
	return broker.ServeHTTP
}
