package db

import "gopkg.in/mgo.v2/bson"

//User db model
type User struct {
	ID       bson.ObjectId `bson:"_id"`
	Username string        `bson:"username"`
	Password string        `bson:"password"`
	Email    string        `bson:"email"`
	Salt     string        `bson:"salt"`
}

type Source struct {
	ID     bson.ObjectId `bson:"_id"`
	UserID uint          `bson:"userId"`
	Name   string        `bson:"name"`
}
