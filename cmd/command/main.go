package main

import (
	"github.com/dmibod/kanban/shared/tools/logger/console"
	"github.com/dmibod/kanban/command"
	"github.com/dmibod/kanban/cmd/shared"
)

func main() {

	l := console.New(
		console.WithPrefix("[COMMAND] "), 
		console.WithDebug(true))

	m := mux.ConfigureMux()

  module := command.Env{Logger: l, Mux: m }

	module.Boot()

	mux.PrintRoutes(l, m)
	mux.StartMux(m, mux.GetPortOrDefault(3000))
}
