package persistence

import (
	"context"

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
func (query LaneListQuery) Operation(ctx context.Context, visitor func(*LaneListModel) error) mongo.Operation {
	return func(col *mgo.Collection) error {
		return mongo.PipeList(ctx, col, query.pipeline(), &LaneListModel{}, func(entity interface{}) error {
			if lane, ok := entity.(*LaneListModel); ok {
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
			"cond":  bson.M{"$in": []string{"$$lane._id", "$children"}}}},
	}}

	unwindLanes := bson.M{"$unwind": bson.M{
		"path":                       "$lanes",
		"includeArrayIndex":          "idx",
		"preserveNullAndEmptyArrays": 1,
	}}

	projectLane := bson.M{"$project": bson.M{
		"_id":         "$lanes._id",
		"kind":        "$lanes.kind",
		"name":        "$lanes.name",
		"description": "$lanes.description",
		"layout":      "$lanes.layout",
		"order":       "$idx",
	}}

	if query.ParentID == "" {
		return []bson.M{matchBoard, projectBoard, unwindLanes, projectLane}
	}

	projectParent := bson.M{"$project": bson.M{
		"lanes": 1,
		"children": bson.M{"$reduce": bson.M{
			"input":        "$lanes",
			"initialValue": []string{},
			"in": bson.M{"$cond": bson.M{
				"if":   bson.M{"$eq": bson.M{"$$this._id": bson.ObjectIdHex(query.ParentID)}},
				"then": "$$this.children",
				"else": "$$value"}}}}}}

	return []bson.M{matchBoard, projectParent, projectBoard, unwindLanes, projectLane}
}

// LaneQuery type
type LaneQuery struct {
	ID      string
	BoardID string
}

// Operation to query lane
func (query LaneQuery) Operation(ctx context.Context, visitor func(*Lane) error) mongo.Operation {
	return func(col *mgo.Collection) error {
		return mongo.PipeOne(ctx, col, query.pipeline(), &Lane{}, func(entity interface{}) error {
			if lane, ok := entity.(*Lane); ok {
				return visitor(lane)
			}

			return ErrInvalidType
		})
	}
}

func (query LaneQuery) pipeline() []bson.M {
	matchBoard := bson.M{"$match": mongo.FromID(query.BoardID)}

	reduceLane := bson.M{"$project": bson.M{
		"_id": 0,
		"lane": bson.M{
			"$reduce": bson.M{
				"input":        "$lanes",
				"initialValue": bson.M{},
				"in": bson.M{"$cond": bson.M{
					"if":   bson.M{"$eq": bson.M{"$$this._id": bson.ObjectIdHex(query.ID)}},
					"then": "$$this",
					"else": "$$value"}}}}}}

	projectLane := bson.M{"$project": bson.M{
		"_id":         "$lanes._id",
		"kind":        "$lanes.kind",
		"name":        "$lanes.name",
		"description": "$lanes.description",
		"layout":      "$lanes.layout",
	}}

	return []bson.M{matchBoard, reduceLane, projectLane}
}

// CreateLaneCommand type
type CreateLaneCommand struct {
	BoardID string
	Lane    *Lane
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
	ID string
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
	ID string
	Field string
	Value interface{}
}

// Operation to update lane
func (command UpdateLaneCommand) Operation(ctx context.Context) mongo.Operation {
	return func(col *mgo.Collection) error {
		return mongo.Update(ctx, col, mongo.FromID(command.BoardID), mongo.Set(command.Field, command.Value))
	}
}

// AttachToLaneCommand type
type AttachToLaneCommand struct {
	BoardID string
	ID string
	ChildID string
}

// Operation to attach to lane
func (command AttachToLaneCommand) Operation(ctx context.Context) mongo.Operation {
	return func(col *mgo.Collection) error {
		return mongo.Update(ctx, col, mongo.FromID(command.BoardID), mongo.AddToSet("children", command.ChildID))
	}
}

// DetachFromLaneCommand type
type DetachFromLaneCommand struct {
	BoardID string
	ID string
	ChildID string
}

// Operation to detach from board
func (command DetachFromLaneCommand) Operation(ctx context.Context) mongo.Operation {
	return func(col *mgo.Collection) error {
		return mongo.Update(ctx, col, mongo.FromID(command.BoardID), mongo.RemoveFromSet("children", command.ChildID))
	}
}

