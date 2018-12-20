package mongo

import (
	"context"
	"errors"

	"github.com/dmibod/kanban/shared/tools/db"
	"github.com/dmibod/kanban/shared/tools/log"
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
	logger   log.Logger
}

// Create creates new document
func (r *Repository) Create(e interface{}) (string, error) {
	var res string
	err := r.executor.Execute(r.ctx, func(col *mongo.Collection) error {
		var e error
		res, e = r.create(col, e)
		return e
	})
	return res, err
}

// FindByID finds document by its id
func (r *Repository) FindByID(id string) (interface{}, error) {
	var res interface{}
	err := r.executor.Execute(r.ctx, func(col *mongo.Collection) error {
		var e error
		res, e = r.findByID(col, id)
		return e
	})
	return res, err
}

// Find dins all documents by criteria
func (r *Repository) Find(c interface{}, v db.Visitor) error {
	return r.executor.Execute(r.ctx, func(col *mongo.Collection) error {
		return r.find(col, c, v)
	})
}

// Count returns count of documents by criteria
func (r *Repository) Count(c interface{}) (int, error) {
	return 0, nil
}

// Update updates document
func (r *Repository) Update(e interface{}) error {
	return nil
}

// Remove removes document
func (r *Repository) Remove(id string) error {
	return nil
}

func (r *Repository) create(col *mongo.Collection, e interface{}) (string, error) {
	res, err := col.InsertOne(context.Background(), e)

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

func (r *Repository) findByID(col *mongo.Collection, id string) (interface{}, error) {
	oid, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		r.logger.Errorln("invalid document id")
		return nil, err
	}

	res := col.FindOne(context.Background(), bson.D{{"_id", bsonx.ObjectID(oid)}})

	e := r.instance()

	err = res.Decode(e)

	if err != nil {
		r.logger.Errorln("cannot decode document")
		return nil, err
	}

	return e, nil
}

func (r *Repository) find(col *mongo.Collection, c interface{}, v db.Visitor) error {
	cur, err := col.Find(context.Background(), c)

	if err != nil {
		r.logger.Errorln("error getting cursor")
		return err
	}

	defer cur.Close(context.Background())

	for cur.Next(context.Background()) {

		e := r.instance()

		err = cur.Decode(e)

		if err != nil {
			r.logger.Errorln("cannot decode document")
			return err
		}

		v(e)
	}

	return nil
}
