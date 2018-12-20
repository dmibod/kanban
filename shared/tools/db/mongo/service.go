package mongo

import (
	"sync"
	"time"

	"github.com/mongodb/mongo-go-driver/mongo/options"

	"context"

	"github.com/dmibod/kanban/shared/tools/logger"
	"github.com/dmibod/kanban/shared/tools/logger/noop"
	"github.com/mongodb/mongo-go-driver/mongo"
)

const defaultAddr = "mongodb://localhost:27017"

// Service declares database service
type Service struct {
	sync.Mutex
	cmu    sync.Mutex
	dmu    sync.Mutex
	client *mongo.Client
	dbs    map[string]*mongo.Database
	cols   map[string]*mongo.Collection
	logger logger.Logger
}

// CreateService creates database service instance
func CreateService(l logger.Logger) *Service {
	if l == nil {
		l = &noop.Logger{}
	}

	return &Service{
		logger: l,
		dbs:    make(map[string]*mongo.Database),
		cols:   make(map[string]*mongo.Collection),
	}
}

// Execute executes operation
func (s *Service) Execute(c *OperationContext, h OperationHandler) error {
	err := s.ensureClient()
	if err != nil {
		s.logger.Errorln("cannot obtain client")
		return err
	}

	err = h(s.getCollection(c))
	if err != nil {
		s.logger.Errorf("%v (%T)\n", err, err)
		s.reset()
	}

	return err
}

func newClient() (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	opts := options.Client()

	opts.SetConnectTimeout(time.Second * 2)
	opts.SetServerSelectionTimeout(time.Second * 2)

	creds := options.Credential{
		AuthSource: "admin",
		Username:   "mongoadmin",
		Password:   "secret",
	}

	opts.SetAuth(creds)

	return mongo.Connect(ctx, defaultAddr, opts)
}

func (s *Service) ensureClient() error {
	s.Lock()
	defer s.Unlock()
	if s.client != nil {
		return nil
	}

	client, err := newClient()
	if err != nil {
		return err
	}

	s.client = client

	return nil
}

func (s *Service) reset() {
	s.Lock()
	defer s.Unlock()

	if s.client != nil {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()

		err := s.client.Ping(ctx, nil)
		if err == nil {
			s.logger.Debugln("ping ok")
			return
		}

		s.logger.Debugln("disconnect client")

		err = s.client.Disconnect(ctx)
		if err != nil {
			s.logger.Errorln("error disconnect client", err)
		}
	}

	s.logger.Debugln("reset client")

	s.dmu.Lock()
	defer s.dmu.Unlock()
	s.dbs = make(map[string]*mongo.Database)

	s.cmu.Lock()
	defer s.cmu.Unlock()
	s.cols = make(map[string]*mongo.Collection)

	s.client = nil
}

func (s *Service) getDatabase(name string) *mongo.Database {
	s.dmu.Lock()
	defer s.dmu.Unlock()
	db, ok := s.dbs[name]
	if ok {
		return db
	}
	db = s.client.Database(name)
	s.dbs[name] = db
	return db
}

func (s *Service) getCollection(c *OperationContext) *mongo.Collection {
	s.cmu.Lock()
	defer s.cmu.Unlock()
	col, ok := s.cols[c.col]
	if ok {
		return col
	}
	col = s.getDatabase(c.db).Collection(c.col)
	s.cols[c.col] = col
	return col
}
