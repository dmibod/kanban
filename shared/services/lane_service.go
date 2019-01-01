package services

import (
	"context"
	"errors"

	"gopkg.in/mgo.v2/bson"

	"github.com/dmibod/kanban/shared/kernel"
	"github.com/dmibod/kanban/shared/persistence"
	"github.com/dmibod/kanban/shared/tools/db"
	"github.com/dmibod/kanban/shared/tools/logger"
)

// LanePayload payload
type LanePayload struct {
	Name   string
	Layout string
	Type   string
}

// LaneModel model
type LaneModel struct {
	ID     kernel.Id
	Name   string
	Layout string
	Type   string
}

// LaneService interface
type LaneService interface {
	// Create lane
	Create(context.Context, *LanePayload) (kernel.Id, error)
	// Update lane
	Update(context.Context, *LaneModel) (*LaneModel, error)
	// Remove lane
	Remove(context.Context, kernel.Id) error
	// GetByID gets lane by id
	GetByID(context.Context, kernel.Id) (*LaneModel, error)
}

type laneService struct {
	logger.Logger
	db.Repository
}

// Create lane
func (s *laneService) Create(ctx context.Context, payload *LanePayload) (kernel.Id, error) {
	entity := mapPayloadToEntity(payload)
	id, err := s.Repository.Create(ctx, entity)
	if err != nil {
		s.Errorln(err)
		return "", err
	}

	return kernel.Id(id), nil
}

// Update lane
func (s *laneService) Update(ctx context.Context, model *LaneModel) (*LaneModel, error) {
	entity := mapModelToEntity(model)
	err := s.Repository.Update(ctx, entity)
	if err != nil {
		s.Errorln(err)
		return nil, err
	}

	return mapEntityToModel(entity), nil
}

// Remove lane
func (s *laneService) Remove(ctx context.Context, id kernel.Id) error {
	err := s.Repository.Remove(ctx, string(id))
	if err != nil {
		s.Errorln(err)
	}

	return err
}

// GetByID gets lane by id
func (s *laneService) GetByID(ctx context.Context, id kernel.Id) (*LaneModel, error) {
	entity, err := s.Repository.FindByID(ctx, string(id))
	if err != nil {
		s.Errorln(err)
		return nil, err
	}

	lane, ok := entity.(*persistence.LaneEntity)
	if !ok {
		s.Errorf("invalid type %T\n", entity)
		return nil, errors.New("Invalid type")
	}

	return mapEntityToModel(lane), nil
}

func mapEntityToModel(entity *persistence.LaneEntity) *LaneModel {
	return &LaneModel{
		ID:     kernel.Id(entity.ID.Hex()),
		Name:   entity.Name,
		Type:   entity.Type,
		Layout: entity.Layout,
	}
}

func mapModelToEntity(model *LaneModel) *persistence.LaneEntity {
	return &persistence.LaneEntity{
		ID:     bson.ObjectIdHex(string(model.ID)),
		Name:   model.Name,
		Type:   model.Type,
		Layout: model.Layout,
	}
}

func mapPayloadToEntity(model *LanePayload) *persistence.LaneEntity {
	return &persistence.LaneEntity{
		Name:   model.Name,
		Type:   model.Type,
		Layout: model.Layout,
	}
}
