package board

import (
	"gopkg.in/mgo.v2/bson"
)

// BoardEntity entity
type BoardEntity struct {
	ID          bson.ObjectId `bson:"_id,omitempty"`
	Owner       string        `bson:"owner"`
	Name        string        `bson:"name"`
	Description string        `bson:"description"`
	Layout      string        `bson:"layout"`
	Shared      bool          `bson:"shared"`
	Children    []string      `bson:"children"`
}

// BoardProjection type
type BoardProjection struct {
	ID          bson.ObjectId `bson:"_id,omitempty"`
	Owner       string        `bson:"owner"`
	Name        string        `bson:"name"`
	Description string        `bson:"description"`
	Layout      string        `bson:"layout"`
	Shared      bool          `bson:"shared"`
}

// Board entity
type Board struct {
	ID          bson.ObjectId   `bson:"_id,omitempty"`
	Owner       string          `bson:"owner"`
	Name        string          `bson:"name"`
	Description string          `bson:"description"`
	Layout      string          `bson:"layout"`
	Shared      bool            `bson:"shared"`
	Children    []bson.ObjectId `bson:"children"`
	Lanes       []Lane          `bson:"lanes"`
	Cards       []Card          `bson:"cards"`
}

// Lane entity
type Lane struct {
	ID          bson.ObjectId   `bson:"_id,omitempty"`
	Kind        string          `bson:"kind"`
	Name        string          `bson:"name"`
	Description string          `bson:"description"`
	Layout      string          `bson:"layout"`
	Children    []bson.ObjectId `bson:"children"`
}

// Card entity
type Card struct {
	ID          bson.ObjectId `bson:"_id,omitempty"`
	Name        string        `bson:"name"`
	Description string        `bson:"description"`
}
