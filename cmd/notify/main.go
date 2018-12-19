package main

import (
	"github.com/dmibod/kanban/shared/tools/log/logger"
	"github.com/dmibod/kanban/notify"
	"github.com/dmibod/kanban/shared/tools/mux/http"
)

func main() {

	l := logger.New(logger.WithPrefix("[NOTIFY.] "), logger.WithDebug(true))
	m := http.New(http.WithPort(http.GetPortOrDefault(3001)))

	notify.Boot(m, l)

	m.Start()
}
