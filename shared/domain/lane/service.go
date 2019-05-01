package lane

import (
	err "github.com/dmibod/kanban/shared/domain/error"
	"github.com/dmibod/kanban/shared/domain/event"
	"github.com/dmibod/kanban/shared/kernel"
)

// Service - lane domain service
type Service struct {
	event.Bus
}

// CreateService - creates lane domain service
func CreateService(bus event.Bus) *Service {
	if bus == nil {
		return nil
	}

	return &Service{Bus: bus}
}

// Create lane
func (s *Service) Create(id kernel.MemberID, kind string) (*Entity, error) {
	if !id.SetID.IsValid() || !id.ID.IsValid() {
		return nil, err.ErrInvalidID
	}

	if kind != kernel.LKind && kind != kernel.CKind {
		return nil, err.ErrInvalidArgument
	}

	layout := kernel.VLayout

	if kind == kernel.CKind {
		layout = ""
	}

	entity := Entity{
		ID:       id,
		Kind:     kind,
		Layout:   layout,
		Children: []kernel.ID{},
	}

	s.Bus.Register(CreatedEvent{entity})

	return &entity, nil
}

// Delete lane
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
