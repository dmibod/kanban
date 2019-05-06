package card

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
	LaneID  string
	BoardID string
}

// Operation to query card list
func (query ListQuery) Operation(ctx context.Context, visitor func(*models.Card) error) mongo.Operation {
	return func(col *mgo.Collection) error {
		return mongo.PipeList(ctx, col, query.pipeline(), &models.Card{}, func(entity interface{}) error {
			if card, ok := entity.(*models.Card); ok {
				return visitor(card)
			}

			return err.ErrInvalidType
		})
	}
}

func (query ListQuery) pipeline() []bson.M {
	matchBoard := bson.M{"$match": mongo.FromID(query.BoardID)}

	projectLane := bson.M{"$project": bson.M{
		"cards": 1,
		"children": bson.M{"$reduce": bson.M{
			"input":        "$lanes",
			"initialValue": []bson.ObjectId{},
			"in": bson.M{
				"$cond": []interface{}{
					bson.M{"$eq": []interface{}{"$$this._id", bson.ObjectIdHex(query.LaneID)}},
					"$$this.children",
					"$$value"}}}}}}

	projectBoard := bson.M{"$project": bson.M{
		"_id": 0,
		"cards": bson.M{"$filter": bson.M{
			"input": "$cards",
			"as":    "card",
			"cond":  bson.M{"$in": []interface{}{"$$card._id", "$children"}}}},
	}}

	unwindCards := bson.M{"$unwind": bson.M{
		"path":                       "$cards",
		"includeArrayIndex":          "idx",
		"preserveNullAndEmptyArrays": false,
	}}

	projectCard := bson.M{"$project": bson.M{
		"_id":         "$cards._id",
		"name":        "$cards.name",
		"description": "$cards.description",
		"order":       "$idx",
	}}

	return []bson.M{matchBoard, projectLane, projectBoard, unwindCards, projectCard}
}
