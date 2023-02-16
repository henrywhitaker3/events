package events

import (
	"sync"
)

type Handler struct {
	listeners map[string][]Listener
	stream    chan Event
	close     chan bool
	wg        *sync.WaitGroup
}

func NewHandler() EventHandler {
	return &Handler{
		listeners: make(map[string][]Listener),
		stream:    make(chan Event),
		close:     make(chan bool),
		wg:        &sync.WaitGroup{},
	}
}

func (a *Handler) Register(t string, l Listener) {
	a.listeners[t] = append(a.listeners[t], l)
}

func (a *Handler) Trigger(e Event) {
	a.stream <- e
}

// Run the listeners for the events in their own goroutines
func (a *Handler) run(e Event) {
	for _, l := range a.getListenersForEvent(e) {
		a.wg.Add(1)
		go func(l Listener) {
			l.Handle(e)
			a.wg.Done()
		}(l)
	}
}

func (a *Handler) Watch() {
	a.wg.Add(1)
	go func() {
		for {
			select {
			case <-a.close:
				a.wg.Done()
				return
			case e := <-a.stream:
				a.run(e)
			}
		}
	}()
}

func (a *Handler) Close() {
	a.close <- true
	a.wg.Wait()
	close(a.close)
	close(a.stream)
}

// Get the listeners for an event type
func (a *Handler) getListenersForEvent(e Event) []Listener {
	return a.listeners[e.Tag]
}
