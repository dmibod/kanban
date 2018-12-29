package bus

// MessageHandler interface
type MessageHandler interface {
	// Handle message
	Handle([]byte)
}

// HandleFunc message func
type HandleFunc func([]byte)

// Handle message
func (f HandleFunc) Handle(m []byte) {
	f(m)
}

// Transport interface
type Transport interface {
	// Publish message
	Publish(topic string, message []byte) error
	// Subscribe for messages
	Subscribe(topic string, queue string, handler MessageHandler) (interface{}, error)
	// Unsubscribe subscription
	Unsubscribe(handle interface{}) error
}
