package persistence

import (
	"gopkg.in/mgo.v2/bson"

	"github.com/dmibod/kanban/shared/tools/db"
)

// LaneEntity definition
type LaneEntity struct {
	ID       bson.ObjectId `bson:"_id,omitempty"`
	Layout   string        `bson:"layout"`
	Type     string        `bson:"type"`
	Name     string        `bson:"name"`
	Children []string      `bson:"children"`
}

// LaneRepository interface
type LaneRepository interface {
	db.RepositoryEntity
	db.Repository
}

type laneRepository struct {
	db.Repository
}

func (r *laneRepository) CreateInstance() interface{} {
	return &LaneEntity{}
}

func (r *laneRepository) GetID(entity interface{}) string {
	return entity.(*LaneEntity).ID.Hex()
}

// CreateLaneRepository creates repository
func CreateLaneRepository(f db.RepositoryFactory) db.Repository {
	r := &laneRepository{}
	r.Repository = f.CreateRepository("lanes", r)
	return r
}
