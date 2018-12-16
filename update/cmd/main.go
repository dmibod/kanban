package main

import (
	"github.com/dmibod/kanban/tools/mux/http"
	"github.com/dmibod/kanban/tools/db/mongo"
	"github.com/dmibod/kanban/update"
)

func main() {
	m := http.New(http.WithPort(http.GetPortOrDefault(3003)))
	f := mongo.New(mongo.WithDatabase("kanban"))

	update.Boot(m, f)

  m.Start()
}
