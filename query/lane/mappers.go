package lane

import (
	"github.com/dmibod/kanban/shared/services/lane"
)

// ModelMapper type
type ModelMapper struct {
}

// ModelToPayload mapping
func (ModelMapper) ModelToPayload(m interface{}) interface{} {
	model := m.(*lane.Model)
	return &Lane{
		ID:     string(model.ID),
		Name:   model.Name,
		Type:   model.Type,
		Layout: model.Layout,
	}
}

// ListModelMapper type
type ListModelMapper struct {
}

// ModelToPayload mapping
func (ListModelMapper) ModelToPayload(m interface{}) interface{} {
	model := m.(*lane.ListModel)
	return &Lane{
		ID:     string(model.ID),
		Name:   model.Name,
		Type:   model.Type,
		Layout: model.Layout,
	}
}

// List mapping
func (m ListModelMapper) List(models []*lane.ListModel) []interface{} {
	list := make([]interface{}, len(models))
	for i, model := range models {
		list[i] = model
	}
	return list
}
