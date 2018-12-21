package mongo

import (
	"sync"
	"time"

	"github.com/dmibod/kanban/shared/tools/logger"
	"github.com/dmibod/kanban/shared/tools/logger/noop"
	"gopkg.in/mgo.v2"
)

const defaultAddr = "localhost:27017"

// Service declares database service
type Service struct {
	sync.Mutex
	session *mgo.Session
	logger  logger.Logger
}

// CreateService creates database service instance
func CreateService(l logger.Logger) *Service {
	if l == nil {
		l = &noop.Logger{}
	}

	return &Service{
		logger: l,
	}
}

// Execute executes operation
func (s *Service) Execute(c *OperationContext, h OperationHandler) error {
	err := s.ensureSession(c)
	if err != nil {
		s.logger.Errorln("cannot open session")
		return err
	}

	err = h(c.ctx, s.session.DB(c.db).C(c.col))
	if err != nil {
		s.logger.Errorf("%v (%T)\n", err, err)
		s.invalidate()
	}

	return err
}

func newSession() (*mgo.Session, error) {
	opts := &mgo.DialInfo{
		Addrs:    []string{defaultAddr},
		Timeout:  60 * time.Second,
		Database: "admin",
		Username: "mongoadmin",
		Password: "secret",
	}

	s, err := mgo.DialWithInfo(opts)
	if err != nil {
		s.SetMode(mgo.Monotonic, true)
	}

	return s, err
}

func (s *Service) ensureSession(ctx *OperationContext) error {
	if ctx.session != nil {
		return nil
	}

	s.Lock()
	defer s.Unlock()
	if s.session == nil {
		session, err := newSession()
		if err != nil {
			return err
		}
		s.logger.Debugln("new session")
		s.session = session
	}

	s.logger.Debugln("copy session")
	ctx.session = s.session.Copy()
	go func(){
		<-ctx.ctx.Done()
		s.logger.Debugln("close session")
		ctx.session.Close()
		ctx.session = nil
	}()

	return nil
}

func (s *Service) invalidate() {
	s.Lock()
	defer s.Unlock()

	if s.session != nil {
		err := s.session.Ping()
		if err == nil {
			s.logger.Debugln("ping ok")
			return
		}

		s.logger.Debugln("close session")

		s.session.Close()
		s.session = nil
	}
}
