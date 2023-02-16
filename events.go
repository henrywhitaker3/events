package events

type Event struct {
	Tag  string
	Data any
}

type Listener interface {
	Handle(Event) error
}

type EventHandler interface {
	// Register an listener for an event type
	Register(string, Listener)

	// Trigger an event
	Trigger(Event)

	// Watch the stream chan for incoming events
	Watch()

	// Stop the watcher and waits for all events to be processed
	Close()
}
