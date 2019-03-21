package card

import (
	"gopkg.in/mgo.v2/bson"
)

// CardEntity maps card to/from mongo db
type CardEntity struct {
	ID          bson.ObjectId `bson:"_id,omitempty"`
	Name        string        `bson:"name"`
	Description string        `bson:"description"`
}
