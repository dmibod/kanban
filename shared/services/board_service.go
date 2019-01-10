package services

import (
	"context"
	"errors"

	"github.com/dmibod/kanban/shared/domain"

	"gopkg.in/mgo.v2/bson"

	"github.com/dmibod/kanban/shared/kernel"
	"github.com/dmibod/kanban/shared/persistence"
	"github.com/dmibod/kanban/shared/tools/db"
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

// BoardService interface
type BoardService interface {
	// Create by payload
	Create(context.Context, *BoardPayload) (kernel.Id, error)
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
	// GetByID get by id
	GetByID(context.Context, kernel.Id) (*BoardModel, error)
	// GetByOwner boards
	GetByOwner(context.Context, string) ([]*BoardModel, error)
	// GetAll boards
	GetAll(context.Context) ([]*BoardModel, error)
}

type boardService struct {
	logger.Logger
	db.Repository
	NotificationService
}

// Create by payload
func (s *boardService) Create(ctx context.Context, payload *BoardPayload) (kernel.Id, error) {
	entity := mapBoardPayloadToEntity(payload)

	id, err := s.Repository.Create(ctx, entity)
	if err != nil {
		s.Errorln(err)
		return "", err
	}

	return kernel.Id(id), nil
}

func (s *boardService) checkUpdate(ctx context.Context, aggregate domain.BoardAggregate) error {
	//TODO
	//securityContext := ctx.Value(scKey).(*SecurityContext)
	//if securityContext == nil || !securityContext.IsOwner(aggregate.GetOwner()) { return ErrOperationIsNotAllowed }
	return nil
}

func (s *boardService) update(ctx context.Context, id kernel.Id, operation func(domain.BoardAggregate) error) error {
	return s.NotificationService.Execute(func(e domain.EventRegistry) error {
		aggregate, err := domain.LoadBoard(id, persistence.CreateBoardDomainRepository(ctx, s.Repository), e)
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

// Remove by id
func (s *boardService) Remove(ctx context.Context, id kernel.Id) error {
	err := s.Repository.Remove(ctx, string(id))
	if err != nil {
		s.Errorln(err)
	}

	return err
}

// GetByID get by id
func (s *boardService) GetByID(ctx context.Context, id kernel.Id) (*BoardModel, error) {
	entity, err := s.Repository.FindByID(ctx, string(id))
	if err != nil {
		s.Errorln(err)
		return nil, err
	}

	board, ok := entity.(*persistence.BoardEntity)
	if !ok {
		s.Errorf("invalid type %T\n", entity)
		return nil, errors.New("Invalid type")
	}

	return mapBoardEntityToModel(board), nil
}

// GetAll boards
func (s *boardService) GetAll(ctx context.Context) ([]*BoardModel, error) {
	models := []*BoardModel{}
	err := s.Repository.Find(ctx, nil, func(entity interface{}) error {
		board, ok := entity.(*persistence.BoardEntity)
		if !ok {
			s.Errorf("invalid type %T\n", entity)
			return errors.New("Invalid type")
		}

		models = append(models, mapBoardEntityToModel(board))

		return nil
	})

	if err != nil {
		s.Errorln(err)
		return nil, err
	}

	return models, nil
}

// GetByOwner boards
func (s *boardService) GetByOwner(ctx context.Context, owner string) ([]*BoardModel, error) {
	models := []*BoardModel{}

	var criteria bson.M

	if owner == "" {
		criteria = bson.M{"shared": true}
	} else {
		criteria = bson.M{"$or": []bson.M{bson.M{"shared": true}, bson.M{"owner": owner}}}
	}

	err := s.Repository.Find(ctx, criteria, func(entity interface{}) error {
		board, ok := entity.(*persistence.BoardEntity)
		if !ok {
			s.Errorf("invalid type %T\n", entity)
			return errors.New("Invalid type")
		}

		models = append(models, mapBoardEntityToModel(board))

		return nil
	})

	if err != nil {
		s.Errorln(err)
		return nil, err
	}

	return models, nil
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

func mapBoardModelToEntity(model *BoardModel) *persistence.BoardEntity {
	return &persistence.BoardEntity{
		ID:     bson.ObjectIdHex(string(model.ID)),
		Name:   model.Name,
		Layout: model.Layout,
	}
}

func mapBoardPayloadToEntity(model *BoardPayload) *persistence.BoardEntity {
	return &persistence.BoardEntity{
		Name:   model.Name,
		Layout: model.Layout,
		Owner:  model.Owner,
	}
}
