package process

import (
	"context"

	"github.com/dmibod/kanban/tools/msg"
	"github.com/dmibod/kanban/tools/msg/nats"
	"github.com/dmibod/kanban/tools/log/logger"
)

func Boot(c context.Context) {
	l := logger.New(logger.WithPrefix("[PROCESS] "), logger.WithDebug(true))

	l.Infoln("starting...");

	var t msg.Transport = nats.New()

	env := &Env{In: t.Receive("command"), Out: t.Send("notification")}

	go env.Handle(c)

	l.Infoln("started!");
}
