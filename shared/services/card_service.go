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

// CardPayload represents card fields without id
type CardPayload struct {
	Name string
}

// CardModel represents card at service layer
type CardModel struct {
	ID   kernel.Id
	Name string
}

// CardService interface
type CardService interface {
	// Create card
	Create(context.Context, *CardPayload) (kernel.Id, error)
	// Update card
	Update(context.Context, *CardModel) (*CardModel, error)
	// Remove card
	Remove(context.Context, kernel.Id) error
	// GetByID gets card by id
	GetByID(context.Context, kernel.Id) (*CardModel, error)
}

type cardService struct {
	logger.Logger
	db.Repository
}

// Create card
func (s *cardService) Create(ctx context.Context, p *CardPayload) (kernel.Id, error) {
	e := &persistence.CardEntity{Name: p.Name}
	id, err := s.Repository.Create(ctx, e)
	if err != nil {
		s.Errorln(err)
		return "", err
	}

	return kernel.Id(id), nil
}

// Update card
func (s *cardService) Update(ctx context.Context, c *CardModel) (*CardModel, error) {
	e := &persistence.CardEntity{ID: bson.ObjectIdHex(string(c.ID)), Name: c.Name}
	err := s.Repository.Update(ctx, e)
	if err != nil {
		s.Errorln(err)
		return nil, err
	}

	return &CardModel{
		ID:   kernel.Id(e.ID.Hex()),
		Name: e.Name,
	}, nil
}

// Remove card
func (s *cardService) Remove(ctx context.Context, id kernel.Id) error {
	err := s.Repository.Remove(ctx, string(id))
	if err != nil {
		s.Errorln(err)
	}

	return err
}

// GetByID gets card by id
func (s *cardService) GetByID(ctx context.Context, id kernel.Id) (*CardModel, error) {
	entity, err := s.Repository.FindByID(ctx, string(id))
	if err != nil {
		s.Errorln(err)
		return nil, err
	}

	card, ok := entity.(*persistence.CardEntity)
	if !ok {
		s.Errorf("invalid type %T\n", entity)
		return nil, errors.New("Invalid type")
	}

	return &CardModel{
		ID:   kernel.Id(card.ID.Hex()),
		Name: card.Name,
	}, nil
}
