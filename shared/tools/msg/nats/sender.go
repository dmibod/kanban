package nats

import (
	"context"

	"github.com/dmibod/kanban/shared/tools/logger"
)

type sender struct {
	s   string
	e   OperationExecutor
	ctx *OperationContext
	l   logger.Logger
}

func createSender(s string, c *OperationContext, e OperationExecutor, l logger.Logger) *sender {
	return &sender{
		e:   e,
		s:   s,
		ctx: c,
		l:   l,
	}
}

func (s *sender) Send(msg []byte) error {
	return s.e.Execute(s.ctx, func(ctx context.Context, conn Connection) error {

		err := conn.Publish(s.s, msg)

		if err == nil {
			err = conn.Flush()
		}

		return err
	})
}
