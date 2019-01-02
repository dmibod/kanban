package mongo

import (
	"context"
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

// SessionProvider interface
type SessionProvider interface {
	WithSession(context.Context) context.Context
}
