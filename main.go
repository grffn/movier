package main

import (
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
	authenticated.POST("/create", makeHandler(handlers.CreateHandler))
	authenticated.GET("/sign", handlers.SignHandler)

	authenticated.GET("/documents", makeHandler(handlers.DocumentsHandler))
	authenticated.GET("/categories", makeHandler(handlers.CategoriesHandler))

	router.StaticFile("/", "web/index.html")
	router.Static("/assets", "web/assets")
	router.Run(":" + port)
}

func makeHandler(handler func(*gin.Context, *db.Context)) gin.HandlerFunc {
	return func(context *gin.Context) {
		db := db.CreateContext()
		defer db.Close()
		handler(context, db)

	}
}
