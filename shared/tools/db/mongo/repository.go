package mongo

import (
	"context"

	"gopkg.in/mgo.v2/bson"

	"github.com/dmibod/kanban/shared/tools/db"
	"github.com/dmibod/kanban/shared/tools/logger"
	"gopkg.in/mgo.v2"
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
	err := r.executor.Execute(r.ctx, func(ctx context.Context, col *mgo.Collection) error {
		var e error
		res, e = r.create(ctx, col, entity)
		return e
	})
	return res, err
}

// FindByID finds document by its id
func (r *Repository) FindByID(id string) (interface{}, error) {
	var res interface{}
	err := r.executor.Execute(r.ctx, func(ctx context.Context, col *mgo.Collection) error {
		var e error
		res, e = r.findByID(ctx, col, id)
		return e
	})
	return res, err
}

// Find dins all documents by criteria
func (r *Repository) Find(c interface{}, v db.Visitor) error {
	return r.executor.Execute(r.ctx, func(ctx context.Context, col *mgo.Collection) error {
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

func (r *Repository) create(ctx context.Context, col *mgo.Collection, entity interface{}) (string, error) {
	id := bson.NewObjectId()
	_, err := col.UpsertId(id, entity)
	if err != nil {
		r.logger.Errorln("cannot insert document")
		return "", err
	}

	return id.Hex(), nil
}

func (r *Repository) findByID(ctx context.Context, col *mgo.Collection, id string) (interface{}, error) {
	e := r.instance()
	err := col.FindId(bson.ObjectIdHex(id)).One(e)
	if err != nil {
		return nil, err
	}

	return e, nil
}

func (r *Repository) find(ctx context.Context, col *mgo.Collection, c interface{}, v db.Visitor) error {
	entity := r.instance()

	iter := col.Find(c).Iter()
	for iter.Next(entity) {
		v(entity)
	}

	return iter.Close()
}
