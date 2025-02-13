package api

import (
	"spy-cat-agency/internal/api/handlers"
	"time"

	"github.com/gin-gonic/gin"
)

type Application struct {
	Config Config
}

type Config struct {
	Addr         string
	WriteTimeout time.Duration
	ReadTimeout  time.Duration
	IdleTimeout  time.Duration
}

func Mount(server *gin.Engine) {
	apiV1 := server.Group("/v1")

	apiV1.GET("/healthy")

	cats := apiV1.Group("/cats")
	cats.GET("/", handlers.GetAllCats)      // get all
	cats.GET("/:id", handlers.GetCatByID)   // get by id
	cats.POST("/", handlers.CreateCat)      // create
	cats.PUT("/:id", handlers.UpdateCat)    // update
	cats.DELETE("/:id", handlers.DeleteCat) // delete

	missions := apiV1.Group("/missions")
	missions.GET("/", handlers.GetAllMissions)      // get all
	missions.GET("/:id", handlers.GetMissionByID)   // get by id
	missions.POST("/", handlers.CreateMission)      // create
	missions.PUT("/:id", handlers.UpdateMission)    // update
	missions.DELETE("/:id", handlers.DeleteMission) // delete

	missions.PUT("/:id/assign", handlers.AssignCatForMission)    // assign	cat for mission
	missions.DELETE("/:id/remove", handlers.RemoveCatForMission) // remove	cat from mission

	targets := missions.Group("/:id/targets")
	targets.POST("/", handlers.AddMissionTarget)         // add mission target
	targets.PUT("/:id/note", handlers.AddNoteOnTarget)   // update note on target
	targets.DELETE("/:id", handlers.DeleteMissionTarget) // delete mission target
}
