package events

import (
	"sync"
	"testing"
)

type testListener struct {
	hasRun bool
}

func (t *testListener) Handle(e Event) error {
	t.hasRun = true
	return nil
}

func TestItRegistersEventsToListeners(t *testing.T) {
	eh := newHandler()
	e := Event{
		Tag:  "bongo",
		Data: 1,
	}
	go eh.Watch()
	defer eh.Close()

	if _, ok := eh.listeners[e.Tag]; ok {
		t.Error("there is a listener setup for the event already")
	}

	eh.Register(e.Tag, &testListener{})

	if _, ok := eh.listeners[e.Tag]; !ok {
		t.Error("there is no listener setup for the event already")
	}
}

func TestItDoesntErrorIfNoHandlerHasBeenRegistered(t *testing.T) {
	eh := newHandler()
	e := Event{
		Tag:  "bongo",
		Data: 1,
	}
	go eh.Watch()
	defer eh.Close()

	if _, ok := eh.listeners[e.Tag]; ok {
		t.Error("there is a listener setup for the event already")
	}

	eh.Trigger(e)
}

func TestItCallsARegisteredListener(t *testing.T) {
	eh := newHandler()
	e := Event{
		Tag:  "bongo",
		Data: 1,
	}
	l := &testListener{}
	eh.Watch()
	defer eh.Close()
	eh.Register(e.Tag, l)

	if l.hasRun {
		t.Error("the listener has already run")
	}

	eh.Trigger(e)

	if !l.hasRun {
		t.Error("the listener has not run")
	}
}

func BenchmarkEventHandler(b *testing.B) {
	eh := newHandler()
	e := Event{
		Tag:  "bongo",
		Data: 1,
	}
	l := &testListener{}
	eh.Watch()
	defer eh.Close()
	eh.Register(e.Tag, l)

	for i := 0; i < b.N; i++ {
		eh.Trigger(e)
	}
}

func newHandler() *Handler {
	return &Handler{
		listeners: map[string][]Listener{},
		stream:    make(chan Event),
		close:     make(chan bool),
		wg:        &sync.WaitGroup{},
	}
}
