package event

// Handler interface
type Handler interface {
	Handle(event interface{})
}

// HandleFunc type
type HandleFunc func(interface{})

// Handle event
func (h HandleFunc) Handle(event interface{}) {
	h(event)
}

// Bus interface
type Bus interface {
	Register(event interface{})
	Listen(handler Handler)
	Fire()
}

type bus struct {
	events   []interface{}
	handlers []Handler
}

// Execute handler
func Execute(handler func(Bus) error) error {
	if handler == nil {
		return nil
	}

	b := &bus{
		events:   []interface{}{},
		handlers: []Handler{},
	}

	return handler(b)
}

// Register event
func (b *bus) Register(event interface{}) {
	if event != nil {
		b.events = append(b.events, event)
	}
}

// Listen events
func (b *bus) Listen(handler Handler) {
	if handler != nil {
		b.handlers = append(b.handlers, handler)
	}
}

// Fire events
func (b *bus) Fire() {
	for _, event := range b.events {
		b.notify(event)
	}
	b.events = b.events[:0]
}

func (b *bus) notify(event interface{}) {
	for _, handler := range b.handlers {
		handler.Handle(event)
	}
}
