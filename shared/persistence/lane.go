package persistence

import (
	"context"
	"github.com/dmibod/kanban/shared/persistence/models"

	"github.com/dmibod/kanban/shared/tools/db/mongo"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// LaneListQuery type
type LaneListQuery struct {
	ParentID string
	BoardID  string
}

// Operation to query lane list
func (query LaneListQuery) Operation(ctx context.Context, visitor func(*models.LaneListModel) error) mongo.Operation {
	return func(col *mgo.Collection) error {
		pipline := query.pipeline()
		entity := &models.LaneListModel{}
		return mongo.PipeList(ctx, col, pipline, entity, func(entity interface{}) error {
			if lane, ok := entity.(*models.LaneListModel); ok {
				return visitor(lane)
			}

			return ErrInvalidType
		})
	}
}

func (query LaneListQuery) pipeline() []bson.M {
	matchBoard := bson.M{"$match": mongo.FromID(query.BoardID)}

	projectBoard := bson.M{"$project": bson.M{
		"_id": 0,
		"lanes": bson.M{"$filter": bson.M{
			"input": "$lanes",
			"as":    "lane",
			"cond":  bson.M{"$in": []interface{}{"$$lane._id", "$children"}}}},
	}}

	unwindLanes := bson.M{"$unwind": bson.M{
		"path":                       "$lanes",
		"includeArrayIndex":          "idx",
		"preserveNullAndEmptyArrays": false,
	}}

	//matchLane := bson.M{"$match": bson.M{ "$expr": bson.M{"$in": []interface{}{"$lanes._id", "$children"}}}}

	projectLane := bson.M{"$project": bson.M{
		"_id":         "$lanes._id",
		"kind":        "$lanes.kind",
		"name":        "$lanes.name",
		"description": "$lanes.description",
		"layout":      "$lanes.layout",
		"order":       "$idx",
	}}

	if query.ParentID == "" {
		return []bson.M{matchBoard, projectBoard, unwindLanes /*matchLane,*/, projectLane}
	}

	projectParent := bson.M{"$project": bson.M{
		"lanes": 1,
		"children": bson.M{
			"$reduce": bson.M{
				"input":        "$lanes",
				"initialValue": []bson.ObjectId{},
				"in": bson.M{
					"$cond": []interface{}{
						bson.M{"$eq": []interface{}{"$$this._id", bson.ObjectIdHex(query.ParentID)}},
						"$$this.children",
						"$$value"}}}}}}

	return []bson.M{matchBoard, projectParent, projectBoard, unwindLanes, projectLane}
}

// LaneQuery type
type LaneQuery struct {
	ID      string
	BoardID string
}

// Operation to query lane
func (query LaneQuery) Operation(ctx context.Context, visitor func(*models.Lane) error) mongo.Operation {
	return func(col *mgo.Collection) error {
		pipline := query.pipeline()
		entity := &models.Lane{}
		return mongo.PipeOne(ctx, col, pipline, entity, func(entity interface{}) error {
			if lane, ok := entity.(*models.Lane); ok {
				return visitor(lane)
			}

			return ErrInvalidType
		})
	}
}

func (query LaneQuery) pipeline() []bson.M {
	matchBoard := bson.M{"$match": mongo.FromID(query.BoardID)}

	projectBoard := bson.M{"$project": bson.M{
		"_id": 0,
		"lanes": bson.M{"$filter": bson.M{
			"input": "$lanes",
			"as":    "lane",
			"cond":  bson.M{"$eq": []interface{}{"$$lane._id", bson.ObjectIdHex(query.ID)}}}},
	}}

	unwindLanes := bson.M{"$unwind": bson.M{
		"path":                       "$lanes",
		"includeArrayIndex":          "idx",
		"preserveNullAndEmptyArrays": false,
	}}

	projectLane := bson.M{"$project": bson.M{
		"_id":         "$lanes._id",
		"kind":        "$lanes.kind",
		"name":        "$lanes.name",
		"description": "$lanes.description",
		"layout":      "$lanes.layout",
		"order":       "$idx",
	}}

	return []bson.M{matchBoard, projectBoard, unwindLanes, projectLane}
	/*
		reduceLane := bson.M{"$project": bson.M{
			"_id": 0,
			"lane": bson.M{
				"$reduce": bson.M{
					"input":        "$lanes",
					"initialValue": nil,
					"in": bson.M{
						"$cond": []interface{}{
							bson.M{"$eq": []interface{}{"$$this._id", bson.ObjectIdHex(query.ID)}},
							"$$this",
							"$$value"}}}}}}
		projectLane := bson.M{"$project": bson.M{
			"_id":         "$lane._id",
			"kind":        "$lane.kind",
			"name":        "$lane.name",
			"description": "$lane.description",
			"layout":      "$lane.layout",
		}}

		return []bson.M{matchBoard, reduceLane, projectLane}
	*/
}

// CreateLaneCommand type
type CreateLaneCommand struct {
	BoardID string
	Lane    *models.Lane
}

// Operation to create lane
func (command CreateLaneCommand) Operation(ctx context.Context) mongo.Operation {
	return func(col *mgo.Collection) error {
		return mongo.Update(ctx, col, mongo.FromID(command.BoardID), mongo.AddToSet("lanes", command.Lane))
	}
}

// RemoveLaneCommand type
type RemoveLaneCommand struct {
	BoardID string
	ID      string
}

// Operation to remove lane
func (command RemoveLaneCommand) Operation(ctx context.Context) mongo.Operation {
	return func(col *mgo.Collection) error {
		return mongo.Update(ctx, col, mongo.FromID(command.BoardID), mongo.RemoveFromSet("lanes", mongo.FromID(command.ID)))
	}
}

// UpdateLaneCommand type
type UpdateLaneCommand struct {
	BoardID string
	ID      string
	Field   string
	Value   interface{}
}

// Operation to update lane
func (command UpdateLaneCommand) Operation(ctx context.Context) mongo.Operation {
	return func(col *mgo.Collection) error {
		return mongo.Update(
			ctx,
			col,
			bson.M{"_id": bson.ObjectIdHex(command.BoardID), "lanes._id": bson.ObjectIdHex(command.ID)},
			mongo.Set("lanes.$."+command.Field, command.Value))
	}
}

// AttachToLaneCommand type
type AttachToLaneCommand struct {
	BoardID string
	ID      string
	ChildID string
}

// Operation to attach to lane
func (command AttachToLaneCommand) Operation(ctx context.Context) mongo.Operation {
	return func(col *mgo.Collection) error {
		return mongo.Update(
			ctx,
			col,
			bson.M{"_id": bson.ObjectIdHex(command.BoardID), "lanes._id": bson.ObjectIdHex(command.ID)},
			mongo.AddToSet("lanes.$.children", bson.ObjectIdHex(command.ChildID)))
	}
}

// DetachFromLaneCommand type
type DetachFromLaneCommand struct {
	BoardID string
	ID      string
	ChildID string
}

// Operation to detach from board
func (command DetachFromLaneCommand) Operation(ctx context.Context) mongo.Operation {
	return func(col *mgo.Collection) error {
		return mongo.Update(
			ctx,
			col,
			bson.M{"_id": bson.ObjectIdHex(command.BoardID), "lanes._id": bson.ObjectIdHex(command.ID)},
			mongo.RemoveFromSet("lanes.$.children", bson.ObjectIdHex(command.ChildID)))
	}
}
