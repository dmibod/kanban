package nats

import (
	"context"

	"github.com/nats-io/go-nats"
)

type sender struct {
	s   string
	e   OperationExecutor
	ctx *OperationContext
}

func (s *sender) Send(msg []byte) error {
	return s.e.Execute(s.ctx, func(ctx context.Context, conn *nats.Conn) error {
		err := conn.Publish(s.s, msg)
		if err == nil {
			err = conn.Flush()
		}
		return err
	})
}
