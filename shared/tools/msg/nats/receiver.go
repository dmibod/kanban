package nats

import (
	"context"
	"sync"

	"github.com/dmibod/kanban/shared/tools/msg"
)

type subscription struct {
	h msg.Receive
	q string
	u Subscription
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

	err := r.e.Execute(r.ctx, func(ctx context.Context, conn Connection) error {

		u, e := conn.Subscribe(r.s, q, func(msg []byte) {
			h(msg)
		})

		if e == nil {
			s.u = u
		}

		return e
	})

	r.Lock()
	r.subscriptions = append(r.subscriptions, s)
	r.Unlock()

	return err
}
