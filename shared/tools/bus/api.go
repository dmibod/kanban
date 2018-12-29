package bus

import (
	"context"
	"errors"
	"sync"
	"time"
)

// Errors
var (
	ErrInvalidConnection = errors.New("bus: invalid connection")
	ErrInvalidTransport  = errors.New("bus: invalid transport")
	ErrConnectionFailed  = errors.New("bus: connection failed")
)

// SubscribeQueue for messages
func SubscribeQueue(topic string, queue string, handler MessageHandler) Subscription {
	return defaultBus.subscribe(topic, queue, handler)
}

// Subscribe for messages
func Subscribe(topic string, handler MessageHandler) Subscription {
	return SubscribeQueue(topic, "", handler)
}

// Publish message
func Publish(topic string, message []byte) error {
	return defaultBus.Publish(topic, message)
}

var defaultBus = &bus{
	subscriptions: make(map[int]*subscription),
}

// ConnectAndServe starts bus
func ConnectAndServe(ctx context.Context, conn Connection, t Transport) error {
	if conn == nil {
		return ErrInvalidConnection
	}

	if t == nil {
		return ErrInvalidTransport
	}

	if ctx == nil {
		ctx = context.TODO()
	}

	defaultBus.Connection = conn
	defaultBus.Transport = t

	go defaultBus.serve(ctx)

	if conn.IsConnected() {
		defaultBus.attachAll()
		return nil
	}

	timer := time.NewTicker(time.Second * 1)

	for {
		select {
		case <-ctx.Done():
			return ErrConnectionFailed
		case <-timer.C:
			err := conn.Connect()
			if err == nil {
				defaultBus.attachAll()
				return nil
			}
		}
	}
}

type subscription struct {
	key     int
	topic   string
	queue   string
	handler MessageHandler
	handle  interface{}
	unsub   func(*subscription) error
}

func (s *subscription) Unsubscribe() error {
	return s.unsub(s)
}

type bus struct {
	sync.Mutex
	Connection
	Transport
	state         bool
	subKey        int
	subscriptions map[int]*subscription
}

func (b *bus) subscribe(topic string, queue string, handler MessageHandler) *subscription {
	b.Lock()
	defer b.Unlock()
	b.subKey++
	s := &subscription{
		key:     b.subKey,
		topic:   topic,
		queue:   queue,
		handler: handler,
		unsub:   b.unsubscribe,
	}
	b.subscriptions[b.subKey] = s
	return s
}

func (b *bus) unsubscribe(s *subscription) error {
	b.Lock()
	defer b.Unlock()
	if s != nil {
		delete(b.subscriptions, s.key)
		return b.Unsubscribe(s.handle)
	}
	return nil
}

func (b *bus) attachAll() {
	b.Lock()
	defer b.Unlock()
	if b.state {
		return
	}
	b.state = true
	for _, s := range b.subscriptions {
		b.attachOne(s)
	}
}

func (b *bus) attachOne(s *subscription) error {
	if s != nil {
		h, err := b.Subscribe(s.topic, s.queue, s.handler)
		if err == nil {
			s.handle = h
		}
	}
	return nil
}

func (b *bus) detachAll() {
	b.Lock()
	defer b.Unlock()
	if !b.state {
		return
	}
	b.state = false
	for _, s := range b.subscriptions {
		b.detachOne(s)
	}
}

func (b *bus) detachOne(s *subscription) error {
	if s != nil {
		h := s.handle
		s.handle = nil
		return b.Unsubscribe(h)
	}
	return nil
}

func (b *bus) serve(ctx context.Context) {
	for {
		select {
		case <-b.Status():
			if b.IsConnected() {
				b.attachAll()
			} else {
				b.detachAll()
			}
		case <-ctx.Done():
			return
		}
	}
}
