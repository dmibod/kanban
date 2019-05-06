package lane

import (
	"github.com/dmibod/kanban/shared/services/lane"
)

type laneCreateMapper struct {
}

// PayloadToModel mapping
func (laneCreateMapper) PayloadToModel(p interface{}) interface{} {
	payload := p.(*Lane)
	return &lane.CreateModel{
		Type:        payload.Type,
		Name:        payload.Name,
		Description: payload.Description,
		Layout:      payload.Layout,
	}
}

// ModelToPayload mapping
func (laneCreateMapper) ModelToPayload(m interface{}) interface{} {
	model := m.(*lane.Model)
	return &Lane{
		ID:          string(model.ID),
		Type:        model.Type,
		Name:        model.Name,
		Description: model.Description,
		Layout:      model.Layout,
	}
}
