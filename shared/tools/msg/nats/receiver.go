package nats

import (
	"context"
	"sync"

	"github.com/dmibod/kanban/shared/tools/logger"
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
	notify        chan bool
	watchRunning  bool
	l             logger.Logger
}

func createReceiver(s string, c *OperationContext, e OperationExecutor, l logger.Logger) *receiver {
	return &receiver{
		e:             e,
		s:             s,
		ctx:           c,
		subscriptions: []*subscription{},
		notify:        make(chan bool),
		l:             l,
	}
}

func (r *receiver) Receive(q string, h msg.Receive) error {
	s := &subscription{q: q, h: h}

	err := r.subscribe(s)

	r.Lock()
	defer r.Unlock()

	r.subscriptions = append(r.subscriptions, s)

	if !r.watchRunning {
		go r.watch()
		r.watchRunning = true
	}

	return err
}

func (r *receiver) watch() {
	r.l.Debugln("watch for executor signals")
	r.e.Notify(r.notify)
	for {
		select {
		case <-r.ctx.ctx.Done():
			return
		case alive := <-r.notify:
			r.l.Debugf("signal from executor: %v\n", alive)
			if alive {
				r.recover()
			} else {
				r.release()
			}
		}
	}
}

func (r *receiver) recover() {
	r.Lock()
	defer r.Unlock()
	for _, s := range r.subscriptions {
		if s.u == nil {
			r.l.Debugf("recover: %+v\n", s)
			r.subscribe(s)
		}
	}
}

func (r *receiver) subscribe(s *subscription) error {
	return r.e.Execute(r.ctx, func(ctx context.Context, conn Connection) error {

		u, e := conn.Subscribe(r.s, s.q, func(msg []byte) {
			s.h(msg)
		})

		if e == nil {
			s.u = u
		}

		return e
	})
}

func (r *receiver) release() {
	r.Lock()
	defer r.Unlock()
	for _, s := range r.subscriptions {
		if s.u != nil {
			r.l.Debugf("release: %+v\n", s)
			r.unsubscribe(s)
		}
	}
}

func (r *receiver) unsubscribe(s *subscription) error {
	err := s.u.Unsubscribe()
	s.u = nil
	return err
}
