package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetAllMissions(context *gin.Context) {
	context.JSON(http.StatusNotImplemented, gin.H{"message": "GetAllMissions"})
}

func GetMissionByID(context *gin.Context) {
	context.JSON(http.StatusNotImplemented, gin.H{"message": "GetMissionByID"})
}

func CreateMission(context *gin.Context) {
	context.JSON(http.StatusNotImplemented, gin.H{"message": "CreateMission"})
}

func UpdateMission(context *gin.Context) {
	context.JSON(http.StatusNotImplemented, gin.H{"message": "UpdateMission"})
}

func DeleteMission(context *gin.Context) {
	context.JSON(http.StatusNotImplemented, gin.H{"message": "DeleteMission"})
}

func AssignCatForMission(context *gin.Context) {
	context.JSON(http.StatusNotImplemented, gin.H{"message": "AssignCatForMission"})
}

func RemoveCatForMission(context *gin.Context) {
	context.JSON(http.StatusNotImplemented, gin.H{"message": "RemoveCatForMission"})
}

func AddMissionTarget(context *gin.Context) {
	context.JSON(http.StatusNotImplemented, gin.H{"message": "AddMissionTarget"})
}

func DeleteMissionTarget(context *gin.Context) {
	context.JSON(http.StatusNotImplemented, gin.H{"message": "DeleteMissionTarget"})
}

func AddNoteOnTarget(context *gin.Context) {
	context.JSON(http.StatusNotImplemented, gin.H{"message": "AddNoteOnTarget"})
}
