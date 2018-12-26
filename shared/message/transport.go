package message

import (
	"context"
	"github.com/dmibod/kanban/shared/tools/logger"
	"github.com/dmibod/kanban/shared/tools/msg"
	"github.com/dmibod/kanban/shared/tools/msg/nats"
)

// CreateTransport create new transport
func CreateTransport(ctx context.Context, e nats.OperationExecutor, l logger.Logger) msg.Transport {
	return nats.CreateTransport(ctx, e, l)
}
