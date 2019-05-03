package bus

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/dmibod/kanban/shared/tools/logger"
	"github.com/dmibod/kanban/shared/tools/logger/noop"
)

// Errors
var (
	ErrInvalidState      = errors.New("bus: invalid state")
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

// Logger set logger for default bus
func Logger(l logger.Logger) {
	defaultBus.Lock()
	defaultBus.Logger = l
	defaultBus.Unlock()

	defaultBus.ensureLogger()
}

var defaultBus = &bus{
	subscriptions: make(map[int]*subscription),
}

// ConnectAndServe starts bus
func ConnectAndServe(ctx context.Context, conn Connection, tran Transport) error {
	if conn == nil {
		return ErrInvalidConnection
	}

	if tran == nil {
		return ErrInvalidTransport
	}

	if ctx == nil {
		ctx = context.Background()
	}

	defaultBus.Lock()

	if defaultBus.state {
		defaultBus.Unlock()
		return ErrInvalidState
	}

	defaultBus.state = true

	defaultBus.Connection = conn
	defaultBus.Transport = tran

	defaultBus.Unlock()

	defaultBus.ensureLogger()

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
			if err := conn.Connect(); err == nil {
				defaultBus.attachAll()
				return nil
			}
		}
	}
}

// Disconnect bus
func Disconnect() {
	defaultBus.detachAll()

	defaultBus.Lock()
	defer defaultBus.Unlock()

	if !defaultBus.state {
		return
	}

	defaultBus.Disconnect()
	defaultBus.state = false

	defaultBus.Connection = nil
	defaultBus.Transport = nil
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
	logger.Logger
	state         bool
	connState     bool
	subKey        int
	subscriptions map[int]*subscription
}

func (b *bus) ensureLogger() {
	b.Lock()
	defer b.Unlock()
	if defaultBus.Logger == nil {
		defaultBus.Logger = &noop.Logger{}
	}
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
	if b.connState {
		if err := b.attachOne(s); err != nil {
			b.Errorln(err)
		}
	}
	return s
}

func (b *bus) unsubscribe(s *subscription) error {
	b.Lock()
	defer b.Unlock()
	if s != nil {
		delete(b.subscriptions, s.key)
		if b.connState {
			return b.Unsubscribe(s.handle)
		}
	}
	return nil
}

func (b *bus) attachAll() {
	b.Lock()
	defer b.Unlock()
	if b.connState {
		return
	}
	b.connState = true
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
	if !b.connState {
		return
	}
	b.connState = false
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
