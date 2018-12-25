package mux

import (
	"fmt"
	"github.com/dmibod/kanban/shared/tools/logger"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	//"github.com/go-chi/render"
	"log"
	"net/http"
	"os"
	"strconv"
)

const muxPortEnvVar = "MUX_PORT"

// GetPortOrDefault gets port from environment variable or fallbacks to default one
func GetPortOrDefault(defPort int) int {
	env := os.Getenv(muxPortEnvVar)

	port, err := strconv.Atoi(env)
	if err != nil {
		return defPort
	}

	return port
}

// ConfigureMux configures default mux
func ConfigureMux() *chi.Mux {
	router := chi.NewRouter()

	router.Use(
		//render.SetContentType(render.ContentTypeJSON), // Set content-Type headers as application/json
		//middleware.Logger,                             // Log API request calls
		middleware.RequestLogger(&middleware.DefaultLogFormatter{Logger: log.New(os.Stdout, "", log.LstdFlags), NoColor: true}),
		middleware.DefaultCompress, // Compress results, mostly gzipping assets and json
		middleware.RedirectSlashes, // Redirect slashes to no slash URL versions
		middleware.Recoverer,       // Recover from panics without crashing server
	)

	return router
}

// StartMux starts mux
func StartMux(m *chi.Mux, port int, l logger.Logger) {
	walkFunc := func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		l.Debugf("%s %s\n", method, route) // Walk and print out all routes
		return nil
	}

	if err := chi.Walk(m, walkFunc); err != nil {
		l.Errorf("Logging err: %s\n", err.Error()) // panic if there is an error
		panic(err)
	}

	if err := http.ListenAndServe(fmt.Sprintf(":%v", port), m); err != nil {
		l.Errorf("Mux err: %s\n", err.Error()) // panic if there is an error
		panic(err)
	}
}
