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

	sess := shared.CreateSessionFactory()
	glob := shared.CreateSessionProvider(sess)
	prov := shared.CreateCopySessionProvider(glob)
	exec := shared.CreateOperationExecutor(prov)
	cfac := shared.CreateContextFactory(prov)
	rfac := shared.CreateRepositoryFactory(exec)
	sfac := shared.CreateServiceFactory(rfac)

	ctx, err := cfac.Context(c)
	if err != nil {
		l.Errorln(err)
		cancel()
		return
	}

	module := process.Module{ServiceFactory: sfac, Context: ctx, Logger: l}
	module.Boot()

	shared.StartBus(c, shared.GetNameOrDefault("proc"), shared.CreateLogger("[..BUS..] ", true))

	<-shared.GetInterruptChan()

	l.Debugln("interrupt signal received!")

	glob.Provide().Close(false)

	cancel()

	time.Sleep(time.Second)

	shared.StopBus()

	l.Debugln("done")
}
