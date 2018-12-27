package bus

// Subscription unsubscribe existing subscription
type Subscription interface {
	// Unsubscribe subscription
	Unsubscribe() error
}
