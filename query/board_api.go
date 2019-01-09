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

// Board api
type Board struct {
	ID     string `json:"id,omitempty"`
	Name   string `json:"name,omitempty"`
	Layout string `json:"layout,omitempty"`
	Owner  string `json:"owner,omitempty"`
	Shared bool   `json:"shared,omitempty"`
}

// BoardAPI dependencies
type BoardAPI struct {
	boardService services.BoardService
	laneService  services.LaneService
	logger.Logger
}

// CreateBoardAPI creates API
func CreateBoardAPI(boardService services.BoardService, laneService services.LaneService, l logger.Logger) *BoardAPI {
	return &BoardAPI{
		boardService: boardService,
		laneService:  laneService,
		Logger:       l,
	}
}

// Routes install handlers
func (a *BoardAPI) Routes(router chi.Router) {
	router.Get("/", a.All)
	router.Get("/{BOARDID}", a.Get)
	router.Get("/{BOARDID}/lane", a.GetLanes)
}

// Get by id
func (a *BoardAPI) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "BOARDID")
	op := handlers.Get(id, a, &BoardGetMapper{}, a.Logger)
	handlers.Handle(w, r, op)
}

// All boards
func (a *BoardAPI) All(w http.ResponseWriter, r *http.Request) {
	models, err := a.boardService.GetByOwner(r.Context(), r.URL.Query().Get("owner"))

	if err != nil {
		a.Errorln(err)
		mux.RenderError(w, http.StatusInternalServerError)
		return
	}

	mapper := BoardGetMapper{}
	render.JSON(w, r, mapper.ModelsToPayload(models))
}

// GetLanes by lane
func (a *BoardAPI) GetLanes(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "BOARDID")
	cards, err := a.laneService.GetByBoardID(r.Context(), kernel.Id(id))
	if err != nil {
		a.Errorln(err)
		mux.RenderError(w, http.StatusInternalServerError)
		return
	}
	mapper := LaneGetMapper{}
	render.JSON(w, r, mapper.ModelsToPayload(cards))
}

// GetByID implements handlers.GetService
func (a *BoardAPI) GetByID(ctx context.Context, id kernel.Id) (interface{}, error) {
	return a.boardService.GetByID(ctx, id)
}

// BoardGetMapper mapper
type BoardGetMapper struct {
}

// ModelToPayload mapping
func (BoardGetMapper) ModelToPayload(m interface{}) interface{} {
	model := m.(*services.BoardModel)
	return &Board{
		ID:     string(model.ID),
		Name:   model.Name,
		Layout: model.Layout,
		Owner:  model.Owner,
		Shared: model.Shared,
	}
}

// ModelsToPayload mapping
func (m BoardGetMapper) ModelsToPayload(models []*services.BoardModel) []interface{} {
	items := []interface{}{}
	for _, model := range models {
		items = append(items, m.ModelToPayload(model))
	}
	return items
}
