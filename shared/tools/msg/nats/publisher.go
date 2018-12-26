package nats

import (
	"context"

	"github.com/dmibod/kanban/shared/tools/logger"
)

type publisher struct {
	s string
	OperationExecutor
	ctx *OperationContext
	logger.Logger
}

func createPublisher(s string, c *OperationContext, e OperationExecutor, l logger.Logger) *publisher {
	return &publisher{
		OperationExecutor: e,
		s:                 s,
		ctx:               c,
		Logger:            l,
	}
}

func (p *publisher) Publish(msg []byte) error {
	return p.Execute(p.ctx, func(ctx context.Context, conn Connection) error {

		err := conn.Publish(p.s, msg)

		if err == nil {
			err = conn.Flush()
		}

		return err
	})
}
