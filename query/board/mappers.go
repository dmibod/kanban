package board

import (
	"github.com/dmibod/kanban/shared/services/board"
)

// ListModelMapper type
type ListModelMapper struct {
}

// ModelToPayload mapping
func (ListModelMapper) ModelToPayload(m interface{}) interface{} {
	model := m.(*board.ListModel)
	return &ListModel{
		ID:     string(model.ID),
		Name:   model.Name,
		Owner:  model.Owner,
		Shared: model.Shared,
	}
}

// ModelsToPayload mapping
func (m ListModelMapper) ModelsToPayload(models []*board.ListModel) []interface{} {
	items := []interface{}{}
	for _, model := range models {
		items = append(items, m.ModelToPayload(model))
	}
	return items
}
