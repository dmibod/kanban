package persistence

import (
	"gopkg.in/mgo.v2/bson"
)

// BoardProjection type
type BoardProjection struct {
	ID          bson.ObjectId `bson:"_id,omitempty"`
	Owner       string        `bson:"owner"`
	Name        string        `bson:"name"`
	Description string        `bson:"description"`
	Layout      string        `bson:"layout"`
	Shared      bool          `bson:"shared"`
}
