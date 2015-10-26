package handlers

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"gopkg.in/mgo.v2/bson"

	"github.com/gin-gonic/gin"
	"github.com/grffn/movier/db"
	"github.com/grffn/movier/models"
)

//CreateHandler Create new item
func CreateHandler(context *gin.Context, database *db.Context) {
	model := models.DocModel{}
	err := context.BindJSON(&model)
	if err != nil {
		context.AbortWithStatus(http.StatusBadRequest)
		return
	}
	userID := context.MustGet("userId").(string)
	user := database.FindUser(userID)
	document := db.Document{
		ID:       bson.NewObjectId(),
		UserID:   user.ID,
		Name:     model.Name,
		Tags:     model.Tags,
		Category: model.Category,
		URL:      model.URL,
		MimeType: model.MimeType,
		Authors:  model.Authors,
	}
	database.NewDocument(document)
}

//SignHandler Sign AWS reqwest
func SignHandler(context *gin.Context) {
	fileName := context.Query("filename")
	fileType := context.Query("filetype")

	accessKey := os.Getenv("AWS_ACCESS_KEY_ID")
	secretKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
	bucket := os.Getenv("S3_BUCKET_NAME")

	fileName = url.QueryEscape(fileName)
	expires := time.Now().AddDate(0, 0, 1).Unix()
	amzHeaders := "x-amz-acl:public-read"

	stringToSign := fmt.Sprintf("PUT\n\n%s\n%d\n%s\n/%s/%s", fileType, expires, amzHeaders, bucket, fileName)
	hasher := hmac.New(sha1.New, []byte(accessKey))
	hasher.Write([]byte(stringToSign))

	signature := url.QueryEscape(base64.URLEncoding.EncodeToString(hasher.Sum(nil)))
	query := url.Values{
		"AWSAccessKeyId": []string{secretKey},
		"Expires":        []string{strconv.FormatInt(expires, 10)},
		"Signature":      []string{signature},
	}
	url := fmt.Sprintf("https://%s.s3.amazonaws.com/%s", bucket, fileName)
	context.JSON(http.StatusOK, gin.H{"signedRequest": url + "?" + query.Encode(), "url": url})
}
