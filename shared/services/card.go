package services

import (
	"errors"
	"github.com/dmibod/kanban/shared/persistence"
	"github.com/dmibod/kanban/shared/tools/db"
	"github.com/dmibod/kanban/shared/tools/log"
	"github.com/dmibod/kanban/shared/kernel"
)

// CardPayload represents card fields without id
type CardPayload struct {
	Name string
}

// CardModel represents card at service layer
type CardModel struct {
	ID kernel.Id
	Name string
}

// CardService holds service dependencies
type CardService struct {
	logger     log.Logger
	repository db.Repository
}

// CreateCardService creates new CardService instance
func CreateCardService(l log.Logger, r db.Repository) *CardService {
	return &CardService{
		logger:     l,
		repository: r,
	}
}

// CreateCard creates new card
func (s *CardService) CreateCard(p *CardPayload) (kernel.Id, error) {
	e := &persistence.CardEntity{Name: p.Name}
	id, err := s.repository.Create(e)
	if err != nil {
		s.logger.Errorf("create card error: %v\n%v\n", err, p)
		return "", err
	}

	return kernel.Id(id), nil
}

// GetCardByID reads card from db by its id
func (s *CardService) GetCardByID(id kernel.Id) (*CardModel, error) {
	entity, err := s.repository.FindByID(string(id))
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
