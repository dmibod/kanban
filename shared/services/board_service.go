package services

import (
	"context"

	"github.com/dmibod/kanban/shared/domain"

	"gopkg.in/mgo.v2/bson"

	"github.com/dmibod/kanban/shared/kernel"
	"github.com/dmibod/kanban/shared/persistence"
	"github.com/dmibod/kanban/shared/tools/logger"
)

// BoardPayload represents payload
type BoardPayload struct {
	Name   string
	Layout string
	Owner  string
}

// BoardModel represents model
type BoardModel struct {
	ID     kernel.Id
	Layout string
	Name   string
	Owner  string
	Shared bool
}

// BoardReader interface
type BoardReader interface {
	// GetByID get by id
	GetByID(context.Context, kernel.Id) (*BoardModel, error)
	// GetByOwner boards
	GetByOwner(context.Context, string) ([]*BoardModel, error)
	// GetAll boards
	GetAll(context.Context) ([]*BoardModel, error)
}

// BoardWriter interface
type BoardWriter interface {
	// Create by payload
	Create(context.Context, *BoardPayload) (*BoardModel, error)
	// Layout board
	Layout(context.Context, kernel.Id, string) (*BoardModel, error)
	// Rename board
	Rename(context.Context, kernel.Id, string) (*BoardModel, error)
	// Share board
	Share(context.Context, kernel.Id, bool) (*BoardModel, error)
	// Remove board by id
	Remove(context.Context, kernel.Id) error
	// AppendChild to lane
	AppendChild(context.Context, kernel.Id, kernel.Id) error
	// ExcludeChild from lane
	ExcludeChild(context.Context, kernel.Id, kernel.Id) error
}

// BoardService interface
type BoardService interface {
	BoardReader
	BoardWriter
}

type boardService struct {
	logger.Logger
	persistence.BoardRepository
	NotificationService
}

// GetByID get by id
func (s *boardService) GetByID(ctx context.Context, id kernel.Id) (*BoardModel, error) {
	entity, err := s.BoardRepository.FindBoardByID(ctx, id)
	if err != nil {
		s.Errorln(err)
		return nil, err
	}

	return mapBoardEntityToModel(entity), nil
}

// GetAll boards
func (s *boardService) GetAll(ctx context.Context) ([]*BoardModel, error) {
	return s.getByCriteria(ctx, nil)
}

// GetByOwner boards
func (s *boardService) GetByOwner(ctx context.Context, owner string) ([]*BoardModel, error) {
	return s.getByCriteria(ctx, buildBoardOwnerCriteria(owner))
}

// Create by payload
func (s *boardService) Create(ctx context.Context, payload *BoardPayload) (*BoardModel, error) {
	return s.createAndGet(ctx, payload.Owner, func(aggregate domain.BoardAggregate) error {
		if err := aggregate.Name(payload.Name); err != nil {
			return err
		}
		return aggregate.Layout(payload.Layout)
	})
}

// Layout board
func (s *boardService) Layout(ctx context.Context, id kernel.Id, layout string) (*BoardModel, error) {
	return s.updateAndGet(ctx, id, func(aggregate domain.BoardAggregate) error {
		return aggregate.Layout(layout)
	})
}

// Rename board
func (s *boardService) Rename(ctx context.Context, id kernel.Id, name string) (*BoardModel, error) {
	return s.updateAndGet(ctx, id, func(aggregate domain.BoardAggregate) error {
		return aggregate.Name(name)
	})
}

// Share board
func (s *boardService) Share(ctx context.Context, id kernel.Id, shared bool) (*BoardModel, error) {
	return s.updateAndGet(ctx, id, func(aggregate domain.BoardAggregate) error {
		return aggregate.Shared(shared)
	})
}

// AppendChild to board
func (s *boardService) AppendChild(ctx context.Context, id kernel.Id, childID kernel.Id) error {
	return s.update(ctx, id, func(aggregate domain.BoardAggregate) error {
		return aggregate.AppendChild(childID)
	})
}

// ExcludeChild from board
func (s *boardService) ExcludeChild(ctx context.Context, id kernel.Id, childID kernel.Id) error {
	return s.update(ctx, id, func(aggregate domain.BoardAggregate) error {
		return aggregate.RemoveChild(childID)
	})
}

// Remove by id
func (s *boardService) Remove(ctx context.Context, id kernel.Id) error {
	err := s.BoardRepository.Remove(ctx, string(id))
	if err != nil {
		s.Errorln(err)
	}

	return err
}

func (s *boardService) checkCreate(ctx context.Context, aggregate domain.BoardAggregate) error {
	return nil
}

func (s *boardService) create(ctx context.Context, owner string, operation func(domain.BoardAggregate) error) (kernel.Id, error) {
	id := kernel.EmptyID
	err := s.NotificationService.Execute(func(e domain.EventRegistry) error {
		aggregate, err := domain.NewBoard(owner, s.BoardRepository.DomainRepository(ctx), e)
		if err != nil {
			s.Errorln(err)
			return err
		}

		err = s.checkCreate(ctx, aggregate)
		if err != nil {
			s.Errorln(err)
			return err
		}

		err = operation(aggregate)
		if err != nil {
			s.Errorln(err)
			return err
		}

		err = aggregate.Save()
		if err == nil {
			id = aggregate.GetID()
		}

		return err
	})

	return id, err
}

func (s *boardService) createAndGet(ctx context.Context, owner string, operation func(domain.BoardAggregate) error) (*BoardModel, error) {
	id, err := s.create(ctx, owner, operation)
	if err != nil {
		s.Errorln(err)
		return nil, err
	}

	return s.GetByID(ctx, id)
}

func (s *boardService) checkUpdate(ctx context.Context, aggregate domain.BoardAggregate) error {
	//TODO
	//securityContext := ctx.Value(scKey).(*SecurityContext)
	//if securityContext == nil || !securityContext.IsOwner(aggregate.GetOwner()) { return ErrOperationIsNotAllowed }
	return nil
}

func (s *boardService) update(ctx context.Context, id kernel.Id, operation func(domain.BoardAggregate) error) error {
	return s.NotificationService.Execute(func(e domain.EventRegistry) error {
		aggregate, err := domain.LoadBoard(id, s.BoardRepository.DomainRepository(ctx), e)
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
		if err != nil {
			s.Errorln(err)
			return err
		}

		return aggregate.Save()
	})
}

func (s *boardService) updateAndGet(ctx context.Context, id kernel.Id, operation func(domain.BoardAggregate) error) (*BoardModel, error) {
	err := s.update(ctx, id, operation)
	if err != nil {
		s.Errorln(err)
		return nil, err
	}

	return s.GetByID(ctx, id)
}

func (s *boardService) getByCriteria(ctx context.Context, criteria bson.M) ([]*BoardModel, error) {
	models := []*BoardModel{}
	err := s.BoardRepository.FindBoards(ctx, criteria, func(entity *persistence.BoardEntity) error {
		models = append(models, mapBoardEntityToModel(entity))
		return nil
	})

	if err != nil {
		s.Errorln(err)
		return nil, err
	}

	return models, nil
}

func buildBoardOwnerCriteria(owner string) bson.M {
	if owner == "" {
		return bson.M{"shared": true}
	}

	return bson.M{"$or": []bson.M{bson.M{"shared": true}, bson.M{"owner": owner}}}
}

func mapBoardEntityToModel(entity *persistence.BoardEntity) *BoardModel {
	return &BoardModel{
		ID:     kernel.Id(entity.ID.Hex()),
		Name:   entity.Name,
		Layout: entity.Layout,
		Owner:  entity.Owner,
		Shared: entity.Shared,
	}
}
