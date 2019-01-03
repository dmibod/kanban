package mongo

import (
	"context"
	"sync"

	"github.com/dmibod/kanban/shared/tools/logger"
	"github.com/dmibod/kanban/shared/tools/logger/noop"
	"gopkg.in/mgo.v2"
)

// SessionProvider interface
type SessionProvider interface {
	Provide() Session
}

// CreateSessionProvider instance
func CreateSessionProvider(f SessionFactory, l logger.Logger) SessionProvider {
	if l == nil {
		l = &noop.Logger{}
	}
	return &sessionProvider{
		SessionFactory: f,
		Logger:         l,
	}
}

type sessionProvider struct {
	sync.Mutex
	logger.Logger
	SessionFactory
	session *mgo.Session
}

func (p *sessionProvider) Provide() Session {
	p.Lock()
	defer p.Unlock()
	if p.session == nil {
		session, err := p.SessionFactory.Session()
		if err != nil {
			p.Errorln(err)
			return nil
		}
		p.Debugln("session created")
		p.session = session
	}
	return p
}

func (p *sessionProvider) Session() *mgo.Session {
	return p.session
}

func (p *sessionProvider) Close() {
	p.Lock()
	defer p.Unlock()
	if p.session != nil {
		p.Debugln("session closed")
		p.session.Close()
		p.session = nil
	}
}

func (p *sessionProvider) Release() {
	p.Lock()
	defer p.Unlock()
	if p.session != nil {
		if err := p.session.Ping(); err != nil {
			p.Debugln(err)
			p.Close()
			return
		}
		p.Debugln("ping ok")
	}
}

// CreateCopySessionProvider instance
func CreateCopySessionProvider(p SessionProvider, l logger.Logger) SessionProvider {
	if l == nil {
		l = &noop.Logger{}
	}
	return &copySessionProvider{
		provider: p,
		Logger:   l,
	}
}

type copySessionProvider struct {
	logger.Logger
	provider SessionProvider
}

func (p *copySessionProvider) Provide() Session {
	session := p.provider.Provide()
	if session != nil {
		return &copySession{session: session}
	}
	return nil
}

type copySession struct {
	session Session
	mgo     *mgo.Session
}

func (s *copySession) Session() *mgo.Session {
	if s.mgo == nil {
		session := s.session.Session()
		if session != nil {
			s.mgo = session.Copy()
		}
	}
	return s.mgo
}

func (s *copySession) Release() {
	if s.mgo != nil {
		s.mgo.Close()
		s.mgo = nil
	}
}

func (s *copySession) Close() {
	s.Close()
	s.session.Release()
}

// CreateContextSessionProvider instance
func CreateContextSessionProvider(c context.Context, l logger.Logger) SessionProvider {
	if l == nil {
		l = &noop.Logger{}
	}
	return &contextSessionProvider{
		Context: c,
		Logger:  l,
	}
}

type contextSessionProvider struct {
	context.Context
	logger.Logger
}

func (p *contextSessionProvider) Provide() Session {
	return FromContext(p.Context)
}

func GetSession(providers ...SessionProvider) Session {
	for _, p := range providers {
		session := p.Provide()
		if session != nil {
			return session
		}
	}

	return nil
}
