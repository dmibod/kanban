package board

import (
	"context"
	"net/http"

	laneapi "github.com/dmibod/kanban/query/lane"

	"github.com/dmibod/kanban/shared/tools/mux"
	"github.com/go-chi/render"

	"github.com/dmibod/kanban/shared/handlers"
	"github.com/dmibod/kanban/shared/kernel"
	"github.com/dmibod/kanban/shared/services/board"
	"github.com/dmibod/kanban/shared/services/lane"
	"github.com/go-chi/chi"

	"github.com/dmibod/kanban/shared/tools/logger"
)

// API dependencies
type API struct {
	boardService board.Service
	laneService  lane.Service
	logger.Logger
}

// CreateAPI creates API
func CreateAPI(boardService board.Service, laneService lane.Service, l logger.Logger) *API {
	return &API{
		boardService: boardService,
		laneService:  laneService,
		Logger:       l,
	}
}

// Routes install handlers
func (a *API) Routes(router chi.Router) {
	router.Get("/", a.All)
	router.Get("/{BOARDID}", a.Get)
	router.Get("/{BOARDID}/lane", a.GetLanes)
}

// All - gets all boards
func (a *API) All(w http.ResponseWriter, r *http.Request) {
	owner := r.URL.Query().Get("owner")
	if models, err := a.boardService.GetByOwner(r.Context(), owner); err == nil {
		mapper := ListModelMapper{}
		render.JSON(w, r, mapper.ModelsToPayload(models))
	} else {
		a.Errorln(err)
		mux.RenderError(w, http.StatusInternalServerError)
	}
}

// Get by id
func (a *API) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "BOARDID")
	op := handlers.Get(id, a, &ListModelMapper{}, a.Logger)
	handlers.Handle(w, r, op)
}

// GetLanes by lane
func (a *API) GetLanes(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "BOARDID")
	if cards, err := a.laneService.GetByBoardID(r.Context(), kernel.ID(id)); err == nil {
		mapper := laneapi.ListModelMapper{}
		render.JSON(w, r, mapper.ModelsToPayload(cards))
	} else {
		a.Errorln(err)
		mux.RenderError(w, http.StatusInternalServerError)
	}
}

// GetByID implements handlers.GetService
func (a *API) GetByID(ctx context.Context, id kernel.ID) (interface{}, error) {
	return a.boardService.GetByID(ctx, id)
}
