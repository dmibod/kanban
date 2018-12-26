package nats

import (
	"github.com/dmibod/kanban/shared/tools/msg"
)

// Connection to message broker
type Connection interface {
	// Subscribe new handler
	Subscribe(topic string, queue string, handler func([]byte)) (msg.Subscription, error)
	// Publish message
	Publish(topic string, message []byte) error
	// Flush pending messages
	Flush() error
	// Close connection
	Close()
}
