package main

import (
	"github.com/dmibod/kanban/cmd/shared"
	"github.com/dmibod/kanban/command"
	"github.com/dmibod/kanban/shared/tools/logger/console"
)

func main() {

	l := console.New(
		console.WithPrefix("[COMMAND] "),
		console.WithDebug(true))

	m := mux.ConfigureMux()

	module := command.Module{Logger: l, Mux: m}

	module.Boot()

	mux.StartMux(m, mux.GetPortOrDefault(3000), l)
}
