package lane

import (
	err "github.com/dmibod/kanban/shared/domain/error"
	"github.com/dmibod/kanban/shared/domain/event"
	"github.com/dmibod/kanban/shared/kernel"
)

// Service - lane domain service
type Service struct {
	Repository
	event.Bus
}

// CreateService - creates lane domain service
func CreateService(repository Repository, bus event.Bus) *Service {
	if repository == nil {
		return nil
	}

	if bus == nil {
		return nil
	}

	return &Service{Repository: repository, Bus: bus}
}

// Create lane
func (s *Service) Create(id kernel.ID, kind string) (*Entity, error) {
	if !id.IsValid() {
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

	if err := s.Repository.Create(&entity); err != nil {
		return nil, err
	}

	s.Bus.Register(CreatedEvent{entity})
	s.Bus.Fire()

	return &entity, nil
}

// Delete lane
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
