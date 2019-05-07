package event

import (
	"context"
)

// Handler interface
type Handler interface {
	Handle(context.Context, interface{}) error
}

// HandleFunc type
type HandleFunc func(context.Context, interface{}) error

// Handle event
func (h HandleFunc) Handle(ctx context.Context, event interface{}) error {
	return h(ctx, event)
}

// Bus interface
type Bus interface {
	Register(event interface{})
	Listen(handler Handler)
	Fire(context.Context) error
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
func (b *bus) Fire(ctx context.Context) error {
	for _, event := range b.events {
		if err := b.notify(ctx, event); err != nil {
			return err
		}
	}
	b.events = b.events[:0]
	return nil
}

func (b *bus) notify(ctx context.Context, event interface{}) error {
	for _, handler := range b.handlers {
		if err := handler.Handle(ctx, event); err != nil {
			return err
		}
	}
	return nil
}
