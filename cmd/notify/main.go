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

  module := notify.Module{Logger: l, Mux: m }

	module.Boot()

	mux.StartMux(m, mux.GetPortOrDefault(3001), l)
}
