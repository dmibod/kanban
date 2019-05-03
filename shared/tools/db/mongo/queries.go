package mongo

import (
	"context"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// QueryList function
func QueryList(ctx context.Context, col *mgo.Collection, criteria interface{}, entity interface{}, visit func(interface{}) error) error {
	iter := col.Find(criteria).Iter()
	for iter.Next(entity) {
		if err := visit(entity); err != nil {
			iter.Close()
			return err
		}
	}
	return iter.Close()
}

// QueryOne function
func QueryOne(ctx context.Context, col *mgo.Collection, criteria interface{}, entity interface{}, visit func(interface{}) error) error {
	query := col.Find(criteria)
	if err := query.One(entity); err != nil {
		return err
	}
	return visit(entity)
}

// QueryCount function
func QueryCount(ctx context.Context, col *mgo.Collection, criteria interface{}, visit func(int) error) error {
	query := col.Find(criteria)
	if count, err := query.Count(); err == nil {
		return visit(count)
	} else {
		return err
	}
}

// PipeList function
func PipeList(ctx context.Context, col *mgo.Collection, pipeline []bson.M, entity interface{}, visit func(interface{}) error) error {
	pipe := col.Pipe(pipeline)
	iter := pipe.Iter()
	for iter.Next(entity) {
		if err := visit(entity); err != nil {
			iter.Close()
			return err
		}
	}
	return iter.Close()
}

// PipeOne function
func PipeOne(ctx context.Context, col *mgo.Collection, pipeline []bson.M, entity interface{}, visit func(interface{}) error) error {
	pipe := col.Pipe(pipeline)
	if err := pipe.One(entity); err != nil {
		return err
	}
	return visit(entity)
}
