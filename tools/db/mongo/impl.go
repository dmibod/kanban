package mongo

import (
	"context"
	"errors"
	"log"

	"github.com/dmibod/kanban/tools/db"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/primitive"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/mongo/options"
	"github.com/mongodb/mongo-go-driver/x/bsonx"
)

const DefaultAddr = "mongodb://localhost:27017"

var _ db.RepoFactory = (*RepoFactory)(nil)
var _ db.Repository = (*Repository)(nil)

type RepoFactory struct {
	db *mongo.Database
}

type Repository struct {
	instance db.InstanceFactory
	col     *mongo.Collection
}

func newClient() (*mongo.Client, error) {

	opts := options.Client()

	creds := options.Credential{
		AuthSource: "admin",
		Username:   "mongoadmin",
		Password:   "secret",
	}

	opts.SetAuth(creds)

	return mongo.Connect(context.Background(), "mongodb://localhost:27017", opts)
}

func New(opts ...Option) *RepoFactory {

	var options Options

	for _, o := range opts {
		o(&options)
	}

	client := options.client

	if client == nil {

		newClient, err := newClient()

		if err != nil {
			log.Panicln(err)
		}

		client = newClient
	}

	return &RepoFactory{
		db: client.Database(options.db),
	}
}

func (f *RepoFactory) Create(name string, instance db.InstanceFactory) db.Repository {
	return &Repository{
		instance: instance,
		col:     f.db.Collection(name),
	}
}

func (r *Repository) Create(e interface{}) (string, error) {
	res, err := r.col.InsertOne(context.Background(), e)

	if err != nil {
		log.Println("Cannot insert document")
		return "", err
	}

	id, ok := res.InsertedID.(primitive.ObjectID)

	if !ok {
		log.Println("Invalid document id")
		return "", errors.New("Cannot decode id")
	}

	return id.Hex(), nil
}

func (r *Repository) FindById(id string) (interface{}, error) {
	oid, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		log.Println("Invalid document id")
		return nil, err
	}

	res := r.col.FindOne(context.Background(), bson.D{{"_id", bsonx.ObjectID(oid)}})

	e := r.instance()

	err = res.Decode(e)

	if err != nil {
		log.Println("Cannot decode document")
		return nil, err
	}

	return e, nil
}

func (r *Repository) Find(c interface{}, v db.Visitor) error {
	cur, err := r.col.Find(context.Background(), c)

	if err != nil {
		log.Println("Error getting cursor")
		return err
	}

	defer cur.Close(context.Background())

	for cur.Next(context.Background()) {

		e := r.instance()

		err = cur.Decode(e)

		if err != nil {
			log.Println("Cannot decode document")
			return err
		}

		v(e)
	}

	return nil
}

func (r *Repository) Count(c interface{}) (int, error) {
	return 0, nil
}

func (r *Repository) Update(e interface{}) error {
	return nil
}

func (r *Repository) Remove(id string) error {
	return nil
}
