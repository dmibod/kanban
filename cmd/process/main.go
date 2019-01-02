package main

import (
	"context"
	"github.com/dmibod/kanban/shared/tools/db/mongo"
	"time"

	"github.com/dmibod/kanban/cmd/shared"

	"github.com/dmibod/kanban/process"
)

func main() {
	c, cancel := context.WithCancel(context.Background())

	l := shared.CreateLogger("[PROCESS] ", true)

	slog := shared.CreateLogger("[SESSION] ", true)
	sess := mongo.CreateSessionFactory(mongo.WithLogger(slog))
	exec := shared.CreateExecutor(sess)
	rfac := shared.CreateRepositoryFactory(exec)
	sfac := shared.CreateServiceFactory(rfac)
	cfac := mongo.CreateContextFactory(sess, slog)

	module := process.Module{ServiceFactory: sfac, ContextFactory: cfac, Context: c, Logger: l}
	module.Boot()

	shared.StartBus(c, shared.GetNameOrDefault("proc"), shared.CreateLogger("[..BUS..] ", true))

	<-shared.GetInterruptChan()

	l.Debugln("interrupt signal received!")

	cancel()

	time.Sleep(time.Second)

	shared.StopBus()

	l.Debugln("done")
}
