package mongo

import (
	"context"
	"errors"

	"github.com/dmibod/kanban/shared/tools/logger"
	"github.com/dmibod/kanban/shared/tools/logger/noop"
)

var sessionKey = &struct{}{}

// FromContext gets mongo session from context
func FromContext(ctx context.Context) Session {
	if s, ok := ctx.Value(sessionKey).(Session); ok {
		return s
	}

	return nil
}

// ContextFactory interface
type ContextFactory interface {
	Context(context.Context) (context.Context, error)
}

type contextFactory struct {
	SessionProvider
	logger.Logger
}

// CreateContextFactory instance
func CreateContextFactory(p SessionProvider, l logger.Logger) ContextFactory {
	if l == nil {
		l = &noop.Logger{}
	}
	return &contextFactory{
		SessionProvider: p,
		Logger:          l,
	}
}

func (f *contextFactory) Context(ctx context.Context) (context.Context, error) {
	session := f.getSession(ctx)
	if session == nil {
		err := errors.New("session not found")
		f.Errorln(err)
		return nil, err
	}
	go func() {
		<-ctx.Done()
		f.Debugln("close session")
		session.Close(false)
	}()
	f.Debugln("produce session aware context")
	return context.WithValue(ctx, sessionKey, session), nil
}

func (f *contextFactory) getSession(ctx context.Context) Session {
	return GetSession(CreateContextSessionProvider(ctx, f.Logger), f.SessionProvider)
}
