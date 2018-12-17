package main

import (
	"github.com/dmibod/kanban/command"
	"github.com/dmibod/kanban/shared/tools/mux/http"
)

func main() {
	mux := http.New(http.WithPort(http.GetPortOrDefault(3000)))

	command.Boot(mux)

	mux.Start()
}
