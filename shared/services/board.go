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
	Create(*BoardPayload) (kernel.Id, error)
	// Update model
	Update(*BoardModel) (*BoardModel, error)
	// Remove remove by id
	Remove(kernel.Id) error
	// GetByID get by id
	GetByID(kernel.Id) (*BoardModel, error)
}

type boardService struct {
	context.Context
	logger.Logger
	db.RepositoryFactory
}

// Create by payload
func (s *boardService) Create(p *BoardPayload) (kernel.Id, error) {
	e := &persistence.BoardEntity{Name: p.Name}
	id, err := s.getRepository().Create(e)
	if err != nil {
		s.Errorf("create error: %v\n", err)
		return "", err
	}

	return kernel.Id(id), nil
}

// Update model
func (s *boardService) Update(m *BoardModel) (*BoardModel, error) {
	entity := &persistence.BoardEntity{ID: bson.ObjectIdHex(string(m.ID)), Name: m.Name}
	err := s.getRepository().Update(entity)
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
func (s *boardService) Remove(id kernel.Id) error {
	err := s.getRepository().Remove(string(id))
	if err != nil {
		s.Errorf("remove error: %v\n", err)
	}

	return err
}

// GetByID get by id
func (s *boardService) GetByID(id kernel.Id) (*BoardModel, error) {
	entity, err := s.getRepository().FindByID(string(id))
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

func (s *boardService) getRepository() db.Repository {
	return persistence.CreateBoardRepository(s.Context, s.RepositoryFactory)
}
