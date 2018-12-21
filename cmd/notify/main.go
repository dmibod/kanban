package main

import (
	"github.com/dmibod/kanban/shared/tools/logger/console"
	"github.com/dmibod/kanban/notify"
	"github.com/dmibod/kanban/cmd/shared"
)

func main() {

	l := console.New(
		console.WithPrefix("[NOTIFY.] "), 
		console.WithDebug(true))
		
	m := mux.ConfigureMux()

  module := notify.Env{Logger: l, Mux: m }

	module.Boot()

	mux.PrintRoutes(l, m)
	mux.StartMux(m, mux.GetPortOrDefault(3001))
}
