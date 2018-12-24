package nats

import (
	"context"

	"github.com/nats-io/go-nats"
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

// Operation to be performed on nats connection
type Operation func(context.Context, *nats.Conn) error
