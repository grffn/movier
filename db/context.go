package db

import (
	"log"
	"os"
	"time"

	"github.com/grffn/movier/Godeps/_workspace/src/github.com/jinzhu/gorm"
	_ "github.com/grffn/movier/Godeps/_workspace/src/github.com/lib/pq"
)

type Context struct {
	DB gorm.DB
}

func (c *Context) Init() {
	var err error
	log.Println(os.Getenv("DATABASE_URL"))
	c.DB, err = gorm.Open("postgres", os.Getenv("DATABASE_URL"))

	if err != nil {
		log.Fatalf("Database connection error: %s", err)
	}
}

func (c *Context) InitSchema() {
	c.DB.AutoMigrate(&User{})
}

func CreateContext() *Context {
	var c = Context{}
	c.Init()
	return &c
}

//Model -- GORM base model
type Model struct {
	ID        uint `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

//User db model
type User struct {
	gorm.Model
	Username string
	Password string
	Salt     string
}

//RegisterModel - Registration View Model
type RegisterModel struct {
	Email    string `json:"email" binding:"required"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

//LoginModel - Login View Model
type LoginModel struct {
	UserID   string `json:"userid" binding:"required"`
	Password string `json:"password" binding:"required"`
}
