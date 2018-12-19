package mongo

import (
	"sync"
	"time"

	"github.com/mongodb/mongo-go-driver/mongo/options"

	"context"
	"errors"

	"github.com/dmibod/kanban/shared/tools/db"
	"github.com/dmibod/kanban/shared/tools/log"
	"github.com/dmibod/kanban/shared/tools/log/logger"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/primitive"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/x/bsonx"
)

const defaultAddr = "mongodb://localhost:27017"

// DatabaseService declares database service
type DatabaseService struct {
	sync.Mutex
	cmu sync.Mutex
	dmu sync.Mutex
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

// DatabaseCommand declares database command
type DatabaseCommand struct {
	db  string
	col string
}

// CreateDatabaseCommand creates DatabaseCommand
func CreateDatabaseCommand(db string, col string) *DatabaseCommand {
	return &DatabaseCommand{
		db:  db,
		col: col,
	}
}

// DatabaseCommandHandler declares DatabaseCommand handler
type DatabaseCommandHandler func(*mongo.Collection) error

// Exec executes DatabaseCommand
func (s *DatabaseService) Exec(c *DatabaseCommand, h DatabaseCommandHandler) error {
	err := s.ensureClient()
	if err != nil {
		s.logger.Errorln("cannot obtain client")
		return err
	}

	err = h(s.getCollection(c.db, c.col))
	if err != nil {
		s.logger.Errorf("execution error: %v\n", err)
		s.reset()
	}

	return err
}

func newClient() (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	opts := options.Client()

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
	s.logger.Debugln("reset client")

	s.Lock()
	defer s.Unlock()

	if s.client != nil {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()

		err := s.client.Disconnect(ctx)
		if err != nil {
			s.logger.Errorln("error disconnect client", err)
		}
	}

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

var _ db.Factory = (*Factory)(nil)
var _ db.Repository = (*Repository)(nil)

// Factory declares repository factory
type Factory struct {
	s      *DatabaseService
	db     string
	logger log.Logger
}

// CreateFactory creates new repository factory
func (s *DatabaseService) CreateFactory(opts ...Option) *Factory {

	var options Options

	for _, o := range opts {
		o(&options)
	}

	return &Factory{
		logger: s.logger,
		db:     options.db,
		s:      s,
	}
}

// Repository declares repository
type Repository struct {
	s        *DatabaseService
	instance db.InstanceFactory
	cmd      *DatabaseCommand
	logger   log.Logger
}

// Create creates new repository
func (f *Factory) Create(col string, instance db.InstanceFactory) db.Repository {
	return &Repository{
		s:        f.s,
		instance: instance,
		cmd:      CreateDatabaseCommand(f.db, col),
		logger:   f.logger,
	}
}

// Create creates new document
func (r *Repository) Create(e interface{}) (string, error) {
	var res string
	var err error
	r.s.Exec(r.cmd, func(col *mongo.Collection) error {
		res, err = r.create(col, e)
		return err
	})
	return res, err
}

// FindByID finds document by its id
func (r *Repository) FindByID(id string) (interface{}, error) {
	var res interface{}
	var err error
	r.s.Exec(r.cmd, func(col *mongo.Collection) error {
		res, err = r.findByID(col, id)
		return err
	})
	return res, err
}

// Find dins all documents by criteria
func (r *Repository) Find(c interface{}, v db.Visitor) error {
	return r.s.Exec(r.cmd, func(col *mongo.Collection) error {
		return r.find(col, c, v)
	})
}

// Count returns count of documents by criteria
func (r *Repository) Count(c interface{}) (int, error) {
	return 0, nil
}

// Update updates document
func (r *Repository) Update(e interface{}) error {
	return nil
}

// Remove removes document
func (r *Repository) Remove(id string) error {
	return nil
}

func (r *Repository) create(col *mongo.Collection, e interface{}) (string, error) {
	res, err := col.InsertOne(context.Background(), e)

	if err != nil {
		r.logger.Errorln("cannot insert document")
		return "", err
	}

	id, ok := res.InsertedID.(primitive.ObjectID)

	if !ok {
		r.logger.Errorln("invalid document id")
		return "", errors.New("Cannot decode id")
	}

	return id.Hex(), nil
}

func (r *Repository) findByID(col *mongo.Collection, id string) (interface{}, error) {
	oid, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		r.logger.Errorln("invalid document id")
		return nil, err
	}

	res := col.FindOne(context.Background(), bson.D{{"_id", bsonx.ObjectID(oid)}})

	e := r.instance()

	err = res.Decode(e)

	if err != nil {
		r.logger.Errorln("cannot decode document")
		return nil, err
	}

	return e, nil
}

func (r *Repository) find(col *mongo.Collection, c interface{}, v db.Visitor) error {
	cur, err := col.Find(context.Background(), c)

	if err != nil {
		r.logger.Errorln("error getting cursor")
		return err
	}

	defer cur.Close(context.Background())

	for cur.Next(context.Background()) {

		e := r.instance()

		err = cur.Decode(e)

		if err != nil {
			r.logger.Errorln("cannot decode document")
			return err
		}

		v(e)
	}

	return nil
}
