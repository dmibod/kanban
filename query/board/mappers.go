package board

import (
	"github.com/dmibod/kanban/shared/services/board"
)

// ListModelMapper type
type ListModelMapper struct {
}

// ModelToPayload mapping
func (ListModelMapper) ModelToPayload(m interface{}) interface{} {
	if model, ok := m.(*board.ListModel); ok {
		return &ListModel{
			ID:          string(model.ID),
			Name:        model.Name,
			Description: model.Description,
			Owner:       model.Owner,
			Shared:      model.Shared,
		}
	}

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

// ModelsToPayload mapping
func (m ListModelMapper) ModelsToPayload(models []*board.ListModel) []interface{} {
	items := []interface{}{}
	for _, model := range models {
		items = append(items, m.ModelToPayload(model))
	}
	return items
}
