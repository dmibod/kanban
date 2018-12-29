package main

import (
	"context"
	"time"

	"github.com/dmibod/kanban/cmd/shared"

	"github.com/dmibod/kanban/process"
)

func main() {
	c, cancel := context.WithCancel(context.Background())

	l := shared.CreateLogger("[PROCESS] ", true)

	module := process.Module{Ctx: c, Logger: l}
	module.Boot()

	shared.StartBus(c, shared.GetNameOrDefault("proc"), shared.CreateLogger("[..BUS..] ", true))

	<-shared.GetInterruptChan()

	l.Debugln("interrupt signal received!")

	cancel()

	time.Sleep(time.Second)
}
