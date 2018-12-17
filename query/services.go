package query

import (
	"github.com/dmibod/kanban/shared/persistence"
	"errors"

	"github.com/dmibod/kanban/shared/tools/log"
	"github.com/dmibod/kanban/shared/kernel"
	"github.com/dmibod/kanban/shared/tools/db"
)

// CardModel represents card at service layer
type CardModel struct {
	ID   kernel.Id
	Name string
}

// CardService exposes cards api at service layer
type CardService struct {
	Logger     log.Logger
	Repository interface {
		FindById(string) (interface{}, error)
	}
}

// CreateCardService creates new instance of service
func CreateCardService(l log.Logger, r db.Repository) *CardService {
	return &CardService{
		Logger:     l,
		Repository: r,
	}
}

// GetCardByID reads card from db by its id
func (s *CardService) GetCardByID(id kernel.Id) (*CardModel, error) {
	entity, err := s.Repository.FindById(string(id))
	if err != nil {
		s.Logger.Errorf("error getting card by id %v\n", id)
		return nil, err
	}

	card, ok := entity.(*persistence.CardEntity)
	if !ok {
		s.Logger.Errorf("invalid card type %T\n", card)
		return nil, errors.New("Invalid type")
	}

	return &CardModel{
		ID:   kernel.Id(card.ID.Hex()),
		Name: card.Name,
	}, nil
}
