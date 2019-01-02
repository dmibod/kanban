package mongo

import (
	"context"
	"sync"
	"time"

	"github.com/dmibod/kanban/shared/tools/logger"
	"github.com/dmibod/kanban/shared/tools/logger/noop"
	"gopkg.in/mgo.v2"
)

const (
	defaultURL      = "localhost:27017"
	defaultTimeout  = time.Second
	defaultAuthDB   = "admin"
	defaultUser     = "mongoadmin"
	defaultPassword = "secret"
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

// OperationExecutor executes operation
type OperationExecutor interface {
	Execute(*OperationContext, Operation) error
}

type executor struct {
	sync.Mutex
	logger.Logger
	url      string
	timeout  time.Duration
	authdb   string
	user     string
	password string
	session  *mgo.Session
}

// CreateExecutor creates executor
func CreateExecutor(opts ...Option) OperationExecutor {
	var o options

	for _, opt := range opts {
		opt(&o)
	}

	l := o.Logger
	if l == nil {
		l = &noop.Logger{}
	}

	url := o.url
	if url == "" {
		url = defaultURL
	}

	t := o.timeout
	if t == 0 {
		t = defaultTimeout
	}

	a := o.authdb
	if a == "" {
		a = defaultAuthDB
	}

	u := o.user
	if u == "" {
		u = defaultUser
	}

	p := o.password
	if p == "" {
		p = defaultPassword
	}

	return &executor{
		Logger:   l,
		url:      url,
		timeout:  t,
		authdb:   a,
		user:     u,
		password: p,
	}
}

// WithSession creates context with mongo session
func (e *executor) WithSession(ctx context.Context) context.Context {
	s := FromContext(ctx)
	if s != nil {
		return ctx
	}

	if e.session == nil {
		session, err := e.newSession()
		if err != nil {
			e.Errorln(err)
			return ctx
		}
		e.Debugln("new session")
		e.session = session
	}

	e.Debugln("open request session")
	s = e.session.Copy()
	go func() {
		<-ctx.Done()
		e.Debugln("close request session")
		s.Close()
		s = nil
	}()

	return context.WithValue(ctx, sessionKey, s)
}

// Execute operation
func (e *executor) Execute(c *OperationContext, o Operation) error {
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

func (e *executor) ensureSession(ctx *OperationContext) error {
	if ctx.session != nil {
		return nil
	}
	if s, ok := ctx.Context.Value(sessionKey).(*mgo.Session); ok {
		ctx.session = s
		return nil
	}

	e.Lock()
	defer e.Unlock()

	if e.session == nil {
		session, err := e.newSession()
		if err != nil {
			return err
		}
		e.Debugln("new session")
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

func (e *executor) dropDeadSession() {
	e.Lock()
	defer e.Unlock()

	if e.session != nil {
		err := e.session.Ping()
		if err == nil {
			e.Debugln("ping ok")
			return
		}

		e.Debugln("close session")

		e.session.Close()
		e.session = nil
	}
}

func (e *executor) newSession() (*mgo.Session, error) {
	opts := &mgo.DialInfo{
		Addrs:    []string{e.url},
		Timeout:  e.timeout,
		Database: e.authdb,
		Username: e.user,
		Password: e.password,
	}

	s, err := mgo.DialWithInfo(opts)
	if err == nil {
		s.SetMode(mgo.Monotonic, true)
	}

	return s, err
}
