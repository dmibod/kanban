package card

import (
	"context"
	persistence "github.com/dmibod/kanban/shared/persistence/card"

	tx "github.com/dmibod/kanban/shared/services/event"

	"github.com/dmibod/kanban/shared/domain/card"
	"github.com/dmibod/kanban/shared/domain/event"

	"gopkg.in/mgo.v2/bson"

	"github.com/dmibod/kanban/shared/kernel"
	"github.com/dmibod/kanban/shared/persistence/models"
	"github.com/dmibod/kanban/shared/tools/logger"
)

type service struct {
	logger.Logger
	persistence.Repository
	tx.Service
}

// CreateService instance
func CreateService(s tx.Service, r persistence.Repository, l logger.Logger) Service {
	return &service{
		Logger:     l,
		Repository: r,
		Service:    s,
	}
}

// GetByID gets card by id
func (s *service) GetByID(ctx context.Context, id kernel.MemberID) (*Model, error) {
	var model *Model
	if err := s.Repository.FindByID(ctx, id, func(entity *models.Card) error {
		model = mapPersistentToModel(entity)
		return nil
	}); err != nil {
		s.Errorln(err)
		return nil, err
	}

	return model, nil
}

// GetByLaneID gets cards by lane id
func (s *service) GetByLaneID(ctx context.Context, laneID kernel.MemberID) ([]*Model, error) {
	cards := []*Model{}
	err := s.Repository.FindByParent(ctx, laneID, func(entity *models.Card) error {
		cards = append(cards, mapPersistentToModel(entity))
		return nil
	})

	if err != nil {
		s.Errorln(err)
		return nil, err
	}

	return cards, nil
}

// Create card
func (s *service) Create(ctx context.Context, boardID kernel.ID, model *CreateModel) (kernel.ID, error) {
	return s.create(ctx, boardID, func(aggregate card.Aggregate) error {
		if err := aggregate.Name(model.Name); err != nil {
			return err
		}
		return aggregate.Description(model.Description)
	})
}

// Name card
func (s *service) Name(ctx context.Context, id kernel.MemberID, name string) error {
	return s.update(ctx, id, func(aggregate card.Aggregate) error {
		return aggregate.Name(name)
	})
}

// Describe card
func (s *service) Describe(ctx context.Context, id kernel.MemberID, description string) error {
	return s.update(ctx, id, func(aggregate card.Aggregate) error {
		return aggregate.Description(description)
	})
}

// Remove card
func (s *service) Remove(ctx context.Context, id kernel.MemberID) error {
	var model *card.Entity
	if err := s.Repository.FindByID(ctx, id, func(entity *models.Card) error {
		model = mapPersistentToDomain(id.SetID, entity)
		return nil
	}); err != nil {
		s.Errorln(err)
		return err
	}

	return s.Service.Execute(ctx, func(bus event.Bus) error {
		return card.CreateService(bus).Delete(*model)
	})
}

func (s *service) checkCreate(ctx context.Context) error {
	return nil
}

func (s *service) create(ctx context.Context, boardID kernel.ID, operation func(card.Aggregate) error) (kernel.ID, error) {
	if err := s.checkCreate(ctx); err != nil {
		s.Errorln(err)
		return kernel.EmptyID, err
	}

	id := kernel.MemberID{SetID: boardID, ID: kernel.ID(bson.NewObjectId().Hex())}

	err := s.Service.Execute(ctx, func(bus event.Bus) error {
		domainService := card.CreateService(bus)

		entity, err := domainService.Create(id)
		if err != nil {
			s.Errorln(err)
			return err
		}

		aggregate, err := domainService.Get(*entity)
		if err != nil {
			s.Errorln(err)
			return err
		}

		err = operation(aggregate)
		if err != nil {
			s.Errorln(err)
			return err
		}

		return err
	})

	return id.ID, err
}

func (s *service) checkUpdate(ctx context.Context, aggregate card.Aggregate) error {
	//TODO
	//securityContext := ctx.Value(scKey).(*SecurityContext)
	//if securityContext == nil || !securityContext.IsOwner(aggregate.GetOwner()) { return ErrOperationIsNotAllowed }
	return nil
}

func (s *service) update(ctx context.Context, id kernel.MemberID, operation func(card.Aggregate) error) error {
	var model *card.Entity
	if err := s.Repository.FindByID(ctx, id, func(entity *models.Card) error {
		model = mapPersistentToDomain(id.SetID, entity)
		return nil
	}); err != nil {
		s.Errorln(err)
		return err
	}

	return s.Service.Execute(ctx, func(bus event.Bus) error {
		aggregate, err := card.CreateService(bus).Get(*model)
		if err != nil {
			s.Errorln(err)
			return err
		}

		err = s.checkUpdate(ctx, aggregate)
		if err != nil {
			s.Errorln(err)
			return err
		}

		err = operation(aggregate)
		if err != nil {
			s.Errorln(err)
			return err
		}

		return nil
	})
}

func mapPersistentToModel(entity *models.Card) *Model {
	return &Model{
		ID:          kernel.ID(entity.ID.Hex()),
		Name:        entity.Name,
		Description: entity.Description,
	}
}

func mapPersistentToDomain(boardID kernel.ID, entity *models.Card) *card.Entity {
	return &card.Entity{
		ID:          kernel.MemberID{ID: kernel.ID(entity.ID.Hex()), SetID: boardID},
		Name:        entity.Name,
		Description: entity.Description,
	}
}
