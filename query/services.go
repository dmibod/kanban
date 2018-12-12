package query

import (
	"errors"
	"log"

	"github.com/dmibod/kanban/kernel"
	"github.com/dmibod/kanban/tools/db"
)

type CardService struct {
	repository db.Repository
}

func CreateCardService(repository db.Repository) *CardService {

	return &CardService{
		repository: repository,
	}
}

func (s *CardService) GetCardById(id kernel.Id) (*Card, error) {

	entity, err := s.repository.FindById(string(id))
	if err != nil {
		log.Printf("Error getting card by id %v\n", id)
		return nil, err
	}

	card, ok := entity.(*DbCard)

	if !ok {
		log.Printf("Invalid card type %T\n", card)
		return nil, errors.New("Invalid type")
	}

	return &Card{
		Id:   kernel.Id(card.Id.Hex()),
		Name: card.Name,
	}, nil
}
