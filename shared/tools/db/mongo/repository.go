package mongo

import (
	"context"
	"errors"

	"github.com/dmibod/kanban/shared/tools/db"
	"github.com/dmibod/kanban/shared/tools/logger"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/primitive"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/x/bsonx"
)

var _ db.Repository = (*Repository)(nil)

// Repository declares repository
type Repository struct {
	executor OperationExecutor
	instance db.InstanceFactory
	ctx      *OperationContext
	logger   logger.Logger
}

// Create creates new document
func (r *Repository) Create(entity interface{}) (string, error) {
	var res string
	err := r.executor.Execute(r.ctx, func(ctx context.Context, col *mongo.Collection) error {
		var e error
		res, e = r.create(ctx, col, entity)
		return e
	})
	return res, err
}

// FindByID finds document by its id
func (r *Repository) FindByID(id string) (interface{}, error) {
	var res interface{}
	err := r.executor.Execute(r.ctx, func(ctx context.Context, col *mongo.Collection) error {
		var e error
		res, e = r.findByID(ctx, col, id)
		return e
	})
	return res, err
}

// Find dins all documents by criteria
func (r *Repository) Find(c interface{}, v db.Visitor) error {
	return r.executor.Execute(r.ctx, func(ctx context.Context, col *mongo.Collection) error {
		return r.find(ctx, col, c, v)
	})
}

// Count returns count of documents by criteria
func (r *Repository) Count(c interface{}) (int, error) {
	return 0, nil
}

// Update updates document
func (r *Repository) Update(entity interface{}) error {
	return nil
}

// Remove removes document
func (r *Repository) Remove(id string) error {
	return nil
}

func (r *Repository) create(ctx context.Context, col *mongo.Collection, entity interface{}) (string, error) {
	res, err := col.InsertOne(ctx, entity)

	if err != nil {
		r.logger.Errorln("cannot insert document")
		return "", err
	}

	id, ok := res.InsertedID.(primitive.ObjectID)

	if !ok {
		r.logger.Errorln("invalid document id")
		return "", errors.New("Cannot decode id")
	}

	return id.Hex(), nil
}

func (r *Repository) findByID(ctx context.Context, col *mongo.Collection, id string) (interface{}, error) {
	oid, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		r.logger.Errorln("invalid document id")
		return nil, err
	}

	res := col.FindOne(ctx, bson.D{{"_id", bsonx.ObjectID(oid)}})

	e := r.instance()

	err = res.Decode(e)

	if err != nil {
		r.logger.Errorln("cannot decode document")
		return nil, err
	}

	return e, nil
}

func (r *Repository) find(ctx context.Context, col *mongo.Collection, c interface{}, v db.Visitor) error {
	cur, err := col.Find(ctx, c)

	if err != nil {
		r.logger.Errorln("error getting cursor")
		return err
	}

	defer cur.Close(ctx)

	for cur.Next(ctx) {

		entity := r.instance()

		err = cur.Decode(entity)

		if err != nil {
			r.logger.Errorln("cannot decode document")
			return err
		}

		v(entity)
	}

	return nil
}
