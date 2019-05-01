package persistence

import (
	"gopkg.in/mgo.v2/bson"
)

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

// BoardListModel type
type BoardListModel struct {
	ID          bson.ObjectId `bson:"_id,omitempty"`
	Owner       string        `bson:"owner"`
	Name        string        `bson:"name"`
	Description string        `bson:"description"`
	Layout      string        `bson:"layout"`
	Shared      bool          `bson:"shared"`
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

// LaneListModel type
type LaneListModel struct {
	ID          bson.ObjectId   `bson:"_id,omitempty"`
	Kind        string          `bson:"kind"`
	Name        string          `bson:"name"`
	Description string          `bson:"description"`
	Layout      string          `bson:"layout"`
}

// Card entity
type Card struct {
	ID          bson.ObjectId `bson:"_id,omitempty"`
	Name        string        `bson:"name"`
	Description string        `bson:"description"`
}
