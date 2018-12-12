package process

import (
	"context"

	"github.com/dmibod/kanban/tools/msg"
	"github.com/dmibod/kanban/tools/msg/nats"
)

func Boot(c context.Context) {
	var t msg.Transport = nats.New()

	env := &Env{In: t.Receive("command"), Out: t.Send("notification")}

	go env.Handle(c)
}
