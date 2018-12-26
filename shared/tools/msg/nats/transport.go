package nats

import (
	"context"

	"github.com/dmibod/kanban/shared/tools/logger"

	"github.com/dmibod/kanban/shared/tools/msg"
)

var _ msg.Transport = (*transport)(nil)

type transport struct {
	context.Context
	OperationExecutor
	logger.Logger
}

// CreateTransport creates new transport
func CreateTransport(c context.Context, e OperationExecutor, l logger.Logger) msg.Transport {
	return &transport{
		OperationExecutor: e,
		Context:           c,
		Logger:            l,
	}
}

func (t *transport) Subscriber(subj string) msg.Subscriber {
	return createSubscriber(subj, CreateOperationContext(t.Context), t.OperationExecutor, t.Logger)
}

func (t *transport) Publisher(subj string) msg.Publisher {
	return createPublisher(subj, CreateOperationContext(t.Context), t.OperationExecutor, t.Logger)
}
