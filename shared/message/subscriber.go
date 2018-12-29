package message

import (
	"github.com/dmibod/kanban/shared/tools/bus"
)

// Subscriber interface
type Subscriber interface {
	// Subscribe message handler
	Subscribe(handler bus.MessageHandler) bus.Subscription
}

// CreateSubscriber creates "topic" subscriber
func CreateSubscriber(topic string) Subscriber {
	return &subscriber{
		topic: topic,
	}
}

type subscriber struct {
	topic string
}

// Subscribe message handler
func (s *subscriber) Subscribe(handler bus.MessageHandler) bus.Subscription {
	return bus.Subscribe(s.topic, handler)
}
