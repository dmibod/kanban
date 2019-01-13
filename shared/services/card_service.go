package services

import (
	"context"

	"github.com/dmibod/kanban/shared/domain/card"
	"github.com/dmibod/kanban/shared/domain/event"

	"gopkg.in/mgo.v2/bson"

	"github.com/dmibod/kanban/shared/kernel"
	"github.com/dmibod/kanban/shared/persistence"
	"github.com/dmibod/kanban/shared/tools/logger"
)

// CardPayload represents card fields without id
type CardPayload struct {
	Name        string
	Description string
}

// CardModel represents card at service layer
type CardModel struct {
	ID          kernel.ID
	Name        string
	Description string
}

// CardReader interface
type CardReader interface {
	// GetByID gets card by id
	GetByID(context.Context, kernel.ID) (*CardModel, error)
	// GetAll cards
	GetAll(context.Context) ([]*CardModel, error)
	// GetByLaneID gets cards by lane id
	GetByLaneID(context.Context, kernel.ID) ([]*CardModel, error)
}

// CardWriter interface
type CardWriter interface {
	// Create card
	Create(context.Context, *CardPayload) (*CardModel, error)
	// Name card
	Name(context.Context, kernel.ID, string) (*CardModel, error)
	// Describe card
	Describe(context.Context, kernel.ID, string) (*CardModel, error)
	// Remove card
	Remove(context.Context, kernel.ID) error
}

// CardService interface
type CardService interface {
	CardReader
	CardWriter
}

type cardService struct {
	logger.Logger
	CardRepository *persistence.CardRepository
	LaneRepository *persistence.LaneRepository
	NotificationService
}

// CreateCardService instance
func CreateCardService(s NotificationService, c *persistence.CardRepository, r *persistence.LaneRepository, l logger.Logger) CardService {
	return &cardService{
		Logger:              l,
		CardRepository:      c,
		LaneRepository:      r,
		NotificationService: s,
	}
}

// GetByID gets card by id
func (s *cardService) GetByID(ctx context.Context, id kernel.ID) (*CardModel, error) {
	entity, err := s.CardRepository.FindCardByID(ctx, id)
	if err != nil {
		s.Errorln(err)
		return nil, err
	}

	return mapCardEntityToModel(entity), nil
}

// GetAll cards
func (s *cardService) GetAll(ctx context.Context) ([]*CardModel, error) {
	return s.getByCriteria(ctx, nil)
}

// GetByLaneID gets cards by lane id
func (s *cardService) GetByLaneID(ctx context.Context, laneID kernel.ID) ([]*CardModel, error) {
	entity, err := s.LaneRepository.FindLaneByID(ctx, laneID)
	if err != nil {
		s.Errorln(err)
		return nil, err
	}

	if len(entity.Children) == 0 {
		return []*CardModel{}, nil
	}

	return s.getByCriteria(ctx, buildCardCriteriaByIds(entity.Children))
}

// Create card
func (s *cardService) Create(ctx context.Context, payload *CardPayload) (*CardModel, error) {
	return s.createAndGet(ctx, func(aggregate card.Aggregate) error {
		if err := aggregate.Name(payload.Name); err != nil {
			return err
		}
		return aggregate.Description(payload.Description)
	})
}

// Name board
func (s *cardService) Name(ctx context.Context, id kernel.ID, name string) (*CardModel, error) {
	return s.updateAndGet(ctx, id, func(aggregate card.Aggregate) error {
		return aggregate.Name(name)
	})
}

// Describe board
func (s *cardService) Describe(ctx context.Context, id kernel.ID, description string) (*CardModel, error) {
	return s.updateAndGet(ctx, id, func(aggregate card.Aggregate) error {
		return aggregate.Description(description)
	})
}

// Remove card
func (s *cardService) Remove(ctx context.Context, id kernel.ID) error {
	return event.Execute(func(bus event.Bus) error {
		s.NotificationService.Listen(bus)
		return card.Delete(card.Entity{ID: id}, bus)
	})
}

func (s *cardService) checkCreate(ctx context.Context, aggregate card.Aggregate) error {
	return nil
}

func (s *cardService) create(ctx context.Context, operation func(card.Aggregate) error) (kernel.ID, error) {
	id := kernel.ID(bson.NewObjectId().Hex())

	err := event.Execute(func(bus event.Bus) error {
		s.NotificationService.Listen(bus)
		entity, err := card.Create(id, bus)
		if err != nil {
			s.Errorln(err)
			return err
		}

		aggregate, err := card.New(*entity, bus)
		if err != nil {
			s.Errorln(err)
			return err
		}

		err = s.checkCreate(ctx, aggregate)
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

	return id, err
}

func (s *cardService) createAndGet(ctx context.Context, operation func(card.Aggregate) error) (*CardModel, error) {
	id, err := s.create(ctx, operation)
	if err != nil {
		s.Errorln(err)
		return nil, err
	}

	return s.GetByID(ctx, id)
}

func (s *cardService) checkUpdate(ctx context.Context, aggregate card.Aggregate) error {
	//TODO
	//securityContext := ctx.Value(scKey).(*SecurityContext)
	//if securityContext == nil || !securityContext.IsOwner(aggregate.GetOwner()) { return ErrOperationIsNotAllowed }
	return nil
}

func (s *cardService) update(ctx context.Context, id kernel.ID, operation func(card.Aggregate) error) error {
	return event.Execute(func(bus event.Bus) error {
		s.NotificationService.Listen(bus)
		aggregate, err := card.New(card.Entity{ID: id}, bus)
		if err != nil {
			s.Errorln(err)
			return err
		}

		err = s.checkUpdate(ctx, aggregate)
		if err != nil {
			s.Errorln(err)
			return err
		}

		return operation(aggregate)
	})
}

func (s *cardService) updateAndGet(ctx context.Context, id kernel.ID, operation func(card.Aggregate) error) (*CardModel, error) {
	err := s.update(ctx, id, operation)
	if err != nil {
		s.Errorln(err)
		return nil, err
	}

	return s.GetByID(ctx, id)
}

func (s *cardService) getByCriteria(ctx context.Context, criteria bson.M) ([]*CardModel, error) {
	models := []*CardModel{}
	err := s.CardRepository.FindCards(ctx, criteria, func(entity *persistence.CardEntity) error {
		models = append(models, mapCardEntityToModel(entity))
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

func mapCardEntityToModel(entity *persistence.CardEntity) *CardModel {
	return &CardModel{
		ID:          kernel.ID(entity.ID.Hex()),
		Name:        entity.Name,
		Description: entity.Description,
	}
}
