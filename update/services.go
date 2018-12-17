package update

import (
	"github.com/dmibod/kanban/shared/kernel"
	"github.com/dmibod/kanban/shared/tools/db"
	"github.com/dmibod/kanban/shared/tools/log"
)

// CardPayload represents card fields without id
type CardPayload struct {
	Name string
}

// CardModel represents card at service layer
type CardModel struct {
	CardPayload
	ID kernel.Id
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
	id, err := s.repository.Create(p)
	if err != nil {
		s.logger.Errorf("create card error: %v\n%v\n", err, p)
		return "", err
	}

	return kernel.Id(id), nil
}
