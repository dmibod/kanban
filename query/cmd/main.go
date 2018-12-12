package main

import (
	"github.com/dmibod/kanban/query"
	"github.com/dmibod/kanban/tools/mux/http"
)

func main() {
	mux := http.New(http.WithPort(3002))

	query.Boot(mux)

	mux.Start()
}
