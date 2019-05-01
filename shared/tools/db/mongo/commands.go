package mongo

import (
	"context"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func FromID(id string) bson.M {
	return bson.M{"_id": bson.ObjectIdHex(id)}
}

func Compose(elements ...bson.DocElem) bson.M {
	m := bson.M{}

	for _, e := range elements {
		m[e.Name] = e.Value
	}

	return m
}

func Set(field string, value interface{}) bson.M {
	return bson.M{"$set": bson.M{field: value}}
}

func AddToSet(field string, value interface{}) bson.M {
	return bson.M{"$addToSet": bson.M{field: value}}
}

func RemoveFromSet(field string, value interface{}) bson.M {
	return bson.M{"$pullAll": bson.M{field: []interface{}{value}}}
}

func UpdateInSet(field string, value interface{}) bson.M {
	return bson.M{"$pullAll": bson.M{field: []interface{}{value}}}
}

// Insert function
func Insert(ctx context.Context, col *mgo.Collection, entity interface{}) error {
	return col.Insert(entity)
}

// Update function
func Update(ctx context.Context, col *mgo.Collection, criteria interface{}, entity interface{}) error {
	return col.Update(criteria, entity)
}

// Remove function
func Remove(ctx context.Context, col *mgo.Collection, criteria interface{}) error {
	return col.Remove(criteria)
}
