package mongo

import (
	"context"
	"errors"
	"github.com/dmibod/kanban/shared/tools/logger"
	"github.com/dmibod/kanban/shared/tools/logger/noop"
	"gopkg.in/mgo.v2"
)

// OperationExecutor executes operation
type OperationExecutor interface {
	Execute(*OperationContext, Operation) error
}

type operationExecutor struct {
	logger.Logger
	SessionProvider
}

// CreateExecutor instance
func CreateExecutor(p SessionProvider, l logger.Logger) OperationExecutor {
	if l == nil {
		l = &noop.Logger{}
	}

	return &operationExecutor{
		Logger:          l,
		SessionProvider: p,
	}
}

// Execute operation
func (e *operationExecutor) Execute(c *OperationContext, o Operation) error {
	session := e.getSession(c)
	if session == nil {
		err := errors.New("session not found")
		e.Errorln(err)
		return err
	}

	err := o(session.Session().DB(c.db).C(c.col))
	if err != nil {
		switch err {
		case mgo.ErrNotFound:
		case mgo.ErrCursor:
		default:
			session.Close(true)
		}
	}

	return err
}

func (e *operationExecutor) getSession(ctx context.Context) Session {
	return GetSession(CreateContextSessionProvider(ctx, e.Logger), e.SessionProvider)
}
