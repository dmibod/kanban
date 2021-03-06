package lane

import (
	"context"

	tx "github.com/dmibod/kanban/shared/services/event"

	"github.com/dmibod/kanban/shared/domain/event"
	"github.com/dmibod/kanban/shared/domain/lane"
	"gopkg.in/mgo.v2/bson"

	"github.com/dmibod/kanban/shared/kernel"
	persistence "github.com/dmibod/kanban/shared/persistence/lane"
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

// GetByID gets lane by id
func (s *service) GetByID(ctx context.Context, id kernel.MemberID) (*Model, error) {
	var model *Model
	if err := s.Repository.FindByID(ctx, id, func(entity *models.Lane) error {
		model = mapPersistentToModel(entity)
		return nil
	}); err != nil {
		s.Errorln(err)
		return nil, err
	}

	return model, nil
}

// GetByBoardID gets lanes by board id
func (s *service) GetByBoardID(ctx context.Context, boardID kernel.ID) ([]*ListModel, error) {
	lanes := []*ListModel{}
	err := s.Repository.FindByParent(ctx, boardID.WithID(kernel.EmptyID), func(entity *models.LaneListModel) error {
		lanes = append(lanes, mapPersistentToListModel(entity))
		return nil
	})

	if err != nil {
		s.Errorln(err)
		return nil, err
	}

	return lanes, nil
}

// GetByLaneID gets lanes by lane id
func (s *service) GetByLaneID(ctx context.Context, laneID kernel.MemberID) ([]*ListModel, error) {
	lanes := []*ListModel{}
	err := s.Repository.FindByParent(ctx, laneID, func(entity *models.LaneListModel) error {
		lanes = append(lanes, mapPersistentToListModel(entity))
		return nil
	})

	if err != nil {
		s.Errorln(err)
		return nil, err
	}

	return lanes, nil
}

// Create lane
func (s *service) Create(ctx context.Context, boardID kernel.ID, model *CreateModel) (kernel.ID, error) {
	return s.create(ctx, boardID, model.Type, func(aggregate lane.Aggregate) error {
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
func (s *service) Layout(ctx context.Context, id kernel.MemberID, layout string) error {
	return s.update(ctx, id, func(aggregate lane.Aggregate) error {
		return aggregate.Layout(layout)
	})
}

// Name board
func (s *service) Name(ctx context.Context, id kernel.MemberID, name string) error {
	return s.update(ctx, id, func(aggregate lane.Aggregate) error {
		return aggregate.Name(name)
	})
}

// Describe board
func (s *service) Describe(ctx context.Context, id kernel.MemberID, description string) error {
	return s.update(ctx, id, func(aggregate lane.Aggregate) error {
		return aggregate.Description(description)
	})
}

// AppendChild to board
func (s *service) AppendChild(ctx context.Context, id kernel.MemberID, childID kernel.ID) error {
	return s.update(ctx, id, func(aggregate lane.Aggregate) error {
		return aggregate.AppendChild(childID)
	})
}

// ExcludeChild from board
func (s *service) ExcludeChild(ctx context.Context, id kernel.MemberID, childID kernel.ID) error {
	return s.update(ctx, id, func(aggregate lane.Aggregate) error {
		return aggregate.RemoveChild(childID)
	})
}

// Remove lane
func (s *service) Remove(ctx context.Context, id kernel.MemberID) error {
	var entity *lane.Entity
	if err := s.Repository.FindByID(ctx, id, func(lane *models.Lane) error {
		entity = mapPersistentToDomain(id.SetID, lane)
		return nil
	}); err != nil {
		s.Errorln(err)
		return err
	}

	return s.Service.Execute(ctx, func(bus event.Bus) error {
		return lane.CreateService(bus).Delete(*entity)
	})
}

func (s *service) checkCreate(ctx context.Context) error {
	return nil
}

func (s *service) create(ctx context.Context, boardID kernel.ID, kind string, operation func(lane.Aggregate) error) (kernel.ID, error) {
	if err := s.checkCreate(ctx); err != nil {
		s.Errorln(err)
		return kernel.EmptyID, err
	}

	id := kernel.MemberID{SetID: boardID, ID: kernel.ID(bson.NewObjectId().Hex())}

	err := s.Service.Execute(ctx, func(bus event.Bus) error {
		domainService := lane.CreateService(bus)

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

		return nil
	})

	return id.ID, err
}

func (s *service) checkUpdate(ctx context.Context, aggregate lane.Aggregate) error {
	//TODO
	//securityContext := ctx.Value(scKey).(*SecurityContext)
	//if securityContext == nil || !securityContext.IsOwner(aggregate.GetOwner()) { return ErrOperationIsNotAllowed }
	return nil
}

func (s *service) update(ctx context.Context, id kernel.MemberID, operation func(lane.Aggregate) error) error {
	var entity *lane.Entity
	if err := s.Repository.FindByID(ctx, id, func(lane *models.Lane) error {
		entity = mapPersistentToDomain(id.SetID, lane)
		return nil
	}); err != nil {
		s.Errorln(err)
		return err
	}

	return s.Service.Execute(ctx, func(bus event.Bus) error {
		aggregate, err := lane.CreateService(bus).Get(*entity)
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

func mapPersistentToModel(entity *models.Lane) *Model {
	children := []kernel.ID{}
	for _, id := range entity.Children {
		children = append(children, kernel.ID(id.Hex()))
	}
	return &Model{
		ID:          kernel.ID(entity.ID.Hex()),
		Type:        entity.Kind,
		Name:        entity.Name,
		Description: entity.Description,
		Layout:      entity.Layout,
		Children:    children,
	}
}

func mapPersistentToListModel(entity *models.LaneListModel) *ListModel {
	return &ListModel{
		ID:          kernel.ID(entity.ID.Hex()),
		Type:        entity.Kind,
		Name:        entity.Name,
		Description: entity.Description,
		Layout:      entity.Layout,
	}
}

func mapPersistentToDomain(boardID kernel.ID, entity *models.Lane) *lane.Entity {
	children := []kernel.ID{}
	for _, id := range entity.Children {
		children = append(children, kernel.ID(id.Hex()))
	}
	return &lane.Entity{
		ID:          kernel.MemberID{ID: kernel.ID(entity.ID.Hex()), SetID: boardID},
		Kind:        entity.Kind,
		Name:        entity.Name,
		Description: entity.Description,
		Layout:      entity.Layout,
		Children:    children,
	}
}
