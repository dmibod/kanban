package services

import (
	"context"
	"errors"

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
	CreateCard(*CardPayload) (kernel.Id, error)
	// GetCardByID reads card from db by its id
	GetCardByID(kernel.Id) (*CardModel, error)
}

type cardService struct {
	ctx     context.Context
	logger  logger.Logger
	factory db.Factory
}

// CreateCard creates new card
func (s *cardService) CreateCard(p *CardPayload) (kernel.Id, error) {
	e := &persistence.CardEntity{Name: p.Name}
	id, err := s.getRepository().Create(e)
	if err != nil {
		s.logger.Errorf("create card error: %v\n%v\n", err, p)
		return "", err
	}

	return kernel.Id(id), nil
}

// GetCardByID reads card from db by its id
func (s *cardService) GetCardByID(id kernel.Id) (*CardModel, error) {
	entity, err := s.getRepository().FindByID(string(id))
	if err != nil {
		s.logger.Errorf("error getting card by id %v\n", id)
		return nil, err
	}

	card, ok := entity.(*persistence.CardEntity)
	if !ok {
		s.logger.Errorf("invalid card type %T\n", card)
		return nil, errors.New("Invalid type")
	}

	return &CardModel{
		ID:   kernel.Id(card.ID.Hex()),
		Name: card.Name,
	}, nil
}

func (s *cardService) getRepository() db.Repository {
	return persistence.CreateCardRepository(s.ctx, s.factory)
}
