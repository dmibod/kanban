package mongo

import (
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

// SessionFactory interface
type SessionFactory interface {
	Session() (*mgo.Session, error)
}

type options struct {
	logger.Logger
	url      string
	timeout  time.Duration
	authdb   string
	user     string
	password string
}

// Option initializes options properties
type Option func(*options)

// WithURL initializes url option
func WithURL(u string) Option {
	return func(o *options) {
		o.url = u
	}
}

// WithTimeout initializes timeout
func WithTimeout(t time.Duration) Option {
	return func(o *options) {
		o.timeout = t
	}
}

// WithAuthDb initializes authdb option
func WithAuthDb(db string) Option {
	return func(o *options) {
		o.authdb = db
	}
}

// WithUser initializes user option
func WithUser(u string) Option {
	return func(o *options) {
		o.user = u
	}
}

// WithPassword initializes password option
func WithPassword(p string) Option {
	return func(o *options) {
		o.password = p
	}
}

// WithLogger initializes logger option
func WithLogger(l logger.Logger) Option {
	return func(o *options) {
		o.Logger = l
	}
}

type sessionFactory struct {
	logger.Logger
	url      string
	timeout  time.Duration
	authdb   string
	user     string
	password string
}

// CreateSessionFactory instance
func CreateSessionFactory(opts ...Option) SessionFactory {
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

	return &sessionFactory{
		url:      url,
		timeout:  t,
		authdb:   a,
		user:     u,
		password: p,
		Logger:   l,
	}
}

func (f *sessionFactory) Session() (*mgo.Session, error) {
	opts := &mgo.DialInfo{
		Addrs:    []string{f.url},
		Timeout:  f.timeout,
		Database: f.authdb,
		Username: f.user,
		Password: f.password,
	}

	session, err := mgo.DialWithInfo(opts)
	
	if err != nil {
		f.Errorln(err)
		return nil, err
	}

	f.Debugln("session created")

	session.SetMode(mgo.Monotonic, true)

	return session, nil
}
