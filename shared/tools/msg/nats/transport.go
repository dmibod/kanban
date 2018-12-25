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
	return createReceiver(subj, CreateOperationContext(t.ctx), t.e)
}

func (t *transport) CreateSender(subj string) msg.Sender {
	return createSender(subj, CreateOperationContext(t.ctx), t.e)
}
