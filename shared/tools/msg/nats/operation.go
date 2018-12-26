package nats

import (
	"context"
	"github.com/dmibod/kanban/shared/tools/msg"
)

// OperationContext declares operation context
type OperationContext struct {
	ctx context.Context
}

// CreateOperationContext creates OperationContext
func CreateOperationContext(ctx context.Context) *OperationContext {
	if ctx == nil {
		ctx = context.TODO()
	}
	return &OperationContext{
		ctx: ctx,
	}
}

type Connection interface {
	Subscribe(string, string, func([]byte)) (msg.Subscription, error)
	Publish(string, []byte) error
	Flush() error
	Close()
}

// Operation to be performed on nats connection
type Operation func(context.Context, Connection) error
