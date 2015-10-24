package main

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"log"
	"net/http"
	"os"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	jwt_lib "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/contrib/jwt"
	"github.com/gin-gonic/gin"
	"github.com/grffn/movier/models"
	"golang.org/x/crypto/scrypt"

	"github.com/grffn/movier/db"
)

var (
	secret = os.Getenv("JWT_SECRET")
)

func main() {

	mgo.SetDebug(true)

	port := os.Getenv("PORT")
	router := gin.Default()
	router.POST("/register", makeHandler(registrationHandler))
	router.POST("/login", makeHandler(loginHandler))
	authenticated := router.Group("/")
	authenticated.Use(authHandler)
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

func loginHandler(context *gin.Context, database *db.Context) {
	var model models.LoginModel
	err := context.BindJSON(&model)
	if err != nil {
		log.Println(err)
		context.AbortWithError(http.StatusBadRequest, err)
		return
	}
	user := db.User{}
	query := bson.M{
		"$or": []interface{}{
			bson.M{"username": model.UserID},
			bson.M{"email": model.UserID},
		},
	}
	database.DB().C("users").Find(query).One(&user)
	log.Println(user.Password)
	storedPassword, _ := base64.URLEncoding.DecodeString(user.Password)
	log.Println(storedPassword)
	salt, _ := base64.URLEncoding.DecodeString(user.Salt)
	checkPassword, _ := getPassword([]byte(model.Password), salt)
	if bytes.Compare(storedPassword, checkPassword) == 0 {
		token := jwt_lib.New(jwt_lib.GetSigningMethod("HS256"))
		token.Claims["ID"] = user.Username
		token.Claims["exp"] = time.Now().Add(time.Hour * 1).Unix()
		tokenString, err := token.SignedString([]byte(secret))
		if err != nil {
			context.JSON(500, gin.H{"message": "Could not generate token"})
			return
		}
		context.JSON(200, gin.H{"token": tokenString})
	} else {
		context.JSON(http.StatusUnauthorized, gin.H{"status": "Login or password is incorrect"})
	}
}

func registrationHandler(context *gin.Context, database *db.Context) {
	model := models.RegisterModel{}
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
		ID:       bson.NewObjectId(),
		Username: model.Username,
		Email:    model.Email,
		Password: base64.URLEncoding.EncodeToString(password),
		Salt:     base64.URLEncoding.EncodeToString(salt),
	}

	database.DB().C("users").Insert(user)
	context.JSON(http.StatusOK, "")
}

func authHandler(context *gin.Context) {
	jwt.Auth(secret)(context)
}

func getPassword(password []byte, salt []byte) ([]byte, error) {
	return scrypt.Key(password, salt, 16384, 8, 1, 32)
}
