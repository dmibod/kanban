package nats

import (
	"context"
	"sync"

	"github.com/dmibod/kanban/shared/tools/logger"
	"github.com/dmibod/kanban/shared/tools/msg"
)

type subscription struct {
	h msg.MessageHandler
	q string
	u msg.Subscription
}

func (s *subscription) Unsubscribe() error {
	err := s.u.Unsubscribe()
	s.u = nil
	return err
}

type subscriber struct {
	sync.Mutex
	logger.Logger
	OperationExecutor
	ctx           *OperationContext
	s             string
	subscriptions []*subscription
	notify        chan bool
	watchRunning  bool
}

func createSubscriber(s string, c *OperationContext, e OperationExecutor, l logger.Logger) *subscriber {
	return &subscriber{
		OperationExecutor: e,
		s:                 s,
		ctx:               c,
		subscriptions:     []*subscription{},
		notify:            make(chan bool),
		Logger:            l,
	}
}

func (s *subscriber) Subscribe(q string, h msg.MessageHandler) (msg.Subscription, error) {
	sub := &subscription{q: q, h: h}

	err := s.subscribe(sub)

	s.Lock()
	defer s.Unlock()

	s.subscriptions = append(s.subscriptions, sub)

	if !s.watchRunning {
		go s.watch()
		s.watchRunning = true
	}

	return sub, err
}

func (s *subscriber) watch() {
	s.Debugln("watch for executor signals")
	s.Notify(s.notify)
	for {
		select {
		case <-s.ctx.ctx.Done():
			return
		case alive := <-s.notify:
			s.Debugf("signal from executor: %v\n", alive)
			if alive {
				s.recover()
			} else {
				s.release()
			}
		}
	}
}

func (s *subscriber) recover() {
	s.Lock()
	defer s.Unlock()
	for _, sub := range s.subscriptions {
		if sub.u == nil {
			s.Debugf("recover: %+v\n", sub)
			s.subscribe(sub)
		}
	}
}

func (s *subscriber) subscribe(sub *subscription) error {
	return s.Execute(s.ctx, func(ctx context.Context, conn Connection) error {

		u, e := conn.Subscribe(s.s, sub.q, func(msg []byte) {
			sub.h(msg)
		})

		if e == nil {
			sub.u = u
		}

		return e
	})
}

func (s *subscriber) release() {
	s.Lock()
	defer s.Unlock()
	for _, sub := range s.subscriptions {
		if sub.u != nil {
			s.Debugf("release: %+v\n", sub)
			sub.Unsubscribe()
		}
	}
}
