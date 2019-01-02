package mongo

import (
	"sync"

	"github.com/dmibod/kanban/shared/tools/logger"
	"github.com/dmibod/kanban/shared/tools/logger/noop"
	"gopkg.in/mgo.v2"
)

// OperationExecutor executes operation
type OperationExecutor interface {
	Execute(*OperationContext, Operation) error
}

type operationExecutor struct {
	sync.Mutex
	logger.Logger
	SessionFactory
	session *mgo.Session
}

// CreateExecutor instance
func CreateExecutor(f SessionFactory, l logger.Logger) OperationExecutor {
	if l == nil {
		l = &noop.Logger{}
	}

	return &operationExecutor{
		Logger:         l,
		SessionFactory: f,
	}
}

// Execute operation
func (e *operationExecutor) Execute(c *OperationContext, o Operation) error {
	err := e.ensureSession(c)
	if err != nil {
		e.Errorln("cannot open session")
		return err
	}

	err = o(e.session.DB(c.db).C(c.col))
	if err != nil {
		switch err {
		case mgo.ErrNotFound:
		case mgo.ErrCursor:
		default:
			e.dropDeadSession()
		}
	}

	return err
}

func (e *operationExecutor) ensureSession(ctx *OperationContext) error {
	if ctx.session != nil {
		return nil
	}
	if s := FromContext(ctx.Context); s != nil {
		e.Debugln("operation context is session aware")
		ctx.session = s
		return nil
	}

	e.Lock()
	defer e.Unlock()

	if e.session == nil {
		session, err := e.Session()
		if err != nil {
			e.Errorln(err)
			return err
		}
		e.Debugln("session acquired")
		e.session = session
	}

	e.Debugln("open operation session")
	ctx.session = e.session.Copy()
	go func() {
		<-ctx.Context.Done()
		e.Debugln("close operation session")
		ctx.session.Close()
		ctx.session = nil
	}()

	return nil
}

func (e *operationExecutor) dropDeadSession() {
	e.Lock()
	defer e.Unlock()

	if e.session != nil {
		err := e.session.Ping()
		if err == nil {
			e.Debugln("ping ok")
			return
		}

		e.Errorln(err)
		e.Debugln("close session")

		e.session.Close()
		e.session = nil
	}
}
