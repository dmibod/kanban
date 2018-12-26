package nats

import (
	"github.com/dmibod/kanban/shared/tools/msg"
	"github.com/nats-io/go-nats-streaming"
)

type stanconn struct {
	conn stan.Conn
	opts []stan.SubscriptionOption
}

// CreateStanConnection creates stan connection
func CreateStanConnection(url, clusterID, clientID string, opts ...stan.Option) (*stanconn, error) {
	conn, err := stan.Connect(clusterID, clientID, append(opts, stan.NatsURL(url))...)
	if err != nil {
		return nil, err
	}
	subOpts := []stan.SubscriptionOption{
		stan.DeliverAllAvailable(),
		stan.DurableName("KANBAN"),
	}
	return &stanconn{
		conn: conn,
		opts: subOpts,
	}, nil
}

// Subscribe new handler
func (c *stanconn) Subscribe(s string, q string, h func([]byte)) (msg.Subscription, error) {
	return c.conn.QueueSubscribe(s, q, func(msg *stan.Msg) {
		h(msg.Data)
	}, c.opts...)
}

// Publish message
func (c *stanconn) Publish(s string, m []byte) error {
	return c.conn.Publish(s, m)
}

// Flush pending messages
func (c *stanconn) Flush() error {
	return c.conn.NatsConn().Flush()
}

// Close connection
func (c *stanconn) Close() {
	c.conn.Close()
	c.conn = nil
}
