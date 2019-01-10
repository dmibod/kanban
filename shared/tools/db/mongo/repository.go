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
	OperationExecutor
	logger.Logger
	db.RepositoryEntity
	db  string
	col string
}

// Create new document
func (r *repository) Create(ctx context.Context, entity interface{}) (string, error) {
	var id string

	err := r.execute(ctx, func(col *mgo.Collection) error {
		var opErr error
		id, opErr = r.create(ctx, col, entity)
		return opErr
	})

	return id, err
}

// FindByID finds document by id
func (r *repository) FindByID(ctx context.Context, id string) (interface{}, error) {
	var entity interface{}

	err := r.execute(ctx, func(col *mgo.Collection) error {
		var opErr error
		entity, opErr = r.findByID(ctx, col, id)
		return opErr
	})

	return entity, err
}

// Find documents by criteria
func (r *repository) Find(ctx context.Context, criteria interface{}, v db.EntityVisitor) error {
	return r.execute(ctx, func(col *mgo.Collection) error {
		return r.find(ctx, col, criteria, v)
	})
}

// Count documents by criteria
func (r *repository) Count(ctx context.Context, criteria interface{}) (int, error) {
	var count int

	err := r.execute(ctx, func(col *mgo.Collection) error {
		var opErr error
		count, opErr = r.count(ctx, col, criteria)
		return opErr
	})

	return count, err
}

// Update document
func (r *repository) Update(ctx context.Context, entity interface{}) error {
	return r.execute(ctx, func(col *mgo.Collection) error {
		return r.update(ctx, col, entity)
	})
}

// Remove document by id
func (r *repository) Remove(ctx context.Context, id string) error {
	return r.execute(ctx, func(col *mgo.Collection) error {
		return r.remove(ctx, col, id)
	})
}

func (r *repository) execute(ctx context.Context, o Operation) error {
	c := CreateOperationContext(ctx, r.db, r.col)
	return r.Execute(c, func(col *mgo.Collection) error {
		return o(col)
	})
}

func (r *repository) create(ctx context.Context, col *mgo.Collection, entity interface{}) (string, error) {
	id := bson.NewObjectId()

	_, err := col.UpsertId(id, entity)
	if err != nil {
		r.Errorln("cannot insert document")
		return "", err
	}

	return id.Hex(), nil
}

func (r *repository) update(ctx context.Context, col *mgo.Collection, entity interface{}) error {
	id := bson.ObjectIdHex(r.GetID(entity))
	err := col.UpdateId(id, entity)
	if err != nil {
		r.Errorln("cannot update document")
		return err
	}

	return nil
}

func (r *repository) remove(ctx context.Context, col *mgo.Collection, id string) error {
	return col.RemoveId(bson.ObjectIdHex(id))
}

func (r *repository) findByID(ctx context.Context, col *mgo.Collection, id string) (interface{}, error) {
	entity := r.CreateInstance()

	err := col.FindId(bson.ObjectIdHex(id)).One(entity)
	if err != nil {
		return nil, err
	}

	return entity, nil
}

func (r *repository) find(ctx context.Context, col *mgo.Collection, criteria interface{}, v db.EntityVisitor) error {
	entity := r.CreateInstance()

	iter := col.Find(criteria).Iter()
	for iter.Next(entity) {
		if err := v(entity); err != nil {
			iter.Close()
			return err
		}
	}

	return iter.Close()
}

func (r *repository) count(ctx context.Context, col *mgo.Collection, criteria interface{}) (int, error) {
	return col.Find(criteria).Count()
}
