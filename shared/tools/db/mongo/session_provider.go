package mongo

import (
	"sync"

	"github.com/dmibod/kanban/shared/tools/logger"
	"github.com/dmibod/kanban/shared/tools/logger/noop"
	"gopkg.in/mgo.v2"
)

// SessionProvider interface
type SessionProvider interface {
	Get() (*mgo.Session, error)
	Release()
}

type sessionProvider struct {
	sync.Mutex
	logger.Logger
	SessionFactory
	session *mgo.Session
}

func CreateSessionProvider(f SessionFactory, l logger.Logger) SessionProvider {
	if l == nil {
		l = &noop.Logger{}
	}
	return &sessionProvider{
		SessionFactory: f,
		Logger:         l,
	}
}

func (p *sessionProvider) Get() (*mgo.Session, error) {
	p.Lock()
	defer p.Unlock()
	if p.session == nil {
		session, err := p.SessionFactory.Session()
		if err != nil {
			p.Errorln(err)
			return nil, err
		}
		p.Debugln("create session")
		p.session = session
	}
	return p.session, nil
}

func (p *sessionProvider) Release() {
	p.Lock()
	defer p.Unlock()
	if p.session != nil {
		p.Debugln("release session")
		p.session.Close()
		p.session = nil
	}
}
