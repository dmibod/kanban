package card

import (
	err "github.com/dmibod/kanban/shared/domain/error"
	"github.com/dmibod/kanban/shared/domain/event"
	"github.com/dmibod/kanban/shared/kernel"
)

// Service - card domain service
type Service struct {
	event.Bus
}

// CreateService - creates card domain service
func CreateService(bus event.Bus) *Service {
	if bus == nil {
		return nil
	}

	return &Service{Bus: bus}
}

// Create card
func (s *Service) Create(id kernel.MemberID) (*Entity, error) {
	if !id.SetID.IsValid() || !id.ID.IsValid() {
		return nil, err.ErrInvalidID
	}

	entity := Entity{ID: id}

	s.Bus.Register(CreatedEvent{entity})

	return &entity, nil
}

// Delete card
func (s *Service) Delete(entity Entity) error {
	if !entity.ID.SetID.IsValid() || !entity.ID.ID.IsValid() {
		return err.ErrInvalidID
	}

	s.Bus.Register(DeletedEvent{entity})

	return nil
}

// Get aggregate
func (s *Service) Get(entity Entity) (Aggregate, error) {
	if !entity.ID.SetID.IsValid() || !entity.ID.ID.IsValid() {
		return nil, err.ErrInvalidID
	}

	return &aggregate{
		Entity: entity,
		Bus:    s.Bus,
	}, nil
}
