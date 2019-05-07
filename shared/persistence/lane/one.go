package lane

import (
	"context"
	err "github.com/dmibod/kanban/shared/persistence/error"
	"github.com/dmibod/kanban/shared/persistence/models"

	"github.com/dmibod/kanban/shared/tools/db/mongo"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// OneQuery type
type OneQuery struct {
	ID      string
	BoardID string
}

// Operation to query lane
func (query OneQuery) Operation(ctx context.Context, visitor func(*models.Lane) error) mongo.Operation {
	return func(col *mgo.Collection) error {
		pipline := query.pipeline()
		entity := &models.Lane{}
		return mongo.PipeOne(ctx, col, pipline, entity, func(entity interface{}) error {
			if lane, ok := entity.(*models.Lane); ok {
				return visitor(lane)
			}

			return err.ErrInvalidType
		})
	}
}

func (query OneQuery) pipeline() []bson.M {
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
		"children":    "$lanes.children",
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
