package mongo

import (
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

// OperationExecutor executes operation
type OperationExecutor interface {
	Execute(*OperationContext, Operation) error
}

type executor struct {
	sync.Mutex
	url      string
	timeout  time.Duration
	authdb   string
	user     string
	password string
	session  *mgo.Session
	logger   logger.Logger
}

// CreateExecutor creates executor
func CreateExecutor(opts ...Option) OperationExecutor {
	var o options

	for _, opt := range opts {
		opt(&o)
	}

	l := o.logger
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
		logger:   l,
		url:      url,
		timeout:  t,
		authdb:   a,
		user:     u,
		password: p,
	}
}

// Execute operation
func (e *executor) Execute(c *OperationContext, o Operation) error {
	err := e.ensureSession(c)
	if err != nil {
		e.logger.Errorln("cannot open session")
		return err
	}

	err = o(c.ctx, e.session.DB(c.db).C(c.col))
	if err != nil {
		e.dropDeadSession()
	}

	return err
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

func (e *executor) ensureSession(ctx *OperationContext) error {
	if ctx.session != nil {
		return nil
	}

	e.Lock()
	defer e.Unlock()

	if e.session == nil {
		session, err := e.newSession()
		if err != nil {
			return err
		}
		e.logger.Debugln("new session")
		e.session = session
	}

	e.logger.Debugln("open request session")
	ctx.session = e.session.Copy()
	go func() {
		<-ctx.ctx.Done()
		e.logger.Debugln("close request session")
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
			e.logger.Debugln("ping ok")
			return
		}

		e.logger.Debugln("close session")

		e.session.Close()
		e.session = nil
	}
}
