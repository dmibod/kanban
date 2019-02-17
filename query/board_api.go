package query

import (
	"context"
	"net/http"

	"github.com/dmibod/kanban/shared/tools/mux"
	"github.com/go-chi/render"

	"github.com/dmibod/kanban/shared/handlers"
	"github.com/dmibod/kanban/shared/kernel"
	"github.com/dmibod/kanban/shared/services"
	"github.com/go-chi/chi"

	"github.com/dmibod/kanban/shared/tools/logger"
)

// BoardAPI dependencies
type BoardAPI struct {
	services.BoardService
	services.LaneService
	logger.Logger
}

// CreateBoardAPI creates API
func CreateBoardAPI(boardService services.BoardService, laneService services.LaneService, l logger.Logger) *BoardAPI {
	return &BoardAPI{
		BoardService: boardService,
		LaneService:  laneService,
		Logger:       l,
	}
}

// Routes install handlers
func (a *BoardAPI) Routes(router chi.Router) {
	router.Get("/", a.All)
	router.Get("/{BOARDID}", a.Get)
	router.Get("/{BOARDID}/lane", a.GetLanes)
}

// All - gets all boards
func (a *BoardAPI) All(w http.ResponseWriter, r *http.Request) {
	owner := r.URL.Query().Get("owner")
	models, err := a.BoardService.GetByOwner(r.Context(), owner)

	if err != nil {
		a.Errorln(err)
		mux.RenderError(w, http.StatusInternalServerError)
		return
	}

	mapper := BoardGetMapper{}
	render.JSON(w, r, mapper.ModelsToPayload(models))
}

// Get by id
func (a *BoardAPI) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "BOARDID")
	op := handlers.Get(id, a, &BoardGetMapper{}, a.Logger)
	handlers.Handle(w, r, op)
}

// GetLanes by lane
func (a *BoardAPI) GetLanes(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "BOARDID")
	cards, err := a.LaneService.GetByBoardID(r.Context(), kernel.ID(id))
	if err != nil {
		a.Errorln(err)
		mux.RenderError(w, http.StatusInternalServerError)
		return
	}
	mapper := LaneGetMapper{}
	render.JSON(w, r, mapper.ModelsToPayload(cards))
}

// GetByID implements handlers.GetService
func (a *BoardAPI) GetByID(ctx context.Context, id kernel.ID) (interface{}, error) {
	return a.BoardService.GetByID(ctx, id)
}
