package lane

import (
	"gopkg.in/mgo.v2/bson"
)

// LaneEntity definition
type LaneEntity struct {
	ID          bson.ObjectId `bson:"_id,omitempty"`
	Kind        string        `bson:"kind"`
	Name        string        `bson:"name"`
	Description string        `bson:"description"`
	Layout      string        `bson:"layout"`
	Children    []string      `bson:"children"`
}
