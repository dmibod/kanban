package nats

import (
	"context"
	"sync"

	"github.com/dmibod/kanban/shared/tools/msg"
	"github.com/nats-io/go-nats"
)

type subscription struct {
	h msg.Receive
	q string
	u interface {
		Unsubscribe() error
	}
}

type receiver struct {
	sync.Mutex
	ctx           *OperationContext
	s             string
	e             OperationExecutor
	subscriptions []*subscription
}

func createReceiver(s string, c *OperationContext, e OperationExecutor) *receiver {
	return &receiver{
		e:             e,
		s:             s,
		ctx:           c,
		subscriptions: []*subscription{},
	}
}

func (r *receiver) Receive(q string, h msg.Receive) error {
	s := &subscription{q: q, h: h}

	err := r.e.Execute(r.ctx, func(ctx context.Context, conn *nats.Conn) error {

		u, e := conn.QueueSubscribe(r.s, q, func(msg *nats.Msg) {
			h(msg.Data)
		})

		if e == nil {
			s.u = u
			r.Lock()
			r.subscriptions = append(r.subscriptions, s)
			r.Unlock()
		}

		return e
	})

	return err
}
