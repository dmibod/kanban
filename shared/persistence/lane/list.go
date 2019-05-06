package lane

import (
	"context"
	err "github.com/dmibod/kanban/shared/persistence/error"
	"github.com/dmibod/kanban/shared/persistence/models"

	"github.com/dmibod/kanban/shared/tools/db/mongo"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// ListQuery type
type ListQuery struct {
	ParentID string
	BoardID  string
}

// Operation to query lane list
func (query ListQuery) Operation(ctx context.Context, visitor func(*models.LaneListModel) error) mongo.Operation {
	return func(col *mgo.Collection) error {
		pipline := query.pipeline()
		entity := &models.LaneListModel{}
		return mongo.PipeList(ctx, col, pipline, entity, func(entity interface{}) error {
			if lane, ok := entity.(*models.LaneListModel); ok {
				return visitor(lane)
			}

			return err.ErrInvalidType
		})
	}
}

func (query ListQuery) pipeline() []bson.M {
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
