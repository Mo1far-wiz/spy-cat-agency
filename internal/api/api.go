package api

import (
	"spy-cat-agency/internal/api/handlers"
	"spy-cat-agency/internal/api/middleware"

	"github.com/gin-gonic/gin"
)

func Mount(router *gin.Engine) {
	router.Use(gin.Recovery())

	apiV1 := router.Group("/v1")

	apiV1.Use(middleware.Logger())

	cats := apiV1.Group("/cats")
	cats.Use(middleware.ExtractID("catID"))
	cats.GET("/", handlers.GetAllCats)         // get all
	cats.GET("/:catID", handlers.GetCatByID)   // get by id
	cats.POST("/", handlers.CreateCat)         // create
	cats.PUT("/:catID", handlers.UpdateCat)    // update
	cats.DELETE("/:catID", handlers.DeleteCat) // delete

	missions := apiV1.Group("/missions")
	missions.Use(middleware.ExtractID("missionID"))
	missions.GET("/", handlers.GetAllMissions)             // get all
	missions.POST("/", handlers.CreateMission)             // create
	missions.GET("/:missionID", handlers.GetMissionByID)   // get by id
	missions.PUT("/:missionID", handlers.UpdateMission)    // update
	missions.DELETE("/:missionID", handlers.DeleteMission) // delete

	catMission := missions.Group("/:missionID")
	catMission.Use(middleware.ExtractID("catID"))
	catMission.PUT("/:catID/assign", handlers.AssignCatForMission) // assign cat for mission

	targets := missions.Group("/targets")
	targets.Use(middleware.ExtractID("targetID"))
	targets.POST("/:missionID", handlers.AddMissionTarget)     // add mission target
	targets.POST("note/:targetID", handlers.AddNoteOnTarget)   // update note on target
	targets.PUT("/:targetID", handlers.UpdateMissionTarget)    // update target
	targets.DELETE("/:targetID", handlers.DeleteMissionTarget) // delete mission target
}
