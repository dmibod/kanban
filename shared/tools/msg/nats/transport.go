package nats

import (
	"context"

	"github.com/dmibod/kanban/shared/tools/msg"
)

var _ msg.Transport = (*transport)(nil)

type transport struct {
	e   OperationExecutor
	ctx context.Context
}

func CreateTransport(ctx context.Context, e OperationExecutor) msg.Transport {
	return &transport{
		e:   e,
		ctx: ctx,
	}
}

func (t *transport) CreateReceiver(subj string) msg.Receiver {
	return &receiver{
		e:             t.e,
		s:             subj,
		ctx:           CreateOperationContext(t.ctx),
		subscriptions: []*subscription{},
	}
}

func (t *transport) CreateSender(subj string) msg.Sender {
	return &sender{
		e:   t.e,
		s:   subj,
		ctx: CreateOperationContext(t.ctx),
	}
}
