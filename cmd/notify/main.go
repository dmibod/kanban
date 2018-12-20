package main

import (
	"github.com/dmibod/kanban/shared/tools/logger/console"
	"github.com/dmibod/kanban/notify"
	"github.com/dmibod/kanban/shared/tools/mux/http"
)

func main() {

	l := console.New(console.WithPrefix("[NOTIFY.] "), console.WithDebug(true))
	m := http.New(http.WithPort(http.GetPortOrDefault(3001)))

	notify.Boot(m, l)

	m.Start()
}
