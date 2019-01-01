package services

import (
	"context"
	"errors"

	"gopkg.in/mgo.v2/bson"

	"github.com/dmibod/kanban/shared/kernel"
	"github.com/dmibod/kanban/shared/persistence"
	"github.com/dmibod/kanban/shared/tools/db"
	"github.com/dmibod/kanban/shared/tools/logger"
)

// BoardPayload represents payload
type BoardPayload struct {
	Name string
}

// BoardModel represents model
type BoardModel struct {
	ID   kernel.Id
	Name string
}

// BoardService interface
type BoardService interface {
	// Create by payload
	Create(context.Context, *BoardPayload) (kernel.Id, error)
	// Update model
	Update(context.Context, *BoardModel) (*BoardModel, error)
	// Remove remove by id
	Remove(context.Context, kernel.Id) error
	// GetByID get by id
	GetByID(context.Context, kernel.Id) (*BoardModel, error)
	// GetAll boards
	GetAll(context.Context) ([]*BoardModel, error)
}

type boardService struct {
	logger.Logger
	db.Repository
}

// Create by payload
func (s *boardService) Create(ctx context.Context, p *BoardPayload) (kernel.Id, error) {
	e := &persistence.BoardEntity{Name: p.Name}
	id, err := s.Repository.Create(ctx, e)
	if err != nil {
		s.Errorf("create error: %v\n", err)
		return "", err
	}

	return kernel.Id(id), nil
}

// Update model
func (s *boardService) Update(ctx context.Context, m *BoardModel) (*BoardModel, error) {
	entity := &persistence.BoardEntity{ID: bson.ObjectIdHex(string(m.ID)), Name: m.Name}
	err := s.Repository.Update(ctx, entity)
	if err != nil {
		s.Errorf("update error: %v\n", err)
		return nil, err
	}

	return &BoardModel{
		ID:   kernel.Id(entity.ID.Hex()),
		Name: entity.Name,
	}, nil
}

// Remove remove by id
func (s *boardService) Remove(ctx context.Context, id kernel.Id) error {
	err := s.Repository.Remove(ctx, string(id))
	if err != nil {
		s.Errorf("remove error: %v\n", err)
	}

	return err
}

// GetByID get by id
func (s *boardService) GetByID(ctx context.Context, id kernel.Id) (*BoardModel, error) {
	entity, err := s.Repository.FindByID(ctx, string(id))
	if err != nil {
		s.Errorf("error getting by id %v\n", id)
		return nil, err
	}

	board, ok := entity.(*persistence.BoardEntity)
	if !ok {
		s.Errorf("invalid type %T\n", board)
		return nil, errors.New("Invalid type")
	}

	return &BoardModel{
		ID:   kernel.Id(board.ID.Hex()),
		Name: board.Name,
	}, nil
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

		model := &BoardModel{
			ID:   kernel.Id(board.ID.Hex()),
			Name: board.Name,
		}

		models = append(models, model)

		return nil
	})

	if err != nil {
		s.Errorln(err)
		return nil, err
	}

	return models, nil
}
