package services

import (
	"context"

	"github.com/dmibod/kanban/shared/domain/board"
	"github.com/dmibod/kanban/shared/domain/event"

	"gopkg.in/mgo.v2/bson"

	"github.com/dmibod/kanban/shared/kernel"
	"github.com/dmibod/kanban/shared/persistence"
	"github.com/dmibod/kanban/shared/tools/logger"
)

// BoardPayload represents payload
type BoardPayload struct {
	Owner       string
	Name        string
	Description string
	Layout      string
}

// BoardModel represents model
type BoardModel struct {
	ID          kernel.ID
	Owner       string
	Name        string
	Description string
	Shared      bool
	Layout      string
}

// BoardReader interface
type BoardReader interface {
	// GetByID get by id
	GetByID(context.Context, kernel.ID) (*BoardModel, error)
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
	Layout(context.Context, kernel.ID, string) (*BoardModel, error)
	// Name board
	Name(context.Context, kernel.ID, string) (*BoardModel, error)
	// Describe board
	Describe(context.Context, kernel.ID, string) (*BoardModel, error)
	// Share board
	Share(context.Context, kernel.ID, bool) (*BoardModel, error)
	// Remove board by id
	Remove(context.Context, kernel.ID) error
	// AppendChild to lane
	AppendChild(context.Context, kernel.ID, kernel.ID) error
	// ExcludeChild from lane
	ExcludeChild(context.Context, kernel.ID, kernel.ID) error
}

// BoardService interface
type BoardService interface {
	BoardReader
	BoardWriter
}

type boardService struct {
	logger.Logger
	BoardRepository *persistence.BoardRepository
	NotificationService
}

// CreateBoardService instance
func CreateBoardService(s NotificationService, r *persistence.BoardRepository, l logger.Logger) BoardService {
	return &boardService{
		Logger:              l,
		BoardRepository:     r,
		NotificationService: s,
	}
}

// GetByID get by id
func (s *boardService) GetByID(ctx context.Context, id kernel.ID) (*BoardModel, error) {
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
	return s.createAndGet(ctx, payload.Owner, func(aggregate board.Aggregate) error {
		if err := aggregate.Name(payload.Name); err != nil {
			return err
		}
		if err := aggregate.Description(payload.Description); err != nil {
			return err
		}
		return aggregate.Layout(payload.Layout)
	})
}

// Layout board
func (s *boardService) Layout(ctx context.Context, id kernel.ID, layout string) (*BoardModel, error) {
	return s.updateAndGet(ctx, id, func(aggregate board.Aggregate) error {
		return aggregate.Layout(layout)
	})
}

// Name board
func (s *boardService) Name(ctx context.Context, id kernel.ID, name string) (*BoardModel, error) {
	return s.updateAndGet(ctx, id, func(aggregate board.Aggregate) error {
		return aggregate.Name(name)
	})
}

// Describe board
func (s *boardService) Describe(ctx context.Context, id kernel.ID, description string) (*BoardModel, error) {
	return s.updateAndGet(ctx, id, func(aggregate board.Aggregate) error {
		return aggregate.Description(description)
	})
}

// Share board
func (s *boardService) Share(ctx context.Context, id kernel.ID, shared bool) (*BoardModel, error) {
	return s.updateAndGet(ctx, id, func(aggregate board.Aggregate) error {
		return aggregate.Shared(shared)
	})
}

// AppendChild to board
func (s *boardService) AppendChild(ctx context.Context, id kernel.ID, childID kernel.ID) error {
	return s.update(ctx, id, func(aggregate board.Aggregate) error {
		return aggregate.AppendChild(childID)
	})
}

// ExcludeChild from board
func (s *boardService) ExcludeChild(ctx context.Context, id kernel.ID, childID kernel.ID) error {
	return s.update(ctx, id, func(aggregate board.Aggregate) error {
		return aggregate.RemoveChild(childID)
	})
}

// Remove by id
func (s *boardService) Remove(ctx context.Context, id kernel.ID) error {
	return event.Execute(func(bus event.Bus) error {
		s.NotificationService.Listen(bus)
		s.BoardRepository.Listen(ctx, bus)
		return board.Delete(board.Entity{ID: id}, bus)
	})
}

func (s *boardService) checkCreate(ctx context.Context) error {
	return nil
}

func (s *boardService) create(ctx context.Context, owner string, operation func(board.Aggregate) error) (kernel.ID, error) {
	if err := s.checkCreate(ctx); err != nil {
		s.Errorln(err)
		return kernel.EmptyID, err
	}

	id := kernel.ID(bson.NewObjectId().Hex())

	err := event.Execute(func(bus event.Bus) error {
		s.NotificationService.Listen(bus)
		s.BoardRepository.Listen(ctx, bus)

		entity, err := board.Create(id, owner, bus)
		if err != nil {
			s.Errorln(err)
			return err
		}

		aggregate, err := board.New(*entity, bus)
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

func (s *boardService) createAndGet(ctx context.Context, owner string, operation func(board.Aggregate) error) (*BoardModel, error) {
	id, err := s.create(ctx, owner, operation)
	if err != nil {
		s.Errorln(err)
		return nil, err
	}

	return s.GetByID(ctx, id)
}

func (s *boardService) checkUpdate(ctx context.Context, aggregate board.Aggregate) error {
	//TODO
	//securityContext := ctx.Value(scKey).(*SecurityContext)
	//if securityContext == nil || !securityContext.IsOwner(aggregate.GetOwner()) { return ErrOperationIsNotAllowed }
	return nil
}

func (s *boardService) update(ctx context.Context, id kernel.ID, operation func(board.Aggregate) error) error {
	entity, err := s.BoardRepository.FindBoardByID(ctx, id)
	if err != nil {
		s.Errorln(err)
		return err
	}

	return event.Execute(func(bus event.Bus) error {
		s.NotificationService.Listen(bus)
		s.BoardRepository.Listen(ctx, bus)

		aggregate, err := board.New(mapBoardEntityToEntity(entity), bus)
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
			aggregate.Save()
		}

		return err
	})
}

func (s *boardService) updateAndGet(ctx context.Context, id kernel.ID, operation func(board.Aggregate) error) (*BoardModel, error) {
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
		ID:          kernel.ID(entity.ID.Hex()),
		Owner:       entity.Owner,
		Name:        entity.Name,
		Description: entity.Description,
		Layout:      entity.Layout,
		Shared:      entity.Shared,
	}
}

func mapBoardEntityToEntity(entity *persistence.BoardEntity) board.Entity {
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
