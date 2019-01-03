package mongo

import "gopkg.in/mgo.v2"

type Session interface {
	Session() *mgo.Session
	Release()
	Close()
}
