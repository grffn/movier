package db

import (
	"log"
	"os"
	"time"

	"gopkg.in/mgo.v2"
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
