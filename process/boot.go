package process

import (
	"context"

	"github.com/dmibod/kanban/shared/tools/msg"
	"github.com/dmibod/kanban/shared/tools/msg/nats"
	"github.com/dmibod/kanban/shared/tools/log/logger"
)

func Boot(c context.Context) {
	l := logger.New(logger.WithPrefix("[PROCESS] "), logger.WithDebug(true))

	l.Debugln("starting...");

	var t msg.Transport = nats.New()

	env := &Env{ Logger: l, In: t.Receive("command"), Out: t.Send("notification")}

	go env.Handle(c)

	l.Debugln("started!");
}
