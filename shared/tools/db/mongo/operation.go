package mongo

import (
	"context"
	"github.com/mongodb/mongo-go-driver/mongo"
)

// OperationContext declares operation context
type OperationContext struct {
	ctx context.Context
	db  string
	col string
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
type OperationHandler func(context.Context, *mongo.Collection) error

// OperationExecutor executes operation
type OperationExecutor interface {
	 Execute(*OperationContext, OperationHandler) error
}