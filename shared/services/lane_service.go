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
	ID     kernel.ID
	Name   string
	Layout string
	Type   string
}

// LaneService interface
type LaneService interface {
	// Create lane
	Create(context.Context, *LanePayload) (*LaneModel, error)
	// Update lane
	Update(context.Context, *LaneModel) (*LaneModel, error)
	// Remove lane
	Remove(context.Context, kernel.ID) error
	// GetByID gets lane by id
	GetByID(context.Context, kernel.ID) (*LaneModel, error)
	// GetAll lanes
	GetAll(context.Context) ([]*LaneModel, error)
	// GetByLaneID gets lanes by lane id
	GetByLaneID(context.Context, kernel.ID) ([]*LaneModel, error)
	// GetByBoardID gets lanes by board id
	GetByBoardID(context.Context, kernel.ID) ([]*LaneModel, error)
	// AppendChild to lane
	AppendChild(context.Context, kernel.ID, kernel.ID) error
	// ExcludeChild from lane
	ExcludeChild(context.Context, kernel.ID, kernel.ID) error
}

type laneService struct {
	logger.Logger
	boardRepository db.Repository
	laneRepository  db.Repository
}

// AppendChild to lane
func (s *laneService) AppendChild(ctx context.Context, id kernel.ID, childID kernel.ID) error {
	entity, err := s.laneRepository.FindByID(ctx, string(id))
	if err != nil {
		s.Errorln(err)
		return err
	}

	lane, ok := entity.(*persistence.LaneEntity)
	if !ok {
		s.Errorf("invalid type %T\n", entity)
		return errors.New("Invalid type")
	}

	child := string(childID)
	for _, val := range lane.Children {
		if val == child {
			return nil
		}
	}

	lane.Children = append(lane.Children, string(childID))

	err = s.laneRepository.Update(ctx, lane)
	if err != nil {
		s.Errorln(err)
		return err
	}

	return nil
}

// ExcludeChild from lane
func (s *laneService) ExcludeChild(ctx context.Context, id kernel.ID, childID kernel.ID) error {
	entity, err := s.laneRepository.FindByID(ctx, string(id))
	if err != nil {
		s.Errorln(err)
		return err
	}

	lane, ok := entity.(*persistence.LaneEntity)
	if !ok {
		s.Errorf("invalid type %T\n", entity)
		return errors.New("Invalid type")
	}

	child := string(childID)
	for idx, val := range lane.Children {
		if val == child {
			lane.Children = append(lane.Children[:idx], lane.Children[idx+1:]...)
			err = s.laneRepository.Update(ctx, lane)
			if err != nil {
				s.Errorln(err)
				return err
			}
			return nil
		}
	}

	return nil
}

// Create lane
func (s *laneService) Create(ctx context.Context, payload *LanePayload) (*LaneModel, error) {
	entity := mapPayloadToEntity(payload)
	id, err := s.laneRepository.Create(ctx, entity)
	if err != nil {
		s.Errorln(err)
		return nil, err
	}

	return s.GetByID(ctx, kernel.ID(id))
}

// Update lane
func (s *laneService) Update(ctx context.Context, model *LaneModel) (*LaneModel, error) {
	entity := mapModelToEntity(model)
	err := s.laneRepository.Update(ctx, entity)
	if err != nil {
		s.Errorln(err)
		return nil, err
	}

	return mapEntityToModel(entity), nil
}

// Remove lane
func (s *laneService) Remove(ctx context.Context, id kernel.ID) error {
	err := s.laneRepository.Remove(ctx, string(id))
	if err != nil {
		s.Errorln(err)
	}

	return err
}

// GetByID gets lane by id
func (s *laneService) GetByID(ctx context.Context, id kernel.ID) (*LaneModel, error) {
	entity, err := s.laneRepository.FindByID(ctx, string(id))
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

// GetAll lanes
func (s *laneService) GetAll(ctx context.Context) ([]*LaneModel, error) {
	models := []*LaneModel{}
	err := s.laneRepository.Find(ctx, nil, func(entity interface{}) error {
		lane, ok := entity.(*persistence.LaneEntity)
		if !ok {
			s.Errorf("invalid type %T\n", entity)
			return errors.New("Invalid type")
		}

		models = append(models, mapEntityToModel(lane))

		return nil
	})

	if err != nil {
		s.Errorln(err)
		return nil, err
	}

	return models, nil
}

// GetByLaneID gets lanes by lane id
func (s *laneService) GetByLaneID(ctx context.Context, laneID kernel.ID) ([]*LaneModel, error) {
	laneEntity, err := s.laneRepository.FindByID(ctx, string(laneID))
	if err != nil {
		s.Errorln(err)
		return nil, err
	}

	lane, ok := laneEntity.(*persistence.LaneEntity)
	if !ok {
		s.Errorf("invalid type %T\n", laneEntity)
		return nil, errors.New("Invalid type")
	}

	if len(lane.Children) == 0 {
		return []*LaneModel{}, nil
	}

	criteria := []bson.M{}

	for _, id := range lane.Children {
		criteria = append(criteria, bson.M{"_id": bson.ObjectIdHex(id)})
	}

	models := []*LaneModel{}
	err = s.laneRepository.Find(ctx, bson.M{"$or": criteria}, func(entity interface{}) error {
		lane, ok := entity.(*persistence.LaneEntity)
		if !ok {
			s.Errorf("invalid type %T\n", entity)
			return errors.New("Invalid type")
		}

		model := mapEntityToModel(lane)

		models = append(models, model)

		return nil
	})

	if err != nil {
		s.Errorln(err)
		return nil, err
	}

	return models, nil
}

// GetByBoardID gets lanes by board id
func (s *laneService) GetByBoardID(ctx context.Context, boardID kernel.ID) ([]*LaneModel, error) {
	boardEntity, err := s.boardRepository.FindByID(ctx, string(boardID))
	if err != nil {
		s.Errorln(err)
		return nil, err
	}

	board, ok := boardEntity.(*persistence.BoardEntity)
	if !ok {
		s.Errorf("invalid type %T\n", boardEntity)
		return nil, errors.New("Invalid type")
	}

	if len(board.Children) == 0 {
		return []*LaneModel{}, nil
	}

	criteria := []bson.M{}

	for _, id := range board.Children {
		criteria = append(criteria, bson.M{"_id": bson.ObjectIdHex(id)})
	}

	models := []*LaneModel{}
	err = s.laneRepository.Find(ctx, bson.M{"$or": criteria}, func(entity interface{}) error {
		lane, ok := entity.(*persistence.LaneEntity)
		if !ok {
			s.Errorf("invalid type %T\n", entity)
			return errors.New("Invalid type")
		}

		model := mapEntityToModel(lane)

		models = append(models, model)

		return nil
	})

	if err != nil {
		s.Errorln(err)
		return nil, err
	}

	return models, nil
}

func mapEntityToModel(entity *persistence.LaneEntity) *LaneModel {
	return &LaneModel{
		ID:     kernel.ID(entity.ID.Hex()),
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
