package mongo

import (
	"context"

	"gopkg.in/mgo.v2"
)

// OperationContext declares operation context
type OperationContext struct {
	context.Context
	db  string
	col string
}

// CreateOperationContext creates OperationContext
func CreateOperationContext(ctx context.Context, db string, col string) *OperationContext {
	return &OperationContext{
		db:      db,
		col:     col,
		Context: ctx,
	}
}

// Operation to be performed on mongo collection
type Operation func(*mgo.Collection) error
