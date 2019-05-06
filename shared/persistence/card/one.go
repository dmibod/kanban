package card

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

// Operation to query card
func (query OneQuery) Operation(ctx context.Context, visitor func(*models.Card) error) mongo.Operation {
	return func(col *mgo.Collection) error {
		return mongo.PipeOne(ctx, col, query.pipeline(), &models.Card{}, func(entity interface{}) error {
			if card, ok := entity.(*models.Card); ok {
				return visitor(card)
			}

			return err.ErrInvalidType
		})
	}
}

func (query OneQuery) pipeline() []bson.M {
	matchBoard := bson.M{"$match": mongo.FromID(query.BoardID)}

	projectBoard := bson.M{"$project": bson.M{
		"_id": 0,
		"cards": bson.M{"$filter": bson.M{
			"input": "$cards",
			"as":    "card",
			"cond":  bson.M{"$eq": []interface{}{"$$card._id", bson.ObjectIdHex(query.ID)}}}},
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

	return []bson.M{matchBoard, projectBoard, unwindCards, projectCard}
	/*
		reduceCard := bson.M{"$project": bson.M{
			"_id": 0,
			"card": bson.M{
				"$reduce": bson.M{
					"input":        "$cards",
					"initialValue": nil,
					"in": bson.M{
						"$cond": []interface{}{
							bson.M{"$eq": []interface{}{"$$this._id", bson.ObjectIdHex(query.ID)}},
							"$$this",
							"$$value"}}}}}}

		projectCard := bson.M{"$project": bson.M{
			"_id":         "$card._id",
			"name":        "$card.name",
			"description": "$card.description",
		}}

		return []bson.M{matchBoard, reduceCard, projectCard}
	*/
}
