package main

import (
	"github.com/dmibod/kanban/shared/tools/logger/console"
	"github.com/dmibod/kanban/notify"
	utils "github.com/dmibod/kanban/shared/tools/mux"
)

func main() {

	l := console.New(
		console.WithPrefix("[NOTIFY.] "), 
		console.WithDebug(true))
		
	m := utils.ConfigureMux()

  module := notify.Env{Logger: l, Mux: m }

	module.Boot()

	utils.PrintRoutes(l, m)

	utils.StartMux(m, utils.GetPortOrDefault(3001))
}
