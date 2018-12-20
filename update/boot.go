package update

import (
	"github.com/go-chi/chi"
	"github.com/dmibod/kanban/shared/services"
	"github.com/dmibod/kanban/shared/persistence"
	"github.com/dmibod/kanban/shared/tools/db"
	"github.com/dmibod/kanban/shared/tools/logger"
)

// Env holds module dependencies
type Env struct {
	Mux    *chi.Mux
	Factory db.Factory
	Logger  logger.Logger
}

// Boot - adds update module handlers to mux
func (e *Env) Boot() {

	repository := persistence.CreateCardRepository(e.Factory)
	service    := services.CreateCardService(e.Logger, repository)

	api := CreateAPI(e.Logger, service)

	e.Mux.Route("/v1", func(r chi.Router) {
		r.Mount("/api/card", api.Routes())
	})

	e.Logger.Debugln("endpoints registered")
}
