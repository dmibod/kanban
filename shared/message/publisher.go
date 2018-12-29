package message

import (
	"github.com/dmibod/kanban/shared/tools/bus"
)

// Publisher interface
type Publisher interface {
	// Publish message
	Publish(message []byte) error
}

// CreatePublisher creates "topic" publisher
func CreatePublisher(topic string) Publisher {
	return &publisher{
		topic: topic,
	}
}

type publisher struct {
	topic string
}

// Publish message
func (p *publisher) Publish(message []byte) error {
	return bus.Publish(p.topic, message)
}
