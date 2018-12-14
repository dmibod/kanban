package main

import (
	"github.com/dmibod/kanban/tools/mux/http"
	"github.com/dmibod/kanban/update"
)

func main() {
	mux := http.New(http.WithPort(http.GetPortOrDefault(3003)))

	update.Boot(mux)

  mux.Start()
}
