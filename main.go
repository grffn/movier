package main

import (
	"log"
	"os"

	"gopkg.in/mgo.v2"

	"github.com/gin-gonic/gin"

	"github.com/grffn/movier/db"
	"github.com/grffn/movier/handlers"
)

func main() {

	mgo.SetDebug(true)

	port := os.Getenv("PORT")
	router := gin.Default()
	router.POST("/register", makeHandler(handlers.RegistrationHandler))
	router.POST("/login", makeHandler(handlers.LoginHandler))
	authenticated := router.Group("/")
	authenticated.Use(handlers.AuthHandler)
	authenticated.POST("/create", func(context *gin.Context) {
		log.Println("Hello, World")
	})
	router.Run(":" + port)
}

func makeHandler(handler func(*gin.Context, *db.Context)) gin.HandlerFunc {
	return func(context *gin.Context) {
		db := db.CreateContext()
		defer db.Close()
		handler(context, db)

	}
}
