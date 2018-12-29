package nats

import (
	"errors"
	"sync"
	"time"

	"github.com/dmibod/kanban/shared/tools/logger/noop"

	"github.com/dmibod/kanban/shared/tools/bus"
	"github.com/dmibod/kanban/shared/tools/logger"
	"github.com/nats-io/go-nats"
)

// Errors
var (
	ErrNotConnected = errors.New("nats: not connected")
)

const (
	defaultReconnectDelay = time.Second
)

var _ bus.Connection = (*Connection)(nil)

// Connection interface
type Connection struct {
	mu       sync.Mutex
	status   chan struct{}
	url      string
	name     string
	logger   logger.Logger
	natsConn *nats.Conn
	natsOpts []nats.Option
}

// CreateConnection creates new connection
func CreateConnection(opts ...Option) *Connection {
	var o options

	o.natsOpts = []nats.Option{}

	for _, opt := range opts {
		opt(&o)
	}

	url := o.url
	if url == "" {
		url = nats.DefaultURL
	}

	l := o.Logger
	if l == nil {
		l = &noop.Logger{}
	}

	var conn *Connection

	o.natsOpts = append(o.natsOpts, nats.DisconnectHandler(func(nc *nats.Conn) { conn.status <- struct{}{} }))
	o.natsOpts = append(o.natsOpts, nats.ReconnectHandler(func(nc *nats.Conn) { conn.status <- struct{}{} }))

	conn = &Connection{
		logger:   l,
		url:      url,
		status:   make(chan struct{}, 1),
		natsOpts: o.natsOpts,
	}

	return conn
}

// Connect to broker
func (c *Connection) Connect() error {
	if c.IsConnected() {
		return nil
	}

	c.mu.Lock()
	c.mu.Unlock()

	c.logger.Debugln("connect nats")

	natsConn, err := nats.Connect(c.url, c.natsOpts...)
	if err != nil {
		c.logger.Errorln(err)
		return err
	}

	c.logger.Debugln("nats connected")
	c.natsConn = natsConn

	return nil
}

// Disconnect from broker
func (c *Connection) Disconnect() {
	if !c.IsConnected() {
		return
	}

	c.mu.Lock()
	c.mu.Unlock()

	c.logger.Debugln("close nats connection")
	c.natsConn.Close()
	c.natsConn = nil
}

// IsConnected status
func (c *Connection) IsConnected() bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.natsConn != nil && c.natsConn.IsConnected()
}

// Status of connection
func (c *Connection) Status() <-chan struct{} {
	return c.status
}

// Publish message
func (c *Connection) Publish(topic string, message []byte) error {
	if !c.IsConnected() {
		return ErrNotConnected
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.natsConn.Publish(topic, message)
}

// Subscribe for messages
func (c *Connection) Subscribe(topic string, queue string, handler bus.MessageHandler) (interface{}, error) {
	if !c.IsConnected() {
		return nil, ErrNotConnected
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.natsConn.QueueSubscribe(topic, queue, func(msg *nats.Msg) {
		handler.Handle(msg.Data)
	})
}

// Unsubscribe subscription
func (c *Connection) Unsubscribe(handle interface{}) error {
	if !c.IsConnected() {
		return ErrNotConnected
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	if s, ok := handle.(nats.Subscription); ok {
		return s.Unsubscribe()
	}
	return nil
}
