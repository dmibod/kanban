package main

import (
	"github.com/dmibod/kanban/tools/mux/http"
	"github.com/dmibod/kanban/command"
)

func main() {
	mux := http.New(http.WithPort(3000))

	command.Boot(mux)

  mux.Start()
}