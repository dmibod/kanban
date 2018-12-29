package main

import (
	"github.com/dmibod/kanban/cmd/shared"
	"github.com/dmibod/kanban/update"
)

func main() {

	m := shared.ConfigureMux()

	module := update.Module{
		Logger:  shared.CreateLogger("[UPDATE.]", true),
		Factory: shared.CreateServiceFactory(),
		Mux:     m,
	}

	module.Boot(true)

	shared.StartMux(m, shared.GetPortOrDefault(8003), shared.CreateLogger("[..MUX..]", true))
}
