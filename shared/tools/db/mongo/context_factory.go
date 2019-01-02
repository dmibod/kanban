package mongo

import (
	"context"

	"github.com/dmibod/kanban/shared/tools/logger"
	"github.com/dmibod/kanban/shared/tools/logger/noop"
	"gopkg.in/mgo.v2"
)

type sessionKeyType struct{}

var sessionKey = &sessionKeyType{}

// FromContext gets mongo session from context
func FromContext(ctx context.Context) *mgo.Session {
	if s, ok := ctx.Value(sessionKey).(*mgo.Session); ok {
		return s
	}

	return nil
}

// ContextFactory interface
type ContextFactory interface {
	Context(context.Context) (context.Context, error)
}

type contextFactory struct {
	SessionFactory
	logger.Logger
}

// CreateContextFactory instance
func CreateContextFactory(f SessionFactory, l logger.Logger) ContextFactory {
	if l == nil {
		l = &noop.Logger{}
	}
	return &contextFactory{
		SessionFactory: f,
		Logger:         l,
	}
}

func (f *contextFactory) Context(ctx context.Context) (context.Context, error) {
	session := FromContext(ctx)
	if session != nil {
		f.Debugln("context is session aware")
		return ctx, nil
	}
	session, err := f.Session()
	if err != nil {
		f.Errorln(err)
		return nil, err
	}
	go func() {
		<-ctx.Done()
		f.Debugln("close context session")
		session.Close()
		session = nil
	}()
	f.Debugln("produce session aware context")
	return context.WithValue(ctx, sessionKey, session), nil
}
