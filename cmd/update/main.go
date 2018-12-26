package main

import (
	"github.com/dmibod/kanban/cmd/shared"
	"github.com/dmibod/kanban/shared/persistence"
	"github.com/dmibod/kanban/shared/services"
	"github.com/dmibod/kanban/shared/tools/logger"
	"github.com/dmibod/kanban/shared/tools/logger/console"
	"github.com/dmibod/kanban/update"
)

func main() {

	f := persistence.CreateFactory(
		persistence.CreateService(createLogger("[BRK.MGO]", true)),
		createLogger("[MONGO..]", true))

	m := mux.ConfigureMux()

	module := update.Module{
		Logger:  createLogger("[UPDATE.]", true),
		Factory: services.CreateFactory(createLogger("[SERVICE]", true), f),
		Mux:     m,
	}

	module.Boot(true)

	mux.StartMux(m, mux.GetPortOrDefault(8003), createLogger("[..MUX..]", true))
}

func createLogger(prefix string, debug bool) logger.Logger {
	return console.New(console.WithPrefix(prefix), console.WithDebug(debug))
}
