package services

import (
	"context"
	"github.com/dmibod/kanban/shared/domain"

	"gopkg.in/mgo.v2/bson"

	"github.com/dmibod/kanban/shared/kernel"
	"github.com/dmibod/kanban/shared/persistence"
	"github.com/dmibod/kanban/shared/tools/logger"
)

// LanePayload payload
type LanePayload struct {
	Type        string
	Name        string
	Description string
	Layout      string
}

// LaneModel model
type LaneModel struct {
	ID          kernel.ID
	Type        string
	Name        string
	Description string
	Layout      string
}

// LaneReader interface
type LaneReader interface {
	// GetByID get by id
	GetByID(context.Context, kernel.ID) (*LaneModel, error)
	// GetAll lanes
	GetAll(context.Context) ([]*LaneModel, error)
	// GetByLaneID gets lanes by lane id
	GetByLaneID(context.Context, kernel.ID) ([]*LaneModel, error)
	// GetByBoardID gets lanes by board id
	GetByBoardID(context.Context, kernel.ID) ([]*LaneModel, error)
}

// LaneWriter interface
type LaneWriter interface {
	// Create lane
	Create(context.Context, *LanePayload) (*LaneModel, error)
	// Layout lane
	Layout(context.Context, kernel.ID, string) (*LaneModel, error)
	// Name lane
	Name(context.Context, kernel.ID, string) (*LaneModel, error)
	// Describe lane
	Describe(context.Context, kernel.ID, string) (*LaneModel, error)
	// Remove lane
	Remove(context.Context, kernel.ID) error
	// AppendChild to lane
	AppendChild(context.Context, kernel.ID, kernel.ID) error
	// ExcludeChild from lane
	ExcludeChild(context.Context, kernel.ID, kernel.ID) error
}

// LaneService interface
type LaneService interface {
	LaneReader
	LaneWriter
}

type laneService struct {
	logger.Logger
	persistence.BoardRepository
	persistence.LaneRepository
	NotificationService
}

// CreateLaneService instance
func CreateLaneService(s NotificationService, r persistence.LaneRepository, b persistence.BoardRepository, l logger.Logger) LaneService {
	return &laneService{
		Logger:              l,
		BoardRepository:     b,
		LaneRepository:      r,
		NotificationService: s,
	}
}

// GetByID gets lane by id
func (s *laneService) GetByID(ctx context.Context, id kernel.ID) (*LaneModel, error) {
	entity, err := s.LaneRepository.FindLaneByID(ctx, id)
	if err != nil {
		s.Errorln(err)
		return nil, err
	}

	return mapLaneEntityToModel(entity), nil
}

// GetAll lanes
func (s *laneService) GetAll(ctx context.Context) ([]*LaneModel, error) {
	return s.getByCriteria(ctx, nil)
}

// GetByLaneID gets lanes by lane id
func (s *laneService) GetByLaneID(ctx context.Context, laneID kernel.ID) ([]*LaneModel, error) {
	entity, err := s.LaneRepository.FindLaneByID(ctx, laneID)
	if err != nil {
		s.Errorln(err)
		return nil, err
	}

	if len(entity.Children) == 0 {
		return []*LaneModel{}, nil
	}

	return s.getByCriteria(ctx, buildLaneCriteriaByIds(entity.Children))
}

// GetByBoardID gets lanes by board id
func (s *laneService) GetByBoardID(ctx context.Context, boardID kernel.ID) ([]*LaneModel, error) {
	entity, err := s.BoardRepository.FindBoardByID(ctx, boardID)
	if err != nil {
		s.Errorln(err)
		return nil, err
	}

	if len(entity.Children) == 0 {
		return []*LaneModel{}, nil
	}

	return s.getByCriteria(ctx, buildLaneCriteriaByIds(entity.Children))
}

// Create lane
func (s *laneService) Create(ctx context.Context, payload *LanePayload) (*LaneModel, error) {
	return s.createAndGet(ctx, payload.Type, func(aggregate domain.LaneAggregate) error {
		if err := aggregate.Name(payload.Name); err != nil {
			return err
		}
		if err := aggregate.Description(payload.Description); err != nil {
			return err
		}
		return aggregate.Layout(payload.Layout)
	})
}

// Layout board
func (s *laneService) Layout(ctx context.Context, id kernel.ID, layout string) (*LaneModel, error) {
	return s.updateAndGet(ctx, id, func(aggregate domain.LaneAggregate) error {
		return aggregate.Layout(layout)
	})
}

// Name board
func (s *laneService) Name(ctx context.Context, id kernel.ID, name string) (*LaneModel, error) {
	return s.updateAndGet(ctx, id, func(aggregate domain.LaneAggregate) error {
		return aggregate.Name(name)
	})
}

// Describe board
func (s *laneService) Describe(ctx context.Context, id kernel.ID, description string) (*LaneModel, error) {
	return s.updateAndGet(ctx, id, func(aggregate domain.LaneAggregate) error {
		return aggregate.Description(description)
	})
}

// AppendChild to board
func (s *laneService) AppendChild(ctx context.Context, id kernel.ID, childID kernel.ID) error {
	return s.update(ctx, id, func(aggregate domain.LaneAggregate) error {
		return aggregate.AppendChild(childID)
	})
}

// ExcludeChild from board
func (s *laneService) ExcludeChild(ctx context.Context, id kernel.ID, childID kernel.ID) error {
	return s.update(ctx, id, func(aggregate domain.LaneAggregate) error {
		return aggregate.RemoveChild(childID)
	})
}

// Remove lane
func (s *laneService) Remove(ctx context.Context, id kernel.ID) error {
	return s.NotificationService.Execute(func(e domain.EventRegistry) error {
		_, err := domain.DeleteLane(id, s.LaneRepository.DomainRepository(ctx), e)
		return err
	})
}

func (s *laneService) checkCreate(ctx context.Context, aggregate domain.LaneAggregate) error {
	return nil
}

func (s *laneService) create(ctx context.Context, owner string, operation func(domain.LaneAggregate) error) (kernel.ID, error) {
	id := kernel.EmptyID
	err := s.NotificationService.Execute(func(e domain.EventRegistry) error {
		aggregate, err := domain.NewLane(owner, s.LaneRepository.DomainRepository(ctx), e)
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

		err = aggregate.Save()
		if err == nil {
			id = aggregate.GetID()
		}

		return err
	})

	return id, err
}

func (s *laneService) createAndGet(ctx context.Context, owner string, operation func(domain.LaneAggregate) error) (*LaneModel, error) {
	id, err := s.create(ctx, owner, operation)
	if err != nil {
		s.Errorln(err)
		return nil, err
	}

	return s.GetByID(ctx, id)
}

func (s *laneService) checkUpdate(ctx context.Context, aggregate domain.LaneAggregate) error {
	//TODO
	//securityContext := ctx.Value(scKey).(*SecurityContext)
	//if securityContext == nil || !securityContext.IsOwner(aggregate.GetOwner()) { return ErrOperationIsNotAllowed }
	return nil
}

func (s *laneService) update(ctx context.Context, id kernel.ID, operation func(domain.LaneAggregate) error) error {
	return s.NotificationService.Execute(func(e domain.EventRegistry) error {
		aggregate, err := domain.LoadLane(id, s.LaneRepository.DomainRepository(ctx), e)
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

		return aggregate.Save()
	})
}

func (s *laneService) updateAndGet(ctx context.Context, id kernel.ID, operation func(domain.LaneAggregate) error) (*LaneModel, error) {
	err := s.update(ctx, id, operation)
	if err != nil {
		s.Errorln(err)
		return nil, err
	}

	return s.GetByID(ctx, id)
}

func (s *laneService) getByCriteria(ctx context.Context, criteria bson.M) ([]*LaneModel, error) {
	models := []*LaneModel{}
	err := s.LaneRepository.FindLanes(ctx, criteria, func(entity *persistence.LaneEntity) error {
		models = append(models, mapLaneEntityToModel(entity))
		return nil
	})

	if err != nil {
		s.Errorln(err)
		return nil, err
	}

	return models, nil
}

func buildLaneCriteriaByIds(ids []string) bson.M {
	criteria := []bson.M{}

	for _, id := range ids {
		criteria = append(criteria, bson.M{"_id": bson.ObjectIdHex(id)})
	}

	return bson.M{"$or": criteria}
}

func mapLaneEntityToModel(entity *persistence.LaneEntity) *LaneModel {
	return &LaneModel{
		ID:          kernel.ID(entity.ID.Hex()),
		Type:        entity.Kind,
		Name:        entity.Name,
		Description: entity.Description,
		Layout:      entity.Layout,
	}
}
