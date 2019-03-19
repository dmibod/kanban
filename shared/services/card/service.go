package card

import (
	"context"

	"github.com/dmibod/kanban/shared/services/notification"

	"github.com/dmibod/kanban/shared/domain/card"
	"github.com/dmibod/kanban/shared/domain/event"

	"gopkg.in/mgo.v2/bson"

	"github.com/dmibod/kanban/shared/kernel"
	"github.com/dmibod/kanban/shared/persistence"
	"github.com/dmibod/kanban/shared/tools/logger"
)

type service struct {
	logger.Logger
	CardRepository *persistence.CardRepository
	LaneRepository *persistence.LaneRepository
	notification.Service
}

// CreateService instance
func CreateService(s notification.Service, c *persistence.CardRepository, r *persistence.LaneRepository, l logger.Logger) Service {
	return &service{
		Logger:         l,
		CardRepository: c,
		LaneRepository: r,
		Service:        s,
	}
}

// GetByID gets card by id
func (s *service) GetByID(ctx context.Context, id kernel.ID) (*Model, error) {
	entity, err := s.CardRepository.FindCardByID(ctx, id)
	if err != nil {
		s.Errorln(err)
		return nil, err
	}

	return mapPersistentToModel(entity), nil
}

// GetAll cards
func (s *service) GetAll(ctx context.Context) ([]*Model, error) {
	return s.getByCriteria(ctx, nil)
}

// GetByLaneID gets cards by lane id
func (s *service) GetByLaneID(ctx context.Context, laneID kernel.ID) ([]*Model, error) {
	entity, err := s.LaneRepository.FindLaneByID(ctx, laneID)
	if err != nil {
		s.Errorln(err)
		return nil, err
	}

	if len(entity.Children) == 0 {
		return []*Model{}, nil
	}

	return s.getByCriteria(ctx, buildCardCriteriaByIds(entity.Children))
}

// Create card
func (s *service) Create(ctx context.Context, model *CreateModel) (kernel.ID, error) {
	return s.create(ctx, func(aggregate card.Aggregate) error {
		if err := aggregate.Name(model.Name); err != nil {
			return err
		}
		return aggregate.Description(model.Description)
	})
}

// Name board
func (s *service) Name(ctx context.Context, id kernel.ID, name string) error {
	return s.update(ctx, id, func(aggregate card.Aggregate) error {
		return aggregate.Name(name)
	})
}

// Describe board
func (s *service) Describe(ctx context.Context, id kernel.ID, description string) error {
	return s.update(ctx, id, func(aggregate card.Aggregate) error {
		return aggregate.Description(description)
	})
}

// Remove card
func (s *service) Remove(ctx context.Context, id kernel.ID) error {
	return event.Execute(func(bus event.Bus) error {
		s.Service.Listen(bus)

		domainService := card.CreateService(s.CardRepository.GetRepository(ctx), bus)

		return domainService.Delete(card.Entity{ID: id})
	})
}

func (s *service) checkCreate(ctx context.Context) error {
	return nil
}

func (s *service) create(ctx context.Context, operation func(card.Aggregate) error) (kernel.ID, error) {
	if err := s.checkCreate(ctx); err != nil {
		s.Errorln(err)
		return kernel.EmptyID, err
	}

	id := kernel.ID(bson.NewObjectId().Hex())

	err := event.Execute(func(bus event.Bus) error {
		s.Service.Listen(bus)

		domainService := card.CreateService(s.CardRepository.GetRepository(ctx), bus)

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

	return id, err
}

func (s *service) checkUpdate(ctx context.Context, aggregate card.Aggregate) error {
	//TODO
	//securityContext := ctx.Value(scKey).(*SecurityContext)
	//if securityContext == nil || !securityContext.IsOwner(aggregate.GetOwner()) { return ErrOperationIsNotAllowed }
	return nil
}

func (s *service) update(ctx context.Context, id kernel.ID, operation func(card.Aggregate) error) error {
	entity, err := s.CardRepository.FindCardByID(ctx, id)
	if err != nil {
		s.Errorln(err)
		return err
	}

	return event.Execute(func(bus event.Bus) error {
		s.Service.Listen(bus)

		domainService := card.CreateService(s.CardRepository.GetRepository(ctx), bus)

		aggregate, err := domainService.Get(mapPersistentToDomain(entity))
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
		if err == nil {
			return aggregate.Save()
		}

		return err
	})
}

func (s *service) getByCriteria(ctx context.Context, criteria bson.M) ([]*Model, error) {
	models := []*Model{}
	err := s.CardRepository.FindCards(ctx, criteria, func(entity *persistence.CardEntity) error {
		models = append(models, mapPersistentToModel(entity))
		return nil
	})

	if err != nil {
		s.Errorln(err)
		return nil, err
	}

	return models, nil
}

func buildCardCriteriaByIds(ids []string) bson.M {
	criteria := []bson.M{}

	for _, id := range ids {
		criteria = append(criteria, bson.M{"_id": bson.ObjectIdHex(id)})
	}

	return bson.M{"$or": criteria}
}

func mapPersistentToModel(entity *persistence.CardEntity) *Model {
	return &Model{
		ID:          kernel.ID(entity.ID.Hex()),
		Name:        entity.Name,
		Description: entity.Description,
	}
}

func mapPersistentToDomain(entity *persistence.CardEntity) card.Entity {
	return card.Entity{
		ID:          kernel.ID(entity.ID.Hex()),
		Name:        entity.Name,
		Description: entity.Description,
	}
}
