package stan

import (
	"errors"
	"sync"
	"time"

	"github.com/dmibod/kanban/shared/tools/logger/noop"

	"github.com/dmibod/kanban/shared/tools/bus"
	"github.com/dmibod/kanban/shared/tools/logger"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/stan"
)

// Errors
var (
	ErrNotConnected = errors.New("stan: not connected")
)

const (
	defaultClusterID      = "test-cluster"
	defaultReconnectDelay = time.Second
)

var _ bus.Connection = (*Connection)(nil)

// Connection interface
type Connection struct {
	sync.Mutex
	logger.Logger
	status    chan struct{}
	url       string
	clusterID string
	clientID  string
	natsConn  *nats.Conn
	stanConn  stan.Conn
	stanOpts  []stan.Option
	natsOpts  []nats.Option
	subOpts   []stan.SubscriptionOption
}

// CreateConnection creates new connection
func CreateConnection(opts ...Option) *Connection {
	var o options

	o.natsOpts = []nats.Option{}
	o.stanOpts = []stan.Option{}

	for _, opt := range opts {
		opt(&o)
	}

	url := o.url
	if url == "" {
		url = stan.DefaultNatsURL
	}

	clusterID := o.clusterID
	if clusterID == "" {
		clusterID = defaultClusterID
	}

	l := o.Logger
	if l == nil {
		l = &noop.Logger{}
	}

	var conn *Connection

	o.natsOpts = append(o.natsOpts, nats.DisconnectHandler(func(nc *nats.Conn) { conn.status <- struct{}{} }))
	o.natsOpts = append(o.natsOpts, nats.ReconnectHandler(func(nc *nats.Conn) { conn.status <- struct{}{} }))
	o.stanOpts = append(o.stanOpts, stan.NatsURL(url))

	subOpts := []stan.SubscriptionOption{stan.StartWithLastReceived()}

	if o.durable != "" {
		subOpts = append(subOpts, stan.DurableName(o.durable))
	}

	conn = &Connection{
		Logger:    l,
		url:       url,
		status:    make(chan struct{}, 1),
		clusterID: clusterID,
		clientID:  o.clientID,
		stanOpts:  o.stanOpts,
		natsOpts:  o.natsOpts,
		subOpts:   subOpts,
	}

	return conn
}

// Connect to broker
func (c *Connection) Connect() error {
	if c.IsConnected() {
		return nil
	}

	c.Lock()
	defer c.Unlock()

	err := c.connectNats()
	if err != nil {
		return err
	}

	err = c.connectStan()
	if err != nil {
		return err
	}

	return nil
}

// Disconnect from broker
func (c *Connection) Disconnect() {
	if !c.IsConnected() {
		return
	}

	c.Lock()
	defer c.Unlock()

	if c.stanConn != nil {
		c.Debugln("close stan connection")
		err := c.stanConn.Close()
		c.stanConn = nil
		if err != nil {
			c.Errorln(err)
		}
	}

	if c.natsConn != nil {
		c.Debugln("close nats connection")
		c.natsConn.Close()
		c.natsConn = nil
	}
}

// IsConnected status
func (c *Connection) IsConnected() bool {
	c.Lock()
	defer c.Unlock()
	return c.stanConn != nil && c.stanConn.NatsConn().IsConnected()
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
	c.Lock()
	defer c.Unlock()
	return c.stanConn.Publish(topic, message)
}

// Subscribe for messages
func (c *Connection) Subscribe(topic string, queue string, handler bus.MessageHandler) (interface{}, error) {
	if !c.IsConnected() {
		return nil, ErrNotConnected
	}
	c.Lock()
	defer c.Unlock()
	return c.stanConn.QueueSubscribe(topic, queue, func(msg *stan.Msg) {
		handler.Handle(msg.Data)
	}, c.subOpts...)
}

// Unsubscribe subscription
func (c *Connection) Unsubscribe(handle interface{}) error {
	if !c.IsConnected() {
		return ErrNotConnected
	}
	c.Lock()
	defer c.Unlock()
	if s, ok := handle.(stan.Subscription); ok {
		return s.Unsubscribe()
	}
	return nil
}

func (c *Connection) connectNats() error {
	if c.natsConn != nil {
		return nil
	}

	c.Debugln("connect nats")
	natsConn, err := nats.Connect(c.url, c.natsOpts...)
	if err != nil {
		c.Errorln(err)
		return err
	}

	c.Debugln("nats connected")
	c.natsConn = natsConn
	c.stanOpts = append(c.stanOpts, stan.NatsConn(natsConn))

	return nil
}

func (c *Connection) connectStan() error {
	if c.stanConn != nil {
		return nil
	}

	c.Debugln("connect stan")
	stanConn, err := stan.Connect(c.clusterID, c.clientID, c.stanOpts...)
	if err != nil {
		c.Errorln(err)
		return err
	}

	c.Debugln("stan connected")
	c.stanConn = stanConn

	return nil
}
