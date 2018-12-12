package main

import (
	"github.com/dmibod/kanban/notify"
	"github.com/dmibod/kanban/tools/mux/http"
)

func main() {
	mux := http.New(http.WithPort(3001))

	notify.Boot(mux)

  mux.Start()
}