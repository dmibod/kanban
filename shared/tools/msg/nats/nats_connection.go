package nats

import (
	"github.com/dmibod/kanban/shared/tools/msg"
	"github.com/nats-io/go-nats"
)

type natsconn struct {
	conn *nats.Conn
}

// CreateNatsConnection creates nats connection
func CreateNatsConnection(url string, opts ...nats.Option) (*natsconn, error) {
	conn, err := nats.Connect(url, opts...)
	if err != nil {
		return nil, err
	}
	return &natsconn{
		conn: conn,
	}, nil
}

// Subscribe new handler
func (c *natsconn) Subscribe(s string, q string, h func([]byte)) (msg.Subscription, error) {
	return c.conn.QueueSubscribe(s, q, func(msg *nats.Msg) {
		h(msg.Data)
	})
}

// Publish message
func (c *natsconn) Publish(s string, m []byte) error {
	return c.conn.Publish(s, m)
}

// Flush pending messages
func (c *natsconn) Flush() error {
	return c.conn.Flush()
}

// Close connection
func (c *natsconn) Close() {
	c.conn.Close()
	c.conn = nil
}
