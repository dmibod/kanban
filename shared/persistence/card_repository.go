package persistence

import (
	"gopkg.in/mgo.v2/bson"

	"github.com/dmibod/kanban/shared/tools/db"
)

// CardEntity maps card to/from mongo db
type CardEntity struct {
	ID   bson.ObjectId `bson:"_id,omitempty"`
	Name string        `bson:"name"`
}

// CardRepository interface
type CardRepository interface {
	db.RepositoryEntity
	db.Repository
}

type cardRepository struct {
	db.Repository
}

func (r *cardRepository) CreateInstance() interface{} {
	return &CardEntity{}
}

func (r *cardRepository) GetID(entity interface{}) string {
	return entity.(*CardEntity).ID.Hex()
}

// CreateCardRepository creates new cards repository
func CreateCardRepository(f db.RepositoryFactory) CardRepository {
	r := &cardRepository{}
	r.Repository = f.CreateRepository("cards", r)
	return r
}
