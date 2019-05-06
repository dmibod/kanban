package models

import (
	"gopkg.in/mgo.v2/bson"
)

// BoardListModel type
type BoardListModel struct {
	ID          bson.ObjectId `bson:"_id,omitempty"`
	Owner       string        `bson:"owner"`
	Name        string        `bson:"name"`
	Description string        `bson:"description"`
	Layout      string        `bson:"layout"`
	Shared      bool          `bson:"shared"`
}

// LaneListModel type
type LaneListModel struct {
	ID          bson.ObjectId `bson:"_id,omitempty"`
	Kind        string        `bson:"kind"`
	Name        string        `bson:"name"`
	Description string        `bson:"description"`
	Layout      string        `bson:"layout"`
}
