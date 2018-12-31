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
	// CreateCard creates new card
	CreateCard(context.Context, *CardPayload) (kernel.Id, error)
	// UpdateCard updates card
	UpdateCard(context.Context, *CardModel) (*CardModel, error)
	// RemoveCard removes card
	RemoveCard(context.Context, kernel.Id) error
	// GetCardByID reads card from db by its id
	GetCardByID(context.Context, kernel.Id) (*CardModel, error)
}

type cardService struct {
	logger.Logger
	db.Repository
}

// CreateCard creates new card
func (s *cardService) CreateCard(ctx context.Context, p *CardPayload) (kernel.Id, error) {
	e := &persistence.CardEntity{Name: p.Name}
	id, err := s.Repository.Create(ctx, e)
	if err != nil {
		s.Errorf("create card error: %v\n%v\n", err, p)
		return "", err
	}

	return kernel.Id(id), nil
}

// UpdateCard updates card
func (s *cardService) UpdateCard(ctx context.Context, c *CardModel) (*CardModel, error) {
	e := &persistence.CardEntity{ID: bson.ObjectIdHex(string(c.ID)), Name: c.Name}
	err := s.Repository.Update(ctx, e)
	if err != nil {
		s.Errorf("update card error: %v\n", err)
		return nil, err
	}

	return &CardModel{
		ID:   kernel.Id(e.ID.Hex()),
		Name: e.Name,
	}, nil
}

// RemoveCard removes card
func (s *cardService) RemoveCard(ctx context.Context, id kernel.Id) error {
	err := s.Repository.Remove(ctx, string(id))
	if err != nil {
		s.Errorf("remove card error: %v\n", err)
	}

	return err
}

// GetCardByID reads card from db by its id
func (s *cardService) GetCardByID(ctx context.Context, id kernel.Id) (*CardModel, error) {
	entity, err := s.Repository.FindByID(ctx, string(id))
	if err != nil {
		s.Errorf("error getting card by id %v\n", id)
		return nil, err
	}

	card, ok := entity.(*persistence.CardEntity)
	if !ok {
		s.Errorf("invalid card type %T\n", card)
		return nil, errors.New("Invalid type")
	}

	return &CardModel{
		ID:   kernel.Id(card.ID.Hex()),
		Name: card.Name,
	}, nil
}
