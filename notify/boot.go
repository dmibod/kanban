package notify

import (
	"github.com/go-chi/chi"
	"github.com/dmibod/kanban/shared/tools/logger"

	"github.com/dmibod/kanban/shared/tools/msg"
	"github.com/dmibod/kanban/shared/tools/msg/nats"
)

// Env holds module dependencies
type Env struct {
	Mux    *chi.Mux
	Logger  logger.Logger
}

// Boot installs notify module handlers to mux
func (e *Env) Boot() {

	var t msg.Transport = nats.New()

	api := CreateAPI(e.Logger, t.Receive("notification"))

	e.Mux.Route("/v1/api/notify", func(r chi.Router) {
		r.Mount("/", api.Routes())
	})

	e.Logger.Debugln("endpoints registered")
}
