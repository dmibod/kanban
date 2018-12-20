package query

import (
	"github.com/go-chi/chi"
	"github.com/dmibod/kanban/shared/persistence"
	"github.com/dmibod/kanban/shared/services"
	"github.com/dmibod/kanban/shared/tools/db"
	"github.com/dmibod/kanban/shared/tools/logger"
)

// Env holds module dependencies
type Env struct {
	Mux    *chi.Mux
	Factory db.Factory
	Logger  logger.Logger
}

// Boot installs Query module http handlers to mux
func (e *Env) Boot() {

	repository := persistence.CreateCardRepository(e.Factory)
	service    := services.CreateCardService(e.Logger, repository)
	
	api := CreateAPI(e.Logger, service)
	api.Routes(e.Mux)

	e.Logger.Debugln("endpoints registered")
}
