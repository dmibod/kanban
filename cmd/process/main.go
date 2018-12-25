package main

import (
	"context"
	"github.com/dmibod/kanban/shared/tools/logger"
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

	module := process.Module{Logger: l, Ctx: ctx, Msg: nats.CreateTransport(ctx, message.CreateService(createLogger("[BRK.NAT] ", true)), l)}

	module.Boot()

	<-c

	l.Debugln("interrupt signal received!")

	cancel()

	time.Sleep(time.Second)
}

func createLogger(prefix string, debug bool) logger.Logger {
	return console.New(console.WithPrefix(prefix), console.WithDebug(debug))
}
