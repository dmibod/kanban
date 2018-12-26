package main

import (
	"context"
	"github.com/dmibod/kanban/shared/tools/logger"
	"os"
	"os/signal"
	"time"

	"github.com/dmibod/kanban/process"
	"github.com/dmibod/kanban/shared/message"
	"github.com/dmibod/kanban/shared/tools/logger/console"
)

func main() {

	ctx, cancel := context.WithCancel(context.Background())

	c := make(chan os.Signal, 1)

	signal.Notify(c, os.Interrupt)

	s := message.CreateService("PROCESS", createLogger("[BRK.NAT] ", true))

	l := createLogger("[PROCESS] ", true)

	t := message.CreateTransport(ctx, s, l)

	module := process.Module{Ctx: ctx, Msg: t, Logger: l}

	module.Boot()

	<-c

	l.Debugln("interrupt signal received!")

	cancel()

	time.Sleep(time.Second)
}

func createLogger(prefix string, debug bool) logger.Logger {
	return console.New(console.WithPrefix(prefix), console.WithDebug(debug))
}
