package msg

// MessageHandler is a callback, which is called upon message arrival
type MessageHandler func([]byte)

// Subscription allows to unsubscribe from existing subscription
type Subscription interface {
	// Unsubscribe from existing subscription
	Unsubscribe() error
}

// Subscriber allows to make subscription for receiving messages
type Subscriber interface {
	// Subscribe for receiving messages
	Subscribe(queue string, handler MessageHandler) (Subscription, error)
}

// Publisher allows to publish messages
type Publisher interface {
	// Publish messages
	Publish(message []byte) error
}

// Transport allows to create a Subscriber or Publisher for specific topic
type Transport interface {
	// Subscriber creates a Subscriber for specific topic
	Subscriber(topic string) Subscriber

	// Publisher creates a Publisher for specific topic
	Publisher(topic string) Publisher
}
