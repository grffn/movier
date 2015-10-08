package main

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"os"

	"github.com/grffn/movier/Godeps/_workspace/src/github.com/gin-gonic/contrib/jwt"
	"github.com/grffn/movier/Godeps/_workspace/src/golang.org/x/crypto/scrypt"

	"github.com/grffn/movier/Godeps/_workspace/src/github.com/gin-gonic/gin"
	_ "github.com/grffn/movier/Godeps/_workspace/src/github.com/lib/pq"
	"github.com/grffn/movier/db"
)

var (
	secret = os.Getenv("JWT_SECRET")
)

func main() {
	port := os.Getenv("PORT")
	c := db.CreateContext()
	c.InitSchema()
	router := gin.Default()
	authenticated := router.Group("/")
	authenticated.Use(authHandler)
	authenticated.POST("/create", func(context *gin.Context) {

	})
	router.Run(":" + port)
}

func makeHandler(handler func(*gin.Context, *db.Context)) gin.HandlerFunc {
	return func(context *gin.Context) {
		db := db.CreateContext()
		handler(context, db)
	}
}

func loginHandler(context *gin.Context, database *db.Context) {
	var model db.LoginModel
	err := context.BindJSON(&model)
	if err != nil {
		context.AbortWithStatus(http.StatusBadRequest)
		return
	}
	user := db.User{}
	database.DB.Where("username = ? or email = ?", model.UserID).First(&user)
	storedPassword, _ := base64.URLEncoding.DecodeString(user.Password)
	salt, _ := base64.URLEncoding.DecodeString(user.Salt)
	checkPassword, _ := getPassword([]byte(model.Password), salt)
	if bytes.Compare(storedPassword, checkPassword) == 0 {
		context.JSON(http.StatusOK, "")
	} else {
		context.JSON(http.StatusUnauthorized, gin.H{"status": "Login or password is incorrect"})
		context.Abort()
	}
}

func registrationHandler(context *gin.Context, database *db.Context) {
	model := db.RegisterModel{}
	err := context.BindJSON(&model)
	if err != nil {
		context.AbortWithStatus(http.StatusBadRequest)
		return
	}
	salt := make([]byte, 128)
	_, err = rand.Read(salt)
	if err != nil {
		context.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	var password []byte
	password, err = getPassword([]byte(model.Password), salt)
	if err != nil {
		context.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	user := db.User{
		Username: model.Username,
		Password: base64.URLEncoding.EncodeToString(password),
		Salt:     base64.URLEncoding.EncodeToString(salt),
	}

	database.DB.Create(&user)
	context.JSON(http.StatusOK, "")
}

func authHandler(context *gin.Context) {
	jwt.Auth(secret)
}

func getPassword(password []byte, salt []byte) ([]byte, error) {
	return scrypt.Key(password, salt, 16384, 8, 1, 32)
}
