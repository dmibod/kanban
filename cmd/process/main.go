package main

import (
	"context"
	"github.com/dmibod/kanban/shared/tools/msg/nats"
	"os"
	"os/signal"
	"time"

	"github.com/dmibod/kanban/process"
	"github.com/dmibod/kanban/shared/message"
	"github.com/dmibod/kanban/shared/tools/logger/console"
)

func main() {
	l := console.New(console.WithPrefix("[PROCESS] "), console.WithDebug(true))

	ctx, cancel := context.WithCancel(context.Background())

	c := make(chan os.Signal, 1)

	signal.Notify(c, os.Interrupt)

	module := process.Module{Logger: l, Ctx: ctx, Msg: nats.CreateTransport(ctx, message.CreateService(l))}

	module.Boot()

	<-c

	l.Debugln("Interrupt signal received!")

	cancel()

	time.Sleep(time.Second)
}
