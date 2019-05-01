package mongo

import (
	"context"

	"gopkg.in/mgo.v2/bson"

	"github.com/dmibod/kanban/shared/tools/logger"
	"gopkg.in/mgo.v2"
)

// Repository type
type Repository struct {
	executor OperationExecutor
	logger.Logger
	db  string
	col string
}

// Create new document
func (r *Repository) Create(ctx context.Context, entity interface{}) (string, error) {
	var id string
	err := r.Execute(ctx, func(col *mgo.Collection) error {
		var opErr error
		id, opErr = r.create(ctx, col, entity)
		return opErr
	})

	return id, err
}

func (r *Repository) create(ctx context.Context, col *mgo.Collection, entity interface{}) (string, error) {
	id := bson.NewObjectId()

	_, err := col.UpsertId(id, entity)
	if err != nil {
		r.Errorln("cannot insert document")
		return "", err
	}

	return id.Hex(), nil
}

// Update document
func (r *Repository) Update(ctx context.Context, id string, op bson.M) error {
	return r.Execute(ctx, func(col *mgo.Collection) error {
		return r.update(ctx, col, id, op)
	})
}

func (r *Repository) update(ctx context.Context, col *mgo.Collection, id string, op bson.M) error {
	err := col.Update(bson.M{"_id": bson.ObjectIdHex(id)}, op)
	if err != nil {
		r.Errorln("cannot update document")
		return err
	}

	return nil
}

// Remove document by id
func (r *Repository) Remove(ctx context.Context, id string) error {
	return r.Execute(ctx, func(col *mgo.Collection) error {
		return r.remove(ctx, col, id)
	})
}

func (r *Repository) remove(ctx context.Context, col *mgo.Collection, id string) error {
	return col.RemoveId(bson.ObjectIdHex(id))
}

func (r *Repository) Execute(ctx context.Context, o Operation) error {
	c := CreateOperationContext(ctx, r.db, r.col)
	return r.executor.Execute(c, func(col *mgo.Collection) error {
		return o(col)
	})
}
