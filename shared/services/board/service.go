package board

import (
	"context"

	"github.com/dmibod/kanban/shared/services/notification"

	"github.com/dmibod/kanban/shared/domain/board"
	"github.com/dmibod/kanban/shared/domain/event"

	"gopkg.in/mgo.v2/bson"

	"github.com/dmibod/kanban/shared/kernel"
	"github.com/dmibod/kanban/shared/persistence"
	"github.com/dmibod/kanban/shared/tools/logger"
)

type service struct {
	logger.Logger
	persistence.Repository
	notification.Service
}

// CreateService instance
func CreateService(s notification.Service, r persistence.Repository, l logger.Logger) Service {
	return &service{
		Logger:     l,
		Repository: r,
		Service:    s,
	}
}

// GetByID get by id
func (s *service) GetByID(ctx context.Context, id kernel.ID) (*Model, error) {
	var model *Model
	if err := s.Repository.FindBoardByID(ctx, id, func (entity *persistence.Board) error {
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
	models := []*ListModel{}
	err := s.Repository.FindBoardsByOwner(ctx, owner, func(entity *persistence.BoardListModel) error {
		models = append(models, mapPersistentToListModel(entity))
		return nil
	})

	if err != nil {
		s.Errorln(err)
		return nil, err
	}

	return models, nil
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
	return event.Execute(func(bus event.Bus) error {
		s.Service.Listen(bus)

		domainService := board.CreateService(bus)

		return domainService.Delete(board.Entity{ID: id})
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

	err := event.Execute(func(bus event.Bus) error {
		s.Service.Listen(bus)

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

		err = operation(aggregate)
		if err != nil {
			s.Errorln(err)
			return err
		}

		return err
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
	var entity *persistence.Board
	if err := s.Repository.FindBoardByID(ctx, id, func(board *persistence.Board) error {
		entity = board
		return nil
	}); err != nil {
		s.Errorln(err)
		return err
	}

	return event.Execute(func(bus event.Bus) error {
		s.Service.Listen(bus)

		domainService := board.CreateService(bus)

		aggregate, err := domainService.Get(mapPersistentToDomain(entity))
		if err != nil {
			s.Errorln(err)
			return err
		}

		err = s.checkUpdate(ctx, aggregate)
		if err != nil {
			s.Errorln(err)
			return err
		}

		err = operation(aggregate)
		if err == nil {
			bus.Fire()
			return nil
		}

		return err
	})
}

func mapPersistentToModel(entity *persistence.Board) *Model {
	return &Model{
		ID:          kernel.ID(entity.ID.Hex()),
		Owner:       entity.Owner,
		Name:        entity.Name,
		Description: entity.Description,
		Layout:      entity.Layout,
		Shared:      entity.Shared,
	}
}

func mapPersistentToListModel(entity *persistence.BoardListModel) *ListModel {
	return &ListModel{
		ID:          kernel.ID(entity.ID.Hex()),
		Owner:       entity.Owner,
		Name:        entity.Name,
		Description: entity.Description,
		Layout:      entity.Layout,
		Shared:      entity.Shared,
	}
}

func mapPersistentToDomain(entity *persistence.Board) board.Entity {
	children := []kernel.ID{}
	for _, id := range entity.Children {
		children = append(children, kernel.ID(id))
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