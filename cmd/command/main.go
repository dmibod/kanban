package main

import (
	"time"
	"context"
	"github.com/dmibod/kanban/cmd/shared"
	"github.com/dmibod/kanban/command"
)

func main() {
	c, cancel := context.WithCancel(context.Background())

	l := shared.CreateLogger("[COMMAND] ", true)
	m := shared.ConfigureMux()

	module := command.Module{Mux: m, Logger: l}
	module.Boot()

	shared.StartBus(c, shared.GetNameOrDefault("cmd"), shared.CreateLogger("[..BUS..] ", true))
	shared.StartMux(m, shared.GetPortOrDefault(8000), shared.CreateLogger("[..MUX..] ", true))
	
	<-shared.GetInterruptChan()

	l.Debugln("interrupt signal received!")

	cancel()

	time.Sleep(time.Second)
}
