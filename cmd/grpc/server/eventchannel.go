package main

import (
	"fmt"
	"math/rand"
	"time"
)

type Metric struct {
	src    string
	dst    string
	metric int
}

type EventHandler struct {
	listenersCh map[string]chan interface{}
	doneCh      chan bool
}

var eventChannel *EventHandler

func init() {
	eventChannel = &EventHandler{
		listenersCh: make(map[string]chan interface{}),
	}
	eventChannel.Start()
}

func (ec *EventHandler) RegisterListener(id string, listenerCh chan interface{}) {
	ec.listenersCh[id] = listenerCh
}

func (ec *EventHandler) UnregisterListener(id string) {
	delete(ec.listenersCh, id)
}

func (ec *EventHandler) Start() {
	go func() {
		ticker := time.NewTicker(300 * time.Millisecond)
		for {
			select {
			case <-ticker.C:
				fmt.Println("Event Channel generates sample metric")
				for key, ch := range ec.listenersCh {
					if ch != nil {
						ec.listenersCh[key] <- Metric{src: "1.1.1.1", dst: "2.2.2.2", metric: rand.Intn(100)}
					}
				}
			case <-ec.doneCh:
				fmt.Println("Received done")
				return
			}
		}

	}()
}
