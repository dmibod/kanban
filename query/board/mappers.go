package board

import (
	"github.com/dmibod/kanban/shared/services/board"
)

// ModelMapper type
type ModelMapper struct {
}

// ModelToPayload mapping
func (ModelMapper) ModelToPayload(m interface{}) interface{} {
	model := m.(*board.Model)
	return &Model{
		ID:          string(model.ID),
		Name:        model.Name,
		Description: model.Description,
		Layout:      model.Layout,
		Owner:       model.Owner,
		Shared:      model.Shared,
	}
}

// ListModelMapper type
type ListModelMapper struct {
}

// ModelToPayload mapping
func (ListModelMapper) ModelToPayload(m interface{}) interface{} {
	model := m.(*board.ListModel)
	return &ListModel{
		ID:          string(model.ID),
		Name:        model.Name,
		Description: model.Description,
		Owner:       model.Owner,
		Shared:      model.Shared,
	}
}

// List mapping
func (m ListModelMapper) List(models []*board.ListModel) []interface{} {
	list := make([]interface{}, len(models))
	for i, model := range models {
		list[i] = model
	}
	return list
}
