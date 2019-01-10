package domain

// EventHandler interface
type EventHandler interface {
	Handle(interface{})
}

// EventSource interface
type EventSource interface {
	Raise()
}

// EventRegistry interface
type EventRegistry interface {
	Register(event interface{})
}

// EventManager type
type EventManager struct {
	events   []interface{}
	handlers []EventHandler
}

// CreateEventManager instance
func CreateEventManager() *EventManager {
	return &EventManager{
		events:   []interface{}{},
		handlers: []EventHandler{},
	}
}

// Register event
func (m *EventManager) Register(event interface{}) {
	if event != nil {
		m.events = append(m.events, event)
	}
}

// Listen events
func (m *EventManager) Listen(handler EventHandler) {
	if handler != nil {
		m.handlers = append(m.handlers, handler)
	}
}

// Raise events
func (m *EventManager) Raise() {
	for _, event := range m.events {
		m.notify(event)
	}
	m.events = m.events[:0]
}

func (m *EventManager) notify(event interface{}) {
	for _, handler := range m.handlers {
		handler.Handle(event)
	}
}
