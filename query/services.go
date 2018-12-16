package query

import (
	"errors"
	"log"

	"github.com/dmibod/kanban/kernel"
	"github.com/dmibod/kanban/tools/db"
)

// CardModel represents card at service layer
type CardModel struct {
	ID   kernel.Id
	Name string
}

// CardService exposes cards api at service layer
type CardService struct {
	Repository interface {
		FindById(string) (interface{}, error)
	}
}

// CreateCardService creates new instance of service
func CreateCardService(r db.Repository) *CardService {
	return &CardService{
		Repository: r,
	}
}

// GetCardByID reads card from db by its id
func (s *CardService) GetCardByID(id kernel.Id) (*CardModel, error) {

	entity, err := s.Repository.FindById(string(id))
	if err != nil {
		log.Printf("Error getting card by id %v\n", id)
		return nil, err
	}

	card, ok := entity.(*CardEntity)

	if !ok {
		log.Printf("Invalid card type %T\n", card)
		return nil, errors.New("Invalid type")
	}

	return &CardModel{
		ID:   kernel.Id(card.ID.Hex()),
		Name: card.Name,
	}, nil
}
