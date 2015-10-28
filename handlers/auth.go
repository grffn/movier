package handlers

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"log"
	"net/http"
	"os"
	"time"

	jwt_lib "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/grffn/movier/db"
	"github.com/grffn/movier/models"
	"github.com/grffn/movier/util"
	"gopkg.in/mgo.v2/bson"
)

var (
	secret = os.Getenv("JWT_SECRET")
)

//LoginHandler Hndler for login method
func LoginHandler(context *gin.Context, database *db.Context) {
	var model models.LoginModel
	err := context.BindJSON(&model)
	if err != nil {
		log.Println(err)
		context.AbortWithError(http.StatusBadRequest, err)
		return
	}
	user := database.FindUser(model.UserID)
	storedPassword, _ := base64.URLEncoding.DecodeString(user.Password)
	salt, _ := base64.URLEncoding.DecodeString(user.Salt)
	checkPassword, _ := util.GeneratePassword([]byte(model.Password), salt)
	if bytes.Compare(storedPassword, checkPassword) == 0 {
		token := jwt_lib.New(jwt_lib.GetSigningMethod("HS256"))
		token.Claims["ID"] = user.Username
		token.Claims["exp"] = time.Now().Add(time.Day * 1).Unix()
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

//RegistrationHandler Hndler for register method
func RegistrationHandler(context *gin.Context, database *db.Context) {
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
	password, err = util.GeneratePassword([]byte(model.Password), salt)
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

//AuthHandler Hndler for authentication method
func AuthHandler(context *gin.Context) {
	util.JWTAuth(secret)(context)
}
