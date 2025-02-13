package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetAllCats(context *gin.Context) {
	context.JSON(http.StatusNotImplemented, gin.H{"message": "GetAllCats"})
}

func GetCatByID(context *gin.Context) {
	context.JSON(http.StatusNotImplemented, gin.H{"message": "GetCatByID"})
}

func CreateCat(context *gin.Context) {
	context.JSON(http.StatusNotImplemented, gin.H{"message": "CreateCat"})
}

func UpdateCat(context *gin.Context) {
	context.JSON(http.StatusNotImplemented, gin.H{"message": "UpdateCat"})
}

func DeleteCat(context *gin.Context) {
	context.JSON(http.StatusNotImplemented, gin.H{"message": "DeleteCat"})
}
