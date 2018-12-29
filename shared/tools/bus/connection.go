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

// Connection interface
type Connection interface {
	// Connect to broker
	Connect() error
	// Disconnect from broker
	Disconnect() 
	// IsConnected status
	IsConnected() bool
	// Status of connection
	Status() <-chan struct{}
	// Publish message
	Publish(topic string, message []byte) error
	// Subscribe for messages
	Subscribe(topic string, queue string, handler MessageHandler) (interface{}, error)
	// Unsubscribe subscription
	Unsubscribe(handle interface{}) error
}
