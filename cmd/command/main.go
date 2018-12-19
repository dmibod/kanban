package main

import (
	"github.com/dmibod/kanban/shared/tools/log/logger"
	"github.com/dmibod/kanban/command"
	"github.com/dmibod/kanban/shared/tools/mux/http"
)

func main() {

	l := logger.New(logger.WithPrefix("[COMMAND] "), logger.WithDebug(true))
  m := http.New(http.WithPort(http.GetPortOrDefault(3000)))

	command.Boot(m, l)

	m.Start()
}
