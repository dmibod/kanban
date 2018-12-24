package main

import (
	"context"
	"github.com/dmibod/kanban/process"
	"github.com/dmibod/kanban/shared/tools/logger/console"
	"os"
	"os/signal"
	"time"
)

func main() {
	l := console.New(console.WithPrefix("[PROCESS] "), console.WithDebug(true))

	ctx, cancel := context.WithCancel(context.Background())

	c := make(chan os.Signal, 1)

	signal.Notify(c, os.Interrupt)

	process.Boot(ctx, l)

	<-c

	l.Debugln("Interrupt signal received!")

	cancel()

	time.Sleep(time.Second)
}
