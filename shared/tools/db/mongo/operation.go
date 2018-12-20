package mongo

import (
	"github.com/mongodb/mongo-go-driver/mongo"
)

// OperationContext declares operation context
type OperationContext struct {
	db  string
	col string
}

// CreateOperationContext creates OperationContext
func CreateOperationContext(db string, col string) *OperationContext {
	return &OperationContext{
		db:  db,
		col: col,
	}
}

// OperationHandler declares Operation handler
type OperationHandler func(*mongo.Collection) error

// OperationExecutor executes operation
type OperationExecutor interface {
	 Execute(*OperationContext, OperationHandler) error
}