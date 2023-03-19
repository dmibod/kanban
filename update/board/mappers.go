package board

import (
	"github.com/dmibod/kanban/shared/services/board"
)

type boardMapper struct {
}

// PayloadToModel mapping
func (boardMapper) PayloadToModel(p interface{}) interface{} {
	payload := p.(*Board)

	return &board.CreateModel{
		Name:   payload.Name,
		Description: payload.Description,
		Layout: payload.Layout,
		Owner:  payload.Owner,
	}
}

func (boardMapper) ModelToPayload(m interface{}) interface{} {
	model := m.(*board.Model)

	return &Board{
		ID:     string(model.ID),
		Name:   model.Name,
		Description: model.Description,
		Layout: model.Layout,
		Owner:  model.Owner,
		Shared: model.Shared,
	}
}
