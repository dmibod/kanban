package main

import (
	"expvar"
	"github.com/dmibod/kanban/cmd/shared"
	"github.com/dmibod/kanban/query"
	"github.com/dmibod/kanban/shared/persistence"
	"github.com/dmibod/kanban/shared/services"
	"github.com/dmibod/kanban/shared/tools/db/mongo"
	"github.com/dmibod/kanban/shared/tools/logger"
	"github.com/dmibod/kanban/shared/tools/logger/console"
	"net/http"
	"net/http/pprof"
)

func main() {

	f := mongo.CreateFactory(
		"kanban",
		persistence.CreateService(createLogger("[BRK.MGO]", true)),
		createLogger("[MONGO..]", true))

	m := mux.ConfigureMux()

	exph := expvar.Handler()
	m.Get("/vars", func(w http.ResponseWriter, r *http.Request) { exph.ServeHTTP(w, r) })
	m.Get("/prof", pprof.Index)

	module := query.Module{
		Logger:  createLogger("[QUERY..]", true),
		Factory: services.CreateFactory(createLogger("[SERVICE]", true), f),
		Mux:     m,
	}

	module.Boot(true)

	mux.StartMux(m, mux.GetPortOrDefault(8002), createLogger("[..MUX..]", true))
}

func createLogger(prefix string, debug bool) logger.Logger {
	return console.New(console.WithPrefix(prefix), console.WithDebug(debug))
}
