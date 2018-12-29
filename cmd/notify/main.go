package main

import (
	"time"
	"context"
	"github.com/dmibod/kanban/cmd/shared"
	"github.com/dmibod/kanban/notify"
)

func main() {
	c, cancel := context.WithCancel(context.Background())

	l := shared.CreateLogger("[NOTIFY.] ", true)
	m := shared.ConfigureMux()

	module := notify.Module{Mux: m, Logger: l}
	module.Boot()

	shared.StartBus(c, shared.GetNameOrDefault("notify"), shared.CreateLogger("[..BUS..] ", true))
	shared.StartMux(m, shared.GetPortOrDefault(8001), shared.CreateLogger("[..MUX..] ", true))

	<-shared.GetInterruptChan()

	l.Debugln("interrupt signal received!")

	cancel()

	time.Sleep(time.Second)
}
