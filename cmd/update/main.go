package main

import (
	"github.com/dmibod/kanban/cmd/shared"
	"github.com/dmibod/kanban/shared/persistence"
	"github.com/dmibod/kanban/shared/services"
	"github.com/dmibod/kanban/update"
)

func main() {

	f := persistence.CreateFactory(
		persistence.CreateService(shared.CreateLogger("[BRK.MGO]", true)),
		shared.CreateLogger("[MONGO..]", true))

	m := shared.ConfigureMux()

	module := update.Module{
		Logger:  shared.CreateLogger("[UPDATE.]", true),
		Factory: services.CreateFactory(shared.CreateLogger("[SERVICE]", true), f),
		Mux:     m,
	}

	module.Boot(true)

	shared.StartMux(m, shared.GetPortOrDefault(8003), shared.CreateLogger("[..MUX..]", true))
}
