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

	e, _ := shared.CreateDatabaseServices()
	rfac := shared.CreateRepositoryFactory(e)
	sfac := shared.CreateServiceFactory(rfac)

	module := process.Module{Factory: sfac, Context: c, Logger: l}
	module.Boot()

	shared.StartBus(c, shared.GetNameOrDefault("proc"), shared.CreateLogger("[..BUS..] ", true))

	<-shared.GetInterruptChan()

	l.Debugln("interrupt signal received!")

	cancel()

	time.Sleep(time.Second)

	shared.StopBus()

	l.Debugln("done")
}
