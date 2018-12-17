package mongo

import (
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/dmibod/kanban/shared/tools/log"
)

type Options struct {
	client *mongo.Client
	db     string
	logger log.Logger
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

// WithLogger initializes Logger option
func WithLogger(l log.Logger) Option {
	return func(o *Options) {
		o.logger = l
	}
}

