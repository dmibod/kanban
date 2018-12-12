package mongo

import (
	"github.com/mongodb/mongo-go-driver/mongo"
)

type Options struct {
	client *mongo.Client
	db     string
}

type Option func(*Options)

func WithClient(c *mongo.Client) Option {
	return func(o *Options) {
		o.client = c
	}
}

func WithDatabase(db string) Option {
	return func(o *Options) {
		o.db = db
	}
}
