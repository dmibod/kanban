package mongo

import (
	"context"

	"gopkg.in/mgo.v2/bson"

	"github.com/dmibod/kanban/shared/tools/logger"
	"gopkg.in/mgo.v2"
)

// Repository type
type Repository struct {
	OperationExecutor
	logger.Logger
	db  string
	col string
}

// FindByID finds document by id
func (r *Repository) FindByID(ctx context.Context, id string, entity interface{}) (interface{}, error) {
	err := r.execute(ctx, func(col *mgo.Collection) error {
		var opErr error
		entity, opErr = r.findByID(ctx, col, id, entity)
		return opErr
	})

	return entity, err
}

func (r *Repository) findByID(ctx context.Context, col *mgo.Collection, id string, entity interface{}) (interface{}, error) {
	err := col.FindId(bson.ObjectIdHex(id)).One(entity)
	if err != nil {
		return nil, err
	}

	return entity, nil
}

// Find documents by criteria
func (r *Repository) Find(ctx context.Context, criteria interface{}, entity interface{}, v func(interface{}) error) error {
	return r.execute(ctx, func(col *mgo.Collection) error {
		return r.find(ctx, col, criteria, entity, v)
	})
}

func (r *Repository) find(ctx context.Context, col *mgo.Collection, criteria interface{}, entity interface{}, v func(interface{}) error) error {
	iter := col.Find(criteria).Iter()
	for iter.Next(entity) {
		if err := v(entity); err != nil {
			iter.Close()
			return err
		}
	}

	return iter.Close()
}

// Count documents by criteria
func (r *Repository) Count(ctx context.Context, criteria interface{}) (int, error) {
	var count int

	err := r.execute(ctx, func(col *mgo.Collection) error {
		var opErr error
		count, opErr = r.count(ctx, col, criteria)
		return opErr
	})

	return count, err
}

func (r *Repository) count(ctx context.Context, col *mgo.Collection, criteria interface{}) (int, error) {
	return col.Find(criteria).Count()
}

// Create new document
func (r *Repository) Create(ctx context.Context, entity interface{}) (string, error) {
	var id string
	err := r.execute(ctx, func(col *mgo.Collection) error {
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
	return r.execute(ctx, func(col *mgo.Collection) error {
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
	return r.execute(ctx, func(col *mgo.Collection) error {
		return r.remove(ctx, col, id)
	})
}

func (r *Repository) remove(ctx context.Context, col *mgo.Collection, id string) error {
	return col.RemoveId(bson.ObjectIdHex(id))
}

func (r *Repository) ExecuteCommands(ctx context.Context, commands []Command) error {
	r.Debugf("executing commands: %v\n", commands)
	return r.execute(ctx, func(col *mgo.Collection) error {
		return r.executeCommands(ctx, col, commands)
	})
}

func (r *Repository) executeCommands(ctx context.Context, col *mgo.Collection, commands []Command) error {
	for _, command := range commands {
		if err := r.executeCommand(ctx, col, command); err != nil {
			r.Errorln(err)
			return err
		}
	}

	return nil
}

func (r *Repository) executeCommand(ctx context.Context, col *mgo.Collection, command Command) error {
	r.Debugf("execute command: %+v\n", &command)

	switch command.Type {
	case InsertCommand:
		return col.Insert(command.Selector, command.Payload)
	case UpdateCommand:
		return col.Update(command.Selector, command.Payload)
	case DeleteCommand:
		return col.Remove(command.Selector)
	default:
		return nil
	}
}

func (r *Repository) execute(ctx context.Context, o Operation) error {
	c := CreateOperationContext(ctx, r.db, r.col)
	return r.Execute(c, func(col *mgo.Collection) error {
		return o(col)
	})
}
