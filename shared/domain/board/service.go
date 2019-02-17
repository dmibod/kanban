package board

import (
	err "github.com/dmibod/kanban/shared/domain/error"
	"github.com/dmibod/kanban/shared/domain/event"
	"github.com/dmibod/kanban/shared/kernel"
)

// Service - board domain service
type Service struct {
	event.Bus
}

// CreateService - creates board domain service
func CreateService(bus event.Bus) *Service {
	if bus == nil {
		return nil
	}

	return &Service{Bus: bus}
}

// Create board
func (s *Service) Create(id kernel.ID, owner string) (*Entity, error) {
	if !id.IsValid() {
		return nil, err.ErrInvalidID
	}

	if owner == "" {
		return nil, err.ErrInvalidArgument
	}

	entity := Entity{
		ID:       id,
		Owner:    owner,
		Layout:   kernel.VLayout,
		Shared:   false,
		Children: []kernel.ID{},
	}

	s.Bus.Register(CreatedEvent{entity})
	s.Bus.Fire()

	return &entity, nil
}

// Delete board
func (s *Service) Delete(entity Entity) error {
	if !entity.ID.IsValid() {
		return err.ErrInvalidID
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
		Entity: entity,
		Bus:    s.Bus,
	}, nil
}
