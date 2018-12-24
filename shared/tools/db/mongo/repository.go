package mongo

import (
	"context"

	"gopkg.in/mgo.v2/bson"

	"github.com/dmibod/kanban/shared/tools/db"
	"github.com/dmibod/kanban/shared/tools/logger"
	"gopkg.in/mgo.v2"
)

var _ db.Repository = (*repository)(nil)

type repository struct {
	executor         OperationExecutor
	instanceFactory  db.InstanceFactory
	instanceIdentity db.InstanceIdentity
	ctx              *OperationContext
	logger           logger.Logger
}

// Create new document
func (r *repository) Create(entity interface{}) (string, error) {
	var id string

	err := r.executor.Execute(r.ctx, func(ctx context.Context, col *mgo.Collection) error {
		var opErr error
		id, opErr = r.create(ctx, col, entity)
		return opErr
	})

	return id, err
}

// FindByID finds document by id
func (r *repository) FindByID(id string) (interface{}, error) {
	var entity interface{}

	err := r.executor.Execute(r.ctx, func(ctx context.Context, col *mgo.Collection) error {
		var opErr error
		entity, opErr = r.findByID(ctx, col, id)
		return opErr
	})

	return entity, err
}

// Find documents by criteria
func (r *repository) Find(criteria interface{}, v db.EntityVisitor) error {
	return r.executor.Execute(r.ctx, func(ctx context.Context, col *mgo.Collection) error {
		return r.find(ctx, col, criteria, v)
	})
}

// Count documents by criteria
func (r *repository) Count(criteria interface{}) (int, error) {
	var count int

	err := r.executor.Execute(r.ctx, func(ctx context.Context, col *mgo.Collection) error {
		var opErr error
		count, opErr = r.count(ctx, col, criteria)
		return opErr
	})

	return count, err
}

// Update document
func (r *repository) Update(entity interface{}) error {
	return r.executor.Execute(r.ctx, func(ctx context.Context, col *mgo.Collection) error {
		return r.update(ctx, col, entity)
	})
}

// Remove document by id
func (r *repository) Remove(id string) error {
	return r.executor.Execute(r.ctx, func(ctx context.Context, col *mgo.Collection) error {
		return r.remove(ctx, col, id)
	})
}

func (r *repository) create(ctx context.Context, col *mgo.Collection, entity interface{}) (string, error) {
	id := bson.NewObjectId()

	_, err := col.UpsertId(id, entity)
	if err != nil {
		r.logger.Errorln("cannot insert document")
		return "", err
	}

	return id.Hex(), nil
}

func (r *repository) update(ctx context.Context, col *mgo.Collection, entity interface{}) error {
	id := bson.ObjectIdHex(r.instanceIdentity(entity))
	err := col.UpdateId(id, entity)
	if err != nil {
		r.logger.Errorln("cannot update document")
		return err
	}

	return nil
}

func (r *repository) remove(ctx context.Context, col *mgo.Collection, id string) error {
	return col.RemoveId(bson.ObjectIdHex(id))
}

func (r *repository) findByID(ctx context.Context, col *mgo.Collection, id string) (interface{}, error) {
	entity := r.instanceFactory()

	err := col.FindId(bson.ObjectIdHex(id)).One(entity)
	if err != nil {
		return nil, err
	}

	return entity, nil
}

func (r *repository) find(ctx context.Context, col *mgo.Collection, criteria interface{}, v db.EntityVisitor) error {
	entity := r.instanceFactory()

	iter := col.Find(criteria).Iter()
	for iter.Next(entity) {
		if v(entity) {
			break
		}
	}

	return iter.Close()
}

func (r *repository) count(ctx context.Context, col *mgo.Collection, criteria interface{}) (int, error) {
	return col.Find(criteria).Count()
}
