package nats

import (
	"context"
)

type sender struct {
	s   string
	e   OperationExecutor
	ctx *OperationContext
}

func createSender(s string, c *OperationContext, e OperationExecutor) *sender {
	return &sender{
		e:   e,
		s:   s,
		ctx: c,
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
