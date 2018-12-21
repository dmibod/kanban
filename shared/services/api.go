package services

import (
	"github.com/dmibod/kanban/shared/kernel"
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
