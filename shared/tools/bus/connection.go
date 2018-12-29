package bus

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
}
