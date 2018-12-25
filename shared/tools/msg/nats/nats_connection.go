package nats

import (
	"github.com/nats-io/go-nats"
)

type natsconn struct {
	conn *nats.Conn
}

func CreateNatsConnection(url string, opts ...nats.Option) (*natsconn, error) {
	conn, err := nats.Connect(url, opts...)
	if err != nil {
		return nil, err
	}
	return &natsconn{
		conn: conn,
	}, nil
}

func (c *natsconn) Subscribe(s string, q string, h func([]byte)) (Subscription, error) {
	return c.conn.QueueSubscribe(s, q, func(msg *nats.Msg) {
		h(msg.Data)
	})
}

func (c *natsconn) Publish(s string, m []byte) error {
	return c.conn.Publish(s, m)
}

func (c *natsconn) Flush() error {
	return c.conn.Flush()
}

func (c *natsconn) Close() {
	c.conn.Close()
	c.conn = nil
}
