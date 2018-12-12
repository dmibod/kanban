package process

import (
	"context"
	"log"

	"github.com/dmibod/kanban/tools/msg"
	"github.com/dmibod/kanban/tools/msg/nats"
)

func Boot(c context.Context) {
	log.Println("Starting processor...");

	var t msg.Transport = nats.New()

	env := &Env{In: t.Receive("command"), Out: t.Send("notification")}

	go env.Handle(c)

	log.Println("Processor started!");
}
