package bus

import (
	"errors"
	"sync"
)

// Errors
var (
	ErrInvalidConnection = errors.New("bus: invalid connection")
)

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

// Subscription unsubscribe existing subscription
type Subscription interface {
	// Unsubscribe subscription
	Unsubscribe() error
}

// SubscribeQueue for messages
func SubscribeQueue(topic string, queue string, handler MessageHandler) Subscription {
	return defaultBus.createSubscription(topic, queue, handler)
}

// Subscribe for messages
func Subscribe(topic string, handler MessageHandler) Subscription {
	return SubscribeQueue(topic, "", handler)
}

// Publish message
func Publish(topic string, message []byte) error {
	return defaultBus.Publish(topic, message)
}

// Connection interface
type Connection interface {
	// Connect to broker
	Connect() <-chan bool
	// IsConnected status
	IsConnected() bool
	// Publish message
	Publish(topic string, message []byte) error
	// Subscribe for messages
	Subscribe(topic string, queue string, handler MessageHandler) (interface{}, error)
	// Unsubscribe subscription
	Unsubscribe(handle interface{}) error
	// Close connection
	Close() <-chan struct{}
}

var defaultBus = &bus{
	subscriptions: make(map[int]*subscription),
}

// ConnectAndServe starts bus
func ConnectAndServe(conn Connection) error {
	if conn == nil {
		return ErrInvalidConnection
	}

	defaultBus.Connection = conn

	if !conn.IsConnected() {
		status := <-conn.Connect()
		if status {
			defaultBus.subscribeAll()
		}
	}

	go defaultBus.serve()

	return nil
}

type subscription struct {
	bus     *bus
	key     int
	topic   string
	queue   string
	handler MessageHandler
	handle  interface{}
}

func (s *subscription) attach() {
	s.bus.attach(s)
}

func (s *subscription) detach() error {
	return s.bus.detach(s)
}

func (s *subscription) Unsubscribe() error {
	return s.bus.unsubscribe(s)
}

type bus struct {
	sync.Mutex
	Connection
	state         bool
	subKey        int
	subscriptions map[int]*subscription
}

func (b *bus) createSubscription(topic string, queue string, handler MessageHandler) *subscription {
	b.Lock()
	defer b.Unlock()
	b.subKey++
	s := &subscription{
		bus:     b,
		key:     b.subKey,
		topic:   topic,
		queue:   queue,
		handler: handler,
	}
	b.subscriptions[b.subKey] = s
	return s
}

func (b *bus) subscribeAll() {
	b.Lock()
	defer b.Unlock()
	if b.state {
		return
	}
	b.state = true
	for _, s := range b.subscriptions {
		s.attach()
	}
}

func (b *bus) unsubscribeAll() {
	b.Lock()
	defer b.Unlock()
	if !b.state {
		return
	}
	b.state = false
	for _, s := range b.subscriptions {
		s.detach()
	}
}

func (b *bus) attach(s *subscription) error {
	if s != nil {
		h, err := b.Connection.Subscribe(s.topic, s.queue, s.handler)
		if err == nil {
			s.handle = h
		}
	}
	return nil
}

func (b *bus) detach(s *subscription) error {
	b.Lock()
	defer b.Unlock()
	if s != nil {
		h := s.handle
		s.handle = nil
		return b.Connection.Unsubscribe(h)
	}
	return nil
}

func (b *bus) unsubscribe(s *subscription) error {
	b.Lock()
	defer b.Unlock()
	if s != nil {
		delete(b.subscriptions, s.key)
		return b.Connection.Unsubscribe(s.handle)
	}
	return nil
}

func (b *bus) serve() {
	for {
		select {
		case status := <-b.Connect():
			if status {
				b.subscribeAll()
			} else {
				b.unsubscribeAll()
			}
		case <-b.Close():
			return
		}
	}
}
