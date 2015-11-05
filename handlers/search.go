package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/grffn/movier/db"
)

func DocumentsHandler(context *gin.Context, database *db.Context) {
	result, err := database.Documents()
	if err != nil {
		context.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	context.JSON(http.StatusOK, result)
}

func CategoriesHandler(context *gin.Context, database *db.Context) {
	result, err := database.Categories()
	if err != nil {
		context.AbortWithError(http.StatusInternalServerError, err)
	}
	context.JSON(http.StatusOK, result)
}

func DocumentsByCategoryHandler(context *gin.Context, database *db.Context) {
	category := context.Param("category")
	result, err := database.DocumentsByCategory(category)
	if err != nil {
		context.AbortWithError(http.StatusInternalServerError, err)
	}
	context.JSON(http.StatusOK, result)
}

func DocumentsByTagHandler(context *gin.Context, database *db.Context) {
	tag := context.Param("tag")
	println(tag)
	result, err := database.DocumentsByTag(tag)
	if err != nil {
		context.AbortWithError(http.StatusInternalServerError, err)
	}
	context.JSON(http.StatusOK, result)
}
