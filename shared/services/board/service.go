package board

import (
	"context"

	"github.com/dmibod/kanban/shared/services/notification"

	"github.com/dmibod/kanban/shared/domain/board"
	"github.com/dmibod/kanban/shared/domain/event"

	"gopkg.in/mgo.v2/bson"

	"github.com/dmibod/kanban/shared/kernel"
	persistence "github.com/dmibod/kanban/shared/persistence/board"
	"github.com/dmibod/kanban/shared/tools/logger"
)

type service struct {
	logger.Logger
	Repository *persistence.Repository
	notification.Service
}

// CreateService instance
func CreateService(s notification.Service, r *persistence.Repository, l logger.Logger) Service {
	return &service{
		Logger:          l,
	  Repository: r,
		Service:         s,
	}
}

// GetByID get by id
func (s *service) GetByID(ctx context.Context, id kernel.ID) (*Model, error) {
	entity, err := s.Repository.FindBoardByID(ctx, id)
	if err != nil {
		s.Errorln(err)
		return nil, err
	}

	return mapPersistentToModel(entity), nil
}

// GetAll boards
func (s *service) GetAll(ctx context.Context) ([]*ListModel, error) {
	return s.getByCriteria(ctx, nil)
}

// GetByOwner boards
func (s *service) GetByOwner(ctx context.Context, owner string) ([]*ListModel, error) {
	return s.getByCriteria(ctx, buildOwnerCriteria(owner))
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

// AppendChild to board
func (s *service) AppendChild(ctx context.Context, id kernel.ID, childID kernel.ID) error {
	return s.update(ctx, id, func(aggregate board.Aggregate) error {
		return aggregate.AppendChild(childID)
	})
}

// ExcludeChild from board
func (s *service) ExcludeChild(ctx context.Context, id kernel.ID, childID kernel.ID) error {
	return s.update(ctx, id, func(aggregate board.Aggregate) error {
		return aggregate.RemoveChild(childID)
	})
}

// Remove by id
func (s *service) Remove(ctx context.Context, id kernel.ID) error {
	return event.Execute(func(bus event.Bus) error {
		s.Service.Listen(bus)

		domainService := board.CreateService(s.Repository.GetRepository(ctx), bus)

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

		domainService := board.CreateService(s.Repository.GetRepository(ctx), bus)

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
	entity, err := s.Repository.FindBoardByID(ctx, id)
	if err != nil {
		s.Errorln(err)
		return err
	}

	return event.Execute(func(bus event.Bus) error {
		s.Service.Listen(bus)

		domainService := board.CreateService(s.Repository.GetRepository(ctx), bus)

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
			return aggregate.Save()
		}

		return err
	})
}

func (s *service) getByCriteria(ctx context.Context, criteria bson.M) ([]*ListModel, error) {
	models := []*ListModel{}
	err := s.Repository.FindBoards(ctx, criteria, func(entity *persistence.BoardEntity) error {
		models = append(models, mapPersistentToListModel(entity))
		return nil
	})

	if err != nil {
		s.Errorln(err)
		return nil, err
	}

	return models, nil
}

func buildOwnerCriteria(owner string) bson.M {
	if owner == "" {
		return bson.M{"shared": true}
	}

	return bson.M{"$or": []bson.M{bson.M{"shared": true}, bson.M{"owner": owner}}}
}

func mapPersistentToModel(entity *persistence.BoardEntity) *Model {
	return &Model{
		ID:          kernel.ID(entity.ID.Hex()),
		Owner:       entity.Owner,
		Name:        entity.Name,
		Description: entity.Description,
		Layout:      entity.Layout,
		Shared:      entity.Shared,
	}
}

func mapPersistentToListModel(entity *persistence.BoardEntity) *ListModel {
	return &ListModel{
		ID:          kernel.ID(entity.ID.Hex()),
		Owner:       entity.Owner,
		Name:        entity.Name,
		Description: entity.Description,
		Layout:      entity.Layout,
		Shared:      entity.Shared,
	}
}

func mapPersistentToDomain(entity *persistence.BoardEntity) board.Entity {
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
