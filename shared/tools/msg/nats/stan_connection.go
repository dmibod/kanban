package nats

import (
	"github.com/nats-io/go-nats-streaming"
)

type stanconn struct {
	conn stan.Conn
	opts []stan.SubscriptionOption
}

func CreateStanConnection(url, clusterID, clientID string, opts ...stan.Option) (Connection, error) {
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

func (c *stanconn) Subscribe(s string, q string, h func([]byte)) (Subscription, error) {
	return c.conn.QueueSubscribe(s, q, func(msg *stan.Msg) {
		h(msg.Data)
	}, c.opts...)
}

func (c *stanconn) Publish(s string, m []byte) error {
	return c.conn.Publish(s, m)
}

func (c *stanconn) Flush() error {
	return c.conn.NatsConn().Flush()
}

func (c *stanconn) Close() {
	c.conn.Close()
	c.conn = nil
}
