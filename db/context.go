package db

import (
	"log"
	"os"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	//Collection of users
	UsersCollection = "users"
	//Collection of documents
	DocumentsCollection = "documents"
)

//Context of application
type Context struct {
	Session *mgo.Session
}

//Init Initialize database connection. Use this method before using database
func (c *Context) Init() {
	if c.Session != nil {
		c.Session = c.Session.Copy()
	} else {
		var err error
		if os.Getenv("debug") == "true" {
			var aLogger *log.Logger
			aLogger = log.New(os.Stderr, "", log.LstdFlags)
			mgo.SetLogger(aLogger)
		}
		c.Session, err = mgo.DialWithTimeout(os.Getenv("MONGO_DB_URL"), time.Duration(5*time.Second))
		if err != nil {
			log.Fatalf("Database connection error: %s", err)
		}
	}
}

//DB Get database
func (c *Context) DB() *mgo.Database {
	return c.Session.DB("movier")
}

func (c *Context) FindUser(userID string) User {
	user := User{}
	query := bson.M{
		"$or": []interface{}{
			bson.M{"username": userID},
			bson.M{"email": userID},
		},
	}
	c.DB().C(UsersCollection).Find(query).One(&user)
	return user
}

func (c *Context) NewUser(user User) error {
	return c.DB().C(UsersCollection).Insert(user)
}

func (c *Context) NewDocument(document Document) error {
	return c.DB().C(DocumentsCollection).Insert(document)
}

func (c *Context) Documents() (documents []Document, err error) {
	err = c.DB().C(DocumentsCollection).Find(nil).All(&documents)
	return
}

func (c *Context) DocumentsByCategory(category string) (documents []Document, err error) {
	err = c.DB().C(DocumentsCollection).Find(bson.M{"category": category}).All(&documents)
	return
}

func (c *Context) DocumentsByTag(tag string) (documents []Document, err error) {
	err = c.DB().C(DocumentsCollection).Find(bson.M{"tags": tag}).All(&documents)
	return
}

func (c *Context) Categories() ([]string, error) {
	var categories []string
	err := c.DB().C(DocumentsCollection).Find(nil).Select(bson.M{"category": 1}).Distinct("category", &categories)
	return categories, err
}

//Close Close database connection
func (c *Context) Close() {
	c.Session.Close()
}

//CreateContext Create context
func CreateContext() *Context {
	var c = Context{}
	c.Init()
	return &c
}
