package mongo

import (
	"sync"
	"time"

	"github.com/mongodb/mongo-go-driver/mongo/options"

	"context"

	"github.com/dmibod/kanban/shared/tools/log"
	"github.com/dmibod/kanban/shared/tools/log/logger"
	"github.com/mongodb/mongo-go-driver/mongo"
)

const defaultAddr = "mongodb://localhost:27017"

// DatabaseService declares database service
type DatabaseService struct {
	sync.Mutex
	cmu    sync.Mutex
	dmu    sync.Mutex
	client *mongo.Client
	dbs    map[string]*mongo.Database
	cols   map[string]*mongo.Collection
	logger log.Logger
}

// CreateDatabaseService creates DatabaseService instance
func CreateDatabaseService(l log.Logger) *DatabaseService {
	if l == nil {
		l = logger.New(logger.WithPrefix("[MONGO] "), logger.WithDebug(true))
	}

	return &DatabaseService{
		logger: l,
		dbs:    make(map[string]*mongo.Database),
		cols:   make(map[string]*mongo.Collection),
	}
}

// Exec executes DatabaseCommand
func (s *DatabaseService) Exec(c *DatabaseCommand, h DatabaseCommandHandler) error {
	err := s.ensureClient()
	if err != nil {
		s.logger.Errorln("cannot obtain client")
		return err
	}

	err = h(s.getCollection(c.db, c.col))
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

func (s *DatabaseService) ensureClient() error {
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

func (s *DatabaseService) reset() {
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

func (s *DatabaseService) getDatabase(n string) *mongo.Database {
	s.dmu.Lock()
	defer s.dmu.Unlock()
	db, ok := s.dbs[n]
	if ok {
		return db
	}
	db = s.client.Database(n)
	s.dbs[n] = db
	return db
}

func (s *DatabaseService) getCollection(db string, n string) *mongo.Collection {
	s.cmu.Lock()
	defer s.cmu.Unlock()
	col, ok := s.cols[n]
	if ok {
		return col
	}
	col = s.getDatabase(db).Collection(n)
	s.cols[n] = col
	return col
}
