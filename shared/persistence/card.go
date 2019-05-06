package persistence

import (
	"context"
	"github.com/dmibod/kanban/shared/persistence/models"

	"github.com/dmibod/kanban/shared/tools/db/mongo"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// CardListQuery type
type CardListQuery struct {
	LaneID  string
	BoardID string
}

// Operation to query card list
func (query CardListQuery) Operation(ctx context.Context, visitor func(*models.Card) error) mongo.Operation {
	return func(col *mgo.Collection) error {
		return mongo.PipeList(ctx, col, query.pipeline(), &models.Card{}, func(entity interface{}) error {
			if card, ok := entity.(*models.Card); ok {
				return visitor(card)
			}

			return ErrInvalidType
		})
	}
}

func (query CardListQuery) pipeline() []bson.M {
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

// CardQuery type
type CardQuery struct {
	ID      string
	BoardID string
}

// Operation to query card
func (query CardQuery) Operation(ctx context.Context, visitor func(*models.Card) error) mongo.Operation {
	return func(col *mgo.Collection) error {
		return mongo.PipeOne(ctx, col, query.pipeline(), &models.Card{}, func(entity interface{}) error {
			if card, ok := entity.(*models.Card); ok {
				return visitor(card)
			}

			return ErrInvalidType
		})
	}
}

func (query CardQuery) pipeline() []bson.M {
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

// CreateCardCommand type
type CreateCardCommand struct {
	BoardID string
	Card    *models.Card
}

// Operation to create card
func (command CreateCardCommand) Operation(ctx context.Context) mongo.Operation {
	return func(col *mgo.Collection) error {
		return mongo.Update(ctx, col, mongo.FromID(command.BoardID), mongo.AddToSet("cards", command.Card))
	}
}

// RemoveCardCommand type
type RemoveCardCommand struct {
	BoardID string
	ID      string
}

// Operation to remove card
func (command RemoveCardCommand) Operation(ctx context.Context) mongo.Operation {
	return func(col *mgo.Collection) error {
		return mongo.Update(ctx, col, mongo.FromID(command.BoardID), mongo.RemoveFromSet("cards", mongo.FromID(command.ID)))
	}
}

// UpdateCardCommand type
type UpdateCardCommand struct {
	BoardID string
	ID      string
	Field   string
	Value   interface{}
}

// Operation to update card
func (command UpdateCardCommand) Operation(ctx context.Context) mongo.Operation {
	return func(col *mgo.Collection) error {
		return mongo.Update(
			ctx,
			col,
			bson.M{"_id": bson.ObjectIdHex(command.BoardID), "cards._id": bson.ObjectIdHex(command.ID)},
			mongo.Set("cards.$."+command.Field, command.Value))
	}
}
