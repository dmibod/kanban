package main

import (
	"fmt"
	"github.com/dmibod/kanban/shared/tools/logger"
	"net/http"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/go-chi/chi"
	"github.com/dmibod/kanban/shared/tools/logger/console"
	"github.com/dmibod/kanban/query"
	"github.com/dmibod/kanban/shared/persistence"
	"github.com/dmibod/kanban/shared/tools/db/mongo"
	port "github.com/dmibod/kanban/shared/tools/mux/http"
)

func main() {

	l := console.New(
		console.WithPrefix("[QUERY..] "),
		console.WithDebug(true))

	f := mongo.CreateFactory(
		mongo.WithDatabase("kanban"), 
		mongo.WithExecutor(persistence.CreateService(l)), 
		mongo.WithLogger(l))

	m := configureMux()

  module := query.Env{Logger: l, Factory: f, Mux: m }

	module.Boot()

	printRoutes(l, m)

	http.ListenAndServe(fmt.Sprintf(":%v", port.GetPortOrDefault(3002)), m)
}

func configureMux() *chi.Mux {
	router := chi.NewRouter()

	router.Use(
		render.SetContentType(render.ContentTypeJSON), // Set content-Type headers as application/json
		middleware.Logger,                             // Log API request calls
		middleware.DefaultCompress,                    // Compress results, mostly gzipping assets and json
		middleware.RedirectSlashes,                    // Redirect slashes to no slash URL versions
		middleware.Recoverer,                          // Recover from panics without crashing server
	)

	return router
}

