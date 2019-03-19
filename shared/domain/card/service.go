package card

import (
	err "github.com/dmibod/kanban/shared/domain/error"
	"github.com/dmibod/kanban/shared/domain/event"
	"github.com/dmibod/kanban/shared/kernel"
)

// Service - card domain service
type Service struct {
	Repository
	event.Bus
}

// CreateService - creates card domain service
func CreateService(repository Repository, bus event.Bus) *Service {
	if repository == nil {
		return nil
	}

	if bus == nil {
		return nil
	}

	return &Service{Repository: repository, Bus: bus}
}

// Create card
func (s *Service) Create(id kernel.ID) (*Entity, error) {
	if !id.IsValid() {
		return nil, err.ErrInvalidID
	}

	entity := Entity{
		ID: id,
	}

	if err := s.Repository.Create(&entity); err != nil {
		return nil, err
	}

	s.Bus.Register(CreatedEvent{entity})
	s.Bus.Fire()

	return &entity, nil
}

// Delete card
func (s *Service) Delete(entity Entity) error {
	if !entity.ID.IsValid() {
		return err.ErrInvalidID
	}

	if err := s.Repository.Delete(&entity); err != nil {
		return err
	}

	s.Bus.Register(DeletedEvent{entity})
	s.Bus.Fire()

	return nil
}

// Get aggregate
func (s *Service) Get(entity Entity) (Aggregate, error) {
	if !entity.ID.IsValid() {
		return nil, err.ErrInvalidID
	}

	return &aggregate{
		Entity:     entity,
		Repository: s.Repository,
		Bus:        s.Bus,
	}, nil
}
