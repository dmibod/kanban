package board

import (
	"github.com/dmibod/kanban/shared/services/board"
)

// ModelMapper type
type ModelMapper struct {
}

func cardToPayload(model board.CardModel) CardModel {
	return CardModel{
		ID:          string(model.ID),
		Name:        model.Name,
		Description: model.Description,
	}
}

func laneToPayload(model board.LaneModel) LaneModel {
	lanes := make([]LaneModel, len(model.Lanes))
	for i, lane := range model.Lanes {
		lanes[i] = laneToPayload(lane)
	}

	cards := make([]CardModel, len(model.Cards))
	for i, card := range model.Cards {
		cards[i] = cardToPayload(card)
	}

	return LaneModel{
		ID:          string(model.ID),
		Name:        model.Name,
		Type:        model.Type,
		Layout:      model.Layout,
		Description: model.Description,
		Lanes:       lanes,
		Cards:       cards,
	}
}

// ModelToPayload mapping
func (ModelMapper) ModelToPayload(m interface{}) interface{} {
	model := m.(*board.Model)

	lanes := make([]LaneModel, len(model.Lanes))
	for i, lane := range model.Lanes {
		lanes[i] = laneToPayload(lane)
	}

	return &Model{
		ID:          string(model.ID),
		Name:        model.Name,
		Description: model.Description,
		Layout:      model.Layout,
		Owner:       model.Owner,
		Shared:      model.Shared,
		Lanes:       lanes,
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
