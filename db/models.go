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

/*Document Document db entry
 */
type Document struct {
	ID       bson.ObjectId `bson:"_id"`
	UserID   bson.ObjectId `bson:"userId"`
	Name     string        `bson:"name"`
	Category string        `bson:"category"`
	Tags     []string      `bson:"tags"`
	URL      string        `bson:"url"`
	MimeType string        `bson:"mimetype"`
	Authors  []string      `bson:"authors"`
}
