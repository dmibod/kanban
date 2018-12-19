package query

import (
	"github.com/dmibod/kanban/shared/persistence"
	"errors"

	"github.com/dmibod/kanban/shared/tools/log"
	"github.com/dmibod/kanban/shared/kernel"
)

// CardModel represents card at service layer
type CardModel struct {
	ID   kernel.Id
	Name string
}

// CardRepository repository expected by service
type CardRepository interface {
	FindByID(string) (interface{}, error)
}
// CardService exposes cards api at service layer
type CardService struct {
	logger     log.Logger
	repository CardRepository
}

// CreateCardService creates new instance of service
func CreateCardService(l log.Logger, r CardRepository) *CardService {
	return &CardService{
		logger:     l,
		repository: r,
	}
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
