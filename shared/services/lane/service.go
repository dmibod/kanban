package lane

import (
	"context"

	"github.com/dmibod/kanban/shared/services/notification"

	"github.com/dmibod/kanban/shared/domain/event"
	"github.com/dmibod/kanban/shared/domain/lane"
	"gopkg.in/mgo.v2/bson"

	"github.com/dmibod/kanban/shared/kernel"
	"github.com/dmibod/kanban/shared/persistence"
	"github.com/dmibod/kanban/shared/tools/logger"
)

type service struct {
	logger.Logger
	BoardRepository *persistence.BoardRepository
	LaneRepository  *persistence.LaneRepository
	notification.Service
}

// CreateService instance
func CreateService(s notification.Service, r *persistence.LaneRepository, b *persistence.BoardRepository, l logger.Logger) Service {
	return &service{
		Logger:          l,
		BoardRepository: b,
		LaneRepository:  r,
		Service:         s,
	}
}

// GetByID gets lane by id
func (s *service) GetByID(ctx context.Context, id kernel.ID) (*Model, error) {
	entity, err := s.LaneRepository.FindLaneByID(ctx, id)
	if err != nil {
		s.Errorln(err)
		return nil, err
	}

	return mapPersistentToModel(entity), nil
}

// GetAll lanes
func (s *service) GetAll(ctx context.Context) ([]*ListModel, error) {
	return s.getByCriteria(ctx, nil)
}

// GetByLaneID gets lanes by lane id
func (s *service) GetByLaneID(ctx context.Context, laneID kernel.ID) ([]*ListModel, error) {
	entity, err := s.LaneRepository.FindLaneByID(ctx, laneID)
	if err != nil {
		s.Errorln(err)
		return nil, err
	}

	if len(entity.Children) == 0 {
		return []*ListModel{}, nil
	}

	return s.getByCriteria(ctx, buildLaneCriteriaByIds(entity.Children))
}

// GetByBoardID gets lanes by board id
func (s *service) GetByBoardID(ctx context.Context, boardID kernel.ID) ([]*ListModel, error) {
	entity, err := s.BoardRepository.FindBoardByID(ctx, boardID)
	if err != nil {
		s.Errorln(err)
		return nil, err
	}

	if len(entity.Children) == 0 {
		return []*ListModel{}, nil
	}

	return s.getByCriteria(ctx, buildLaneCriteriaByIds(entity.Children))
}

// Create lane
func (s *service) Create(ctx context.Context, model *CreateModel) (kernel.ID, error) {
	return s.create(ctx, model.Type, func(aggregate lane.Aggregate) error {
		if err := aggregate.Name(model.Name); err != nil {
			return err
		}
		if err := aggregate.Description(model.Description); err != nil {
			return err
		}
		return aggregate.Layout(model.Layout)
	})
}

// Layout board
func (s *service) Layout(ctx context.Context, id kernel.ID, layout string) error {
	return s.update(ctx, id, func(aggregate lane.Aggregate) error {
		return aggregate.Layout(layout)
	})
}

// Name board
func (s *service) Name(ctx context.Context, id kernel.ID, name string) error {
	return s.update(ctx, id, func(aggregate lane.Aggregate) error {
		return aggregate.Name(name)
	})
}

// Describe board
func (s *service) Describe(ctx context.Context, id kernel.ID, description string) error {
	return s.update(ctx, id, func(aggregate lane.Aggregate) error {
		return aggregate.Description(description)
	})
}

// AppendChild to board
func (s *service) AppendChild(ctx context.Context, id kernel.ID, childID kernel.ID) error {
	return s.update(ctx, id, func(aggregate lane.Aggregate) error {
		return aggregate.AppendChild(childID)
	})
}

// ExcludeChild from board
func (s *service) ExcludeChild(ctx context.Context, id kernel.ID, childID kernel.ID) error {
	return s.update(ctx, id, func(aggregate lane.Aggregate) error {
		return aggregate.RemoveChild(childID)
	})
}

// Remove lane
func (s *service) Remove(ctx context.Context, id kernel.ID) error {
	return event.Execute(func(bus event.Bus) error {
		s.Service.Listen(bus)

		domainService := lane.CreateService(s.LaneRepository.GetRepository(ctx), bus)

		return domainService.Delete(lane.Entity{ID: id})
	})
}

func (s *service) checkCreate(ctx context.Context) error {
	return nil
}

func (s *service) create(ctx context.Context, kind string, operation func(lane.Aggregate) error) (kernel.ID, error) {
	if err := s.checkCreate(ctx); err != nil {
		s.Errorln(err)
		return kernel.EmptyID, err
	}

	id := kernel.ID(bson.NewObjectId().Hex())

	err := event.Execute(func(bus event.Bus) error {
		s.Service.Listen(bus)

		domainService := lane.CreateService(s.LaneRepository.GetRepository(ctx), bus)

		entity, err := domainService.Create(id, kind)
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

func (s *service) checkUpdate(ctx context.Context, aggregate lane.Aggregate) error {
	//TODO
	//securityContext := ctx.Value(scKey).(*SecurityContext)
	//if securityContext == nil || !securityContext.IsOwner(aggregate.GetOwner()) { return ErrOperationIsNotAllowed }
	return nil
}

func (s *service) update(ctx context.Context, id kernel.ID, operation func(lane.Aggregate) error) error {
	entity, err := s.LaneRepository.FindLaneByID(ctx, id)
	if err != nil {
		s.Errorln(err)
		return err
	}

	return event.Execute(func(bus event.Bus) error {
		s.Service.Listen(bus)

		domainService := lane.CreateService(s.LaneRepository.GetRepository(ctx), bus)

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

func (s *service) getByCriteria(ctx context.Context, criteria bson.M) ([]*ListModel, error) {
	models := []*ListModel{}
	err := s.LaneRepository.FindLanes(ctx, criteria, func(entity *persistence.LaneEntity) error {
		models = append(models, mapPersistentToListModel(entity))
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

func mapPersistentToModel(entity *persistence.LaneEntity) *Model {
	return &Model{
		ID:          kernel.ID(entity.ID.Hex()),
		Type:        entity.Kind,
		Name:        entity.Name,
		Description: entity.Description,
		Layout:      entity.Layout,
	}
}

func mapPersistentToListModel(entity *persistence.LaneEntity) *ListModel {
	return &ListModel{
		ID:          kernel.ID(entity.ID.Hex()),
		Type:        entity.Kind,
		Name:        entity.Name,
		Description: entity.Description,
		Layout:      entity.Layout,
	}
}

func mapPersistentToDomain(entity *persistence.LaneEntity) lane.Entity {
	children := []kernel.ID{}
	for _, id := range entity.Children {
		children = append(children, kernel.ID(id))
	}
	return lane.Entity{
		ID:          kernel.ID(entity.ID.Hex()),
		Kind:        entity.Kind,
		Name:        entity.Name,
		Description: entity.Description,
		Layout:      entity.Layout,
		Children:    children,
	}
}
