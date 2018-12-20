package main

import (
	"github.com/dmibod/kanban/shared/tools/logger/console"
	"github.com/dmibod/kanban/command"
	"github.com/dmibod/kanban/shared/tools/mux/http"
)

func main() {

	l := console.New(console.WithPrefix("[COMMAND] "), console.WithDebug(true))
  m := http.New(http.WithPort(http.GetPortOrDefault(3000)))

	command.Boot(m, l)

	m.Start()
}
