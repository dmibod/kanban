package mongo

import (
	"gopkg.in/mgo.v2/bson"
)

// CommandType type
type CommandType int

const (
	InsertCommand CommandType = CommandType(iota)
	UpdateCommand
	DeleteCommand
)

// Command type
type Command struct {
	Selector bson.M
	Type     CommandType
	Payload  interface{}
}

func Insert(id string, entity interface{}) Command {
	return Command{
		Selector: FromID(id),
		Type:     InsertCommand,
		Payload:  entity,
	}
}

func Update(id string, field string, value interface{}) Command {
	return Command{
		Selector: FromID(id),
		Type:     UpdateCommand,
		Payload:  bson.M{"$set": bson.M{field: value}},
	}
}

func Remove(id string) Command {
	return Command{
		Selector: FromID(id),
		Type:     DeleteCommand,
	}
}

func CustomUpdate(id string, updater bson.M) Command {
	return Command{
		Selector: FromID(id),
		Type:     UpdateCommand,
		Payload:  updater,
	}
}

func FromID(id string) bson.M {
	return bson.M{"_id": bson.ObjectIdHex(id)}
}

func AddToSet(field string, value interface{}) bson.M {
	return bson.M{"$addToSet": bson.M{field: value}}
}

func PullFromSet(field string, value interface{}) bson.M {
	return bson.M{"$pullAll": bson.M{field: []interface{}{value}}}
}
