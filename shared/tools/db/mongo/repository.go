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
	executor        OperationExecutor
	instanceFactory db.InstanceFactory
	ctx             *OperationContext
	logger          logger.Logger
}

// Create creates new document
func (r *Repository) Create(entity interface{}) (string, error) {
	var id string

	err := r.executor.Execute(r.ctx, func(ctx context.Context, col *mgo.Collection) error {
		var opErr error
		id, opErr = r.create(ctx, col, entity)
		return opErr
	})

	return id, err
}

// FindByID finds document by its id
func (r *Repository) FindByID(id string) (interface{}, error) {
	var entity interface{}

	err := r.executor.Execute(r.ctx, func(ctx context.Context, col *mgo.Collection) error {
		var opErr error
		entity, opErr = r.findByID(ctx, col, id)
		return opErr
	})
	
	return entity, err
}

// Find dins all documents by criteria
func (r *Repository) Find(criteria interface{}, v db.EntityVisitor) error {
	return r.executor.Execute(r.ctx, func(ctx context.Context, col *mgo.Collection) error {
		return r.find(ctx, col, criteria, v)
	})
}

// Count returns count of documents by criteria
func (r *Repository) Count(criteria interface{}) (int, error) {
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
	entity := r.instanceFactory()

	err := col.FindId(bson.ObjectIdHex(id)).One(entity)
	if err != nil {
		return nil, err
	}

	return entity, nil
}

func (r *Repository) find(ctx context.Context, col *mgo.Collection, criteria interface{}, v db.EntityVisitor) error {
	entity := r.instanceFactory()

	iter := col.Find(criteria).Iter()
	for iter.Next(entity) {
		if v(entity) {
			break
		}
	}

	return iter.Close()
}
