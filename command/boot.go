package command

import (
	"github.com/go-chi/chi"
	"github.com/dmibod/kanban/shared/tools/logger"
	"github.com/dmibod/kanban/shared/tools/msg/nats"
	"github.com/dmibod/kanban/shared/tools/msg"
)

// Env holds module dependencies
type Env struct {
	Mux    *chi.Mux
	Logger  logger.Logger
}

// Boot installs command module handlers to mux
func (e *Env) Boot(){

	var t msg.Transport = nats.New()

	api := CreateAPI(e.Logger, t.Send("command"))

	e.Mux.Route("/v1", func(r chi.Router) {
		r.Mount("/api/commands", api.Routes())
	})

	e.Logger.Debugln("endpoints registered")
}