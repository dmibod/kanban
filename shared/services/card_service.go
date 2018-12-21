package services

import (
	"context"
	"errors"

	"github.com/dmibod/kanban/shared/kernel"
	"github.com/dmibod/kanban/shared/persistence"
	"github.com/dmibod/kanban/shared/tools/db"
	"github.com/dmibod/kanban/shared/tools/logger"
)

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
