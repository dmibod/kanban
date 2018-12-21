package mongo

import (
	"context"
	"gopkg.in/mgo.v2"
)

// OperationContext declares operation context
type OperationContext struct {
	ctx     context.Context
	session *mgo.Session
	db      string
	col     string
}

// CreateOperationContext creates OperationContext
func CreateOperationContext(ctx context.Context, db string, col string) *OperationContext {
	if ctx == nil {
		ctx = context.TODO()
	}
	return &OperationContext{
		ctx: ctx,
		db:  db,
		col: col,
	}
}

// OperationHandler declares Operation handler
type OperationHandler func(context.Context, *mgo.Collection) error

// OperationExecutor executes operation
type OperationExecutor interface {
	Execute(*OperationContext, OperationHandler) error
}
