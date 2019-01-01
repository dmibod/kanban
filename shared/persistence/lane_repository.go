package persistence

import (
	"gopkg.in/mgo.v2/bson"

	"github.com/dmibod/kanban/shared/tools/db"
)

// HLayout horizontal
const HLayout = "H"

// VLayout vertical
const VLayout = "V"

// LType contains lanes
const LType = "L"

// CType containes cards
const CType = "C"

// LaneEntity definition
type LaneEntity struct {
	ID       bson.ObjectId `bson:"_id,omitempty"`
	Layout   string        `bson:"layout"`
	Type     string        `bson:"type"`
	Name     string        `bson:"name"`
	Children []string      `bson:"children"`
}

// CreateLaneRepository creates repository
func CreateLaneRepository(f db.RepositoryFactory) db.Repository {
	instance := func() interface{} {
		return &LaneEntity{}
	}
	identity := func(entity interface{}) string {
		return entity.(*LaneEntity).ID.Hex()
	}
	return f.CreateRepository("lanes", instance, identity)
}
