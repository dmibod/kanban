package nats

import (
	"context"

	"github.com/dmibod/kanban/shared/tools/logger"

	"github.com/dmibod/kanban/shared/tools/msg"
)

var _ msg.Transport = (*transport)(nil)

type transport struct {
	ctx context.Context
	e   OperationExecutor
	l   logger.Logger
}

func CreateTransport(ctx context.Context, e OperationExecutor, l logger.Logger) msg.Transport {
	return &transport{
		ctx: ctx,
		e:   e,
		l:   l,
	}
}

func (t *transport) CreateReceiver(subj string) msg.Receiver {
	return createReceiver(subj, CreateOperationContext(t.ctx), t.e, t.l)
}

func (t *transport) CreateSender(subj string) msg.Sender {
	return createSender(subj, CreateOperationContext(t.ctx), t.e, t.l)
}
