package mongo

import (
	"github.com/mongodb/mongo-go-driver/mongo"
)

type FactoryFn func() interface{}

// Options can be used to create a customized transport.
type Options struct {
	FactoryFn FactoryFn
	Client *mongo.Client
	Db string
	Col string
}

// Option is a function on the options for a nats transport.
type Option func(*Options)

func WithClient(c *mongo.Client) Option {
	return func(o *Options) {
		o.Client = c
	}
}

func WithDatabase(db string) Option {
	return func(o *Options) {
		o.Db = db
	}
}

func WithCollection(c string) Option {
	return func(o *Options) {
		o.Col = c
	}
}

func WithFactory(f FactoryFn) Option {
	return func(o *Options) {
		o.FactoryFn = f
	}
}
