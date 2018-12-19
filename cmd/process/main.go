package main

import (
	"github.com/dmibod/kanban/shared/tools/log/logger"
	"time"
	"os/signal"
	"os"
	"github.com/dmibod/kanban/process"
	"context"
)

func main() {
	l := logger.New(logger.WithPrefix("[PROCESS] "), logger.WithDebug(true))

	ctx, cancel := context.WithCancel(context.Background())

	c := make(chan os.Signal, 1)

	signal.Notify(c, os.Interrupt)

	process.Boot(ctx, l)

	<-c

	l.Debugln("Interrupt signal received!");

	cancel()

	time.Sleep(time.Second)
}
