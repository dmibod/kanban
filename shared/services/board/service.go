package board

import (
	"context"

	"github.com/dmibod/kanban/shared/domain/event"

	tx "github.com/dmibod/kanban/shared/services/event"

	"github.com/dmibod/kanban/shared/domain/board"

	"gopkg.in/mgo.v2/bson"

	"github.com/dmibod/kanban/shared/kernel"
	persistence "github.com/dmibod/kanban/shared/persistence/board"
	"github.com/dmibod/kanban/shared/persistence/models"
	"github.com/dmibod/kanban/shared/tools/logger"
)

type service struct {
	logger.Logger
	persistence.Repository
	tx.Service
}

// CreateService instance
func CreateService(s tx.Service, r persistence.Repository, l logger.Logger) Service {
	return &service{
		Logger:     l,
		Repository: r,
		Service:    s,
	}
}

// GetByID get by id
func (s *service) GetByID(ctx context.Context, id kernel.ID) (*Model, error) {
	var model *Model
	if err := s.Repository.FindByID(ctx, id, func(entity *models.Board) error {
		model = mapPersistentToModel(entity)
		return nil
	}); err != nil {
		s.Errorln(err)
		return nil, err
	}

	return model, nil
}

// GetByOwner boards
func (s *service) GetByOwner(ctx context.Context, owner string) ([]*ListModel, error) {
	boards := []*ListModel{}
	err := s.Repository.FindByOwner(ctx, owner, func(entity *models.BoardListModel) error {
		boards = append(boards, mapPersistentToListModel(entity))
		return nil
	})

	if err != nil {
		s.Errorln(err)
		return nil, err
	}

	return boards, nil
}

// Create by payload
func (s *service) Create(ctx context.Context, model *CreateModel) (kernel.ID, error) {
	return s.create(ctx, model.Owner, func(aggregate board.Aggregate) error {
		if err := aggregate.Name(model.Name); err != nil {
			return err
		}
		if err := aggregate.Description(model.Description); err != nil {
			return err
		}
		if err := aggregate.Shared(model.Shared); err != nil {
			return err
		}
		return aggregate.Layout(model.Layout)
	})
}

// Layout board
func (s *service) Layout(ctx context.Context, id kernel.ID, layout string) error {
	return s.update(ctx, id, func(aggregate board.Aggregate) error {
		return aggregate.Layout(layout)
	})
}

// Name board
func (s *service) Name(ctx context.Context, id kernel.ID, name string) error {
	return s.update(ctx, id, func(aggregate board.Aggregate) error {
		return aggregate.Name(name)
	})
}

// Describe board
func (s *service) Describe(ctx context.Context, id kernel.ID, description string) error {
	return s.update(ctx, id, func(aggregate board.Aggregate) error {
		return aggregate.Description(description)
	})
}

// Share board
func (s *service) Share(ctx context.Context, id kernel.ID, shared bool) error {
	return s.update(ctx, id, func(aggregate board.Aggregate) error {
		return aggregate.Shared(shared)
	})
}

// AppendLane to board
func (s *service) AppendLane(ctx context.Context, id kernel.MemberID) error {
	return s.update(ctx, id.SetID, func(aggregate board.Aggregate) error {
		return aggregate.AppendChild(id.ID)
	})
}

// ExcludeLane from board
func (s *service) ExcludeLane(ctx context.Context, id kernel.MemberID) error {
	return s.update(ctx, id.SetID, func(aggregate board.Aggregate) error {
		return aggregate.RemoveChild(id.ID)
	})
}

// Remove by id
func (s *service) Remove(ctx context.Context, id kernel.ID) error {
	var entity *models.Board
	if err := s.Repository.FindByID(ctx, id, func(board *models.Board) error {
		entity = board
		return nil
	}); err != nil {
		s.Errorln(err)
		return err
	}

	return s.Service.Execute(ctx, func(bus event.Bus) error {
		return board.CreateService(bus).Delete(mapPersistentToDomain(entity))
	})
}

func (s *service) checkCreate(ctx context.Context) error {
	return nil
}

func (s *service) create(ctx context.Context, owner string, operation func(board.Aggregate) error) (kernel.ID, error) {
	if err := s.checkCreate(ctx); err != nil {
		s.Errorln(err)
		return kernel.EmptyID, err
	}

	id := kernel.ID(bson.NewObjectId().Hex())

	err := s.Service.Execute(ctx, func(bus event.Bus) error {
		domainService := board.CreateService(bus)

		entity, err := domainService.Create(id, owner)
		if err != nil {
			s.Errorln(err)
			return err
		}

		aggregate, err := domainService.Get(*entity)
		if err != nil {
			s.Errorln(err)
			return err
		}

		return operation(aggregate)
	})

	return id, err
}

func (s *service) checkUpdate(ctx context.Context, aggregate board.Aggregate) error {
	//TODO
	//securityContext := ctx.Value(scKey).(*SecurityContext)
	//if securityContext == nil || !securityContext.IsOwner(aggregate.GetOwner()) { return ErrOperationIsNotAllowed }
	return nil
}

func (s *service) update(ctx context.Context, id kernel.ID, operation func(board.Aggregate) error) error {
	var entity *models.Board
	if err := s.Repository.FindByID(ctx, id, func(board *models.Board) error {
		entity = board
		return nil
	}); err != nil {
		s.Errorln(err)
		return err
	}

	return s.Service.Execute(ctx, func(bus event.Bus) error {
		aggregate, err := board.CreateService(bus).Get(mapPersistentToDomain(entity))
		if err != nil {
			s.Errorln(err)
			return err
		}

		err = s.checkUpdate(ctx, aggregate)
		if err != nil {
			s.Errorln(err)
			return err
		}

		return operation(aggregate)
	})
}

func mapCardToModel(entity models.Card) CardModel {
	return CardModel{
		ID:          kernel.ID(entity.ID.Hex()),
		Name:        entity.Name,
		Description: entity.Description,
	}
}

func mapLaneToModel(entity models.Lane, lanes map[string]models.Lane, cards map[string]models.Card) LaneModel {
	var childLanes []LaneModel
	var childCards []CardModel

	if entity.Kind == kernel.LKind {
		childLanes = make([]LaneModel, len(entity.Children))
		for i, id := range entity.Children {
			childLanes[i] = mapLaneToModel(lanes[id.Hex()], lanes, cards)
		}
	} else {
		childCards = make([]CardModel, len(entity.Children))
		for i, id := range entity.Children {
			childCards[i] = mapCardToModel(cards[id.Hex()])
		}
	}

	return LaneModel{
		ID:          kernel.ID(entity.ID.Hex()),
		Type:        entity.Kind,
		Name:        entity.Name,
		Description: entity.Description,
		Layout:      entity.Layout,
		Lanes:       childLanes,
		Cards:       childCards,
	}
}

func mapPersistentToModel(entity *models.Board) *Model {
	lanes := make(map[string]models.Lane)
	for _, lane := range entity.Lanes {
		lanes[lane.ID.Hex()] = lane
	}

	cards := make(map[string]models.Card)
	for _, card := range entity.Cards {
		cards[card.ID.Hex()] = card
	}

	children := make([]LaneModel, len(entity.Children))
	for i, id := range entity.Children {
		children[i] = mapLaneToModel(lanes[id.Hex()], lanes, cards)
	}

	return &Model{
		ID:          kernel.ID(entity.ID.Hex()),
		Owner:       entity.Owner,
		Name:        entity.Name,
		Description: entity.Description,
		Layout:      entity.Layout,
		Shared:      entity.Shared,
		Lanes:       children,
	}
}

func mapPersistentToListModel(entity *models.BoardListModel) *ListModel {
	return &ListModel{
		ID:          kernel.ID(entity.ID.Hex()),
		Owner:       entity.Owner,
		Name:        entity.Name,
		Description: entity.Description,
		Layout:      entity.Layout,
		Shared:      entity.Shared,
	}
}

func mapPersistentToDomain(entity *models.Board) board.Entity {
	children := make([]kernel.ID, len(entity.Children))
	for i, id := range entity.Children {
		children[i] = kernel.ID(id.Hex())
	}
	return board.Entity{
		ID:          kernel.ID(entity.ID.Hex()),
		Owner:       entity.Owner,
		Name:        entity.Name,
		Description: entity.Description,
		Layout:      entity.Layout,
		Shared:      entity.Shared,
		Children:    children,
	}
}
