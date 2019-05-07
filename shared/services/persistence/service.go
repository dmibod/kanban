package persistence

import (
	"context"
	"github.com/dmibod/kanban/shared/domain/board"
	"github.com/dmibod/kanban/shared/domain/card"
	"github.com/dmibod/kanban/shared/domain/event"
	"github.com/dmibod/kanban/shared/domain/lane"
	b "github.com/dmibod/kanban/shared/persistence/board"
	c "github.com/dmibod/kanban/shared/persistence/card"
	l "github.com/dmibod/kanban/shared/persistence/lane"
	"github.com/dmibod/kanban/shared/persistence/models"
	"github.com/dmibod/kanban/shared/tools/logger"
	"gopkg.in/mgo.v2/bson"
)

type service struct {
	logger.Logger
	boardRepository b.Repository
	laneRepository  l.Repository
	cardRepository  c.Repository
}

// CreateService instance
func CreateService(b b.Repository, l l.Repository, c c.Repository, log logger.Logger) Service {
	return &service{
		boardRepository: b,
		laneRepository:  l,
		cardRepository:  c,
		Logger:          log,
	}
}

// Listen doamin events
func (s *service) Listen(bus event.Bus) {
	if bus != nil {
		bus.Listen(s)
	}
}

// Handle doamin event
func (s *service) Handle(ctx context.Context, event interface{}) error {
	if event == nil {
		return nil
	}

	s.Debugf("domain event: %+v\n", event)

	if ok, err := s.handleBoardEvent(ctx, event); err != nil {
		s.Error(err)
		return err
	} else if ok {
		return nil
	}

	if ok, err := s.handleLaneEvent(ctx, event); err != nil {
		s.Error(err)
		return err
	} else if ok {
		return nil
	}

	if _, err := s.handleCardEvent(ctx, event); err != nil {
		s.Error(err)
		return err
	}

	return nil
}

func (s *service) handleBoardEvent(ctx context.Context, event interface{}) (bool, error) {
	var err error

	switch e := event.(type) {
	case board.CreatedEvent:
		err = s.boardRepository.Create(ctx, s.mapBoard(&e.Entity))
	case board.DeletedEvent:
		err = s.boardRepository.Remove(ctx, e.Entity.ID.String())
	case board.NameChangedEvent:
		err = s.boardRepository.Update(ctx, e.ID.String(), "name", e.NewValue)
	case board.DescriptionChangedEvent:
		err = s.boardRepository.Update(ctx, e.ID.String(), "description", e.NewValue)
	case board.LayoutChangedEvent:
		err = s.boardRepository.Update(ctx, e.ID.String(), "layout", e.NewValue)
	case board.SharedChangedEvent:
		err = s.boardRepository.Update(ctx, e.ID.String(), "shared", e.NewValue)
	case board.ChildAppendedEvent:
		err = s.boardRepository.Attach(ctx, e.ID.SetID.String(), e.ID.ID.String())
	case board.ChildRemovedEvent:
		err = s.boardRepository.Detach(ctx, e.ID.SetID.String(), e.ID.ID.String())
	default:
		return false, nil
	}

	return true, err
}

func (s *service) handleLaneEvent(ctx context.Context, event interface{}) (bool, error) {
	var err error

	switch e := event.(type) {
	case lane.CreatedEvent:
		err = s.laneRepository.Create(ctx, e.ID.SetID.String(), s.mapLane(&e.Entity))
	case lane.DeletedEvent:
		err = s.laneRepository.Remove(ctx, e.ID.SetID.String(), e.ID.ID.String())
	case lane.NameChangedEvent:
		err = s.laneRepository.Update(ctx, e.ID.SetID.String(), e.ID.ID.String(), "name", e.NewValue)
	case lane.DescriptionChangedEvent:
		err = s.laneRepository.Update(ctx, e.ID.SetID.String(), e.ID.ID.String(), "description", e.NewValue)
	case lane.LayoutChangedEvent:
		err = s.laneRepository.Update(ctx, e.ID.SetID.String(), e.ID.ID.String(), "layout", e.NewValue)
	case lane.ChildAppendedEvent:
		err = s.laneRepository.Attach(ctx, e.ID.SetID.String(), e.ID.ID.String(), e.ChildID.String())
	case lane.ChildRemovedEvent:
		err = s.laneRepository.Detach(ctx, e.ID.SetID.String(), e.ID.ID.String(), e.ChildID.String())
	default:
		return false, nil
	}

	return true, err
}

func (s *service) handleCardEvent(ctx context.Context, event interface{}) (bool, error) {
	var err error

	switch e := event.(type) {
	case card.CreatedEvent:
		err = s.cardRepository.Create(ctx, e.ID.SetID.String(), s.mapCard(&e.Entity))
	case card.DeletedEvent:
		err = s.cardRepository.Remove(ctx, e.ID.SetID.String(), e.ID.ID.String())
	case card.NameChangedEvent:
		err = s.cardRepository.Update(ctx, e.ID.SetID.String(), e.ID.ID.String(), "name", e.NewValue)
	case card.DescriptionChangedEvent:
		err = s.cardRepository.Update(ctx, e.ID.SetID.String(), e.ID.ID.String(), "description", e.NewValue)
	default:
		return false, nil
	}

	return true, err
}

func (*service) mapBoard(entity *board.Entity) *models.Board {
	return &models.Board{
		ID:          bson.ObjectIdHex(entity.ID.String()),
		Owner:       entity.Owner,
		Name:        entity.Name,
		Description: entity.Description,
		Layout:      entity.Layout,
		Shared:      entity.Shared,
		Children:    []bson.ObjectId{},
		Lanes:       []models.Lane{},
		Cards:       []models.Card{},
	}
}

func (*service) mapLane(entity *lane.Entity) *models.Lane {
	return &models.Lane{
		ID:          bson.ObjectIdHex(entity.ID.ID.String()),
		Kind:        entity.Kind,
		Name:        entity.Name,
		Description: entity.Description,
		Layout:      entity.Layout,
		Children:    []bson.ObjectId{},
	}
}

func (*service) mapCard(entity *card.Entity) *models.Card {
	return &models.Card{
		ID:          bson.ObjectIdHex(entity.ID.ID.String()),
		Name:        entity.Name,
		Description: entity.Description,
	}
}
