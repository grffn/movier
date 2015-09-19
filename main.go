package main

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"log"
	"net/http"
	"time"

	"golang.org/x/crypto/scrypt"

	"github.com/StephanDollberg/go-json-rest-middleware-jwt"
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
)

func main() {
	jwtMiddleware := &jwt.JWTMiddleware{
		Key:           []byte("secret key"),
		Realm:         "jwt auth",
		Timeout:       time.Hour,
		MaxRefresh:    time.Hour * 24,
		Authenticator: Authenticate,
	}
	c := context{}
	c.Init()
	c.InitSchema()
	api := rest.NewApi()
	api.Use(rest.DefaultDevStack...)
	api.Use(&rest.IfMiddleware{
		Condition: func(request *rest.Request) bool {
			return request.URL.Path != "/login" &&
				request.URL.Path != "/register"
		},
		IfTrue: jwtMiddleware,
	})
	apiRouter, _ := rest.MakeRouter(
		rest.Post("/login", jwtMiddleware.LoginHandler),
		rest.Post("/register", Register),
		rest.Get("/refresh_token", jwtMiddleware.RefreshHandler),
	)
	api.SetApp(apiRouter)
	log.Fatal(http.ListenAndServe(":5000", api.MakeHandler()))
}

type context struct {
	DB gorm.DB
}

//Model -- GORM base model
type Model struct {
	ID        uint `gorm:"primaty_key"`
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
	Username string `json:"username"`
	Password string `json:"password"`
}

func Authenticate(userID string, password string) bool {
	c := context{}
	c.Init()
	user := User{}
	c.DB.Where("username = ?", userID).First(&user)
	storedPassword, _ := base64.URLEncoding.DecodeString(user.Password)
	salt, _ := base64.URLEncoding.DecodeString(user.Salt)
	checkPassword, _ := getPassword([]byte(password), salt)
	return bytes.Compare(storedPassword, checkPassword) == 0
}

func getPassword(password []byte, salt []byte) ([]byte, error) {
	return scrypt.Key(password, salt, 16384, 8, 1, 32)
}

func Register(w rest.ResponseWriter, r *rest.Request) {
	model := RegisterModel{}
	err := r.DecodeJsonPayload(&model)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	salt := make([]byte, 128)
	_, err = rand.Read(salt)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Println(salt)
	var password []byte
	password, err = getPassword([]byte(model.Password), salt)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Println(password)
	user := User{
		Username: model.Username,
		Password: base64.URLEncoding.EncodeToString(password),
		Salt:     base64.URLEncoding.EncodeToString(salt),
	}
	log.Println(user)
	c := context{}
	c.Init()

	c.DB.Create(&user)
	w.WriteHeader(http.StatusOK)
}

func (c *context) Init() {
	var err error
	c.DB, err = gorm.Open("postgres", "host=localhost port=5432 dbname=movier user=postgres password=123456 sslmode=disable")
	if err != nil {
		log.Fatal("Database connection error: %v", err)
	}
}

func (c *context) InitSchema() {
	c.DB.AutoMigrate(&User{})
}
