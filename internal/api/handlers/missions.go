package handlers

import (
	"errors"
	"log"
	"net/http"
	"spy-cat-agency/internal/application"
	"spy-cat-agency/internal/store"

	"github.com/gin-gonic/gin"
)

type requestMissionComplete struct {
	IsComplete bool `json:"is_complete" validate:"required"`
}

type response struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func newResponse(message string, data ...interface{}) response {
	if len(data) > 0 {
		return response{Message: message, Data: data[0]}
	}
	return response{Message: message}
}

func logError(err error, message string) {
	log.Printf("ERROR: %s: %v", message, err)
}

func GetAllMissions(c *gin.Context) {
	missions, err := application.App.Store.Mission.GetAllWithTargets(c.Request.Context())
	if err != nil {
		logError(err, "failed to get all missions")
		c.JSON(http.StatusInternalServerError, newResponse("Could not get all missions"))
		return
	}
	c.JSON(http.StatusOK, missions)
}

func GetMissionByID(c *gin.Context) {
	mission, err := application.App.Store.Mission.GetByIDWithTargets(c.Request.Context(), c.GetInt64("missionID"))
	if err != nil {
		logError(err, "failed to get mission by ID")
		status := http.StatusInternalServerError
		message := "Could not get mission"

		if errors.Is(err, store.ErrorNotFound) {
			status = http.StatusNotFound
			message = "Mission not found"
		}

		c.JSON(status, newResponse(message))
		return
	}
	c.JSON(http.StatusOK, mission)
}

func CreateMission(c *gin.Context) {
	var mission store.Mission
	if err := c.ShouldBindJSON(&mission); err != nil {
		logError(err, "failed to parse mission data")
		c.JSON(http.StatusUnprocessableEntity, newResponse("Could not parse request data"))
		return
	}

	if err := application.App.Store.Mission.Create(c.Request.Context(), &mission); err != nil {
		logError(err, "failed to create mission")
		c.JSON(http.StatusInternalServerError, newResponse("Could not create mission"))
		return
	}
	c.JSON(http.StatusCreated, mission)
}

func UpdateMission(c *gin.Context) {
	var request requestMissionComplete
	if err := c.ShouldBindJSON(&request); err != nil {
		logError(err, "failed to parse mission update data")
		c.JSON(http.StatusUnprocessableEntity, newResponse("Could not parse request data"))
		return
	}

	ctx := c.Request.Context()
	missionID := c.GetInt64("missionID")

	targets, err := application.App.Store.Mission.GetAllMissionTargets(ctx, missionID)
	if err != nil {
		logError(err, "failed to get mission targets")
		status := http.StatusInternalServerError
		message := "Could not get targets"

		if errors.Is(err, store.ErrorNotFound) {
			status = http.StatusNotFound
			message = "Targets not found"
		}

		c.JSON(status, newResponse(message))
		return
	}

	for _, t := range targets {
		if !t.IsComplete {
			c.JSON(http.StatusBadRequest, newResponse("All targets must be completed first"))
			return
		}
	}

	mission := store.Mission{
		ID:         missionID,
		IsComplete: request.IsComplete,
	}

	if err := application.App.Store.Mission.Update(ctx, &mission); err != nil {
		logError(err, "failed to update mission")
		status := http.StatusInternalServerError
		message := "Could not update mission"

		if errors.Is(err, store.ErrorNotFound) {
			status = http.StatusNotFound
			message = "Mission not found"
		}

		c.JSON(status, newResponse(message))
		return
	}
	c.JSON(http.StatusOK, newResponse("Mission updated"))
}

func DeleteMission(c *gin.Context) {
	ctx := c.Request.Context()
	missionID := c.GetInt64("missionID")

	mission, err := application.App.Store.Mission.GetByID(ctx, missionID)
	if err != nil {
		logError(err, "failed to get mission for deletion")
		status := http.StatusInternalServerError
		message := "Could not get mission"

		if errors.Is(err, store.ErrorNotFound) {
			status = http.StatusNotFound
			message = "Mission not found"
		}

		c.JSON(status, newResponse(message))
		return
	}

	if mission.CatID != nil {
		c.JSON(http.StatusBadRequest, newResponse("Cannot delete mission: spy already assigned"))
		return
	}

	if err := application.App.Store.Mission.Delete(ctx, missionID); err != nil {
		logError(err, "failed to delete mission")
		c.JSON(http.StatusInternalServerError, newResponse("Could not delete mission"))
		return
	}
	c.JSON(http.StatusOK, newResponse("Mission deleted"))
}

func AssignCatForMission(c *gin.Context) {
	ctx := c.Request.Context()
	missionID := c.GetInt64("missionID")
	catID := c.GetInt64("catID")

	cat, err := application.App.Store.Cat.GetByID(ctx, catID)
	if err != nil {
		logError(err, "failed to get cat for mission assignment")
		status := http.StatusInternalServerError
		message := "Could not get cat"

		if errors.Is(err, store.ErrorNotFound) {
			status = http.StatusNotFound
			message = "Cat not found"
		}

		c.JSON(status, newResponse(message))
		return
	}

	hasMissions, err := application.App.Store.Cat.HasIncompleteMission(ctx, catID)
	if err != nil {
		logError(err, "failed to check cat's incomplete missions")
		c.JSON(http.StatusInternalServerError, newResponse("Could not check cat missions"))
		return
	}

	if hasMissions {
		c.JSON(http.StatusBadRequest, newResponse("Cannot assign mission: spy has unfinished business"))
		return
	}

	hasSpy, err := application.App.Store.Mission.HasAssignedSpy(ctx, missionID)
	if err != nil {
		logError(err, "failed to check mission's assigned spy")
		c.JSON(http.StatusInternalServerError, newResponse("Could not check mission assignment"))
		return
	}

	if hasSpy {
		c.JSON(http.StatusBadRequest, newResponse("Mission already has an assigned spy"))
		return
	}

	if err := application.App.Store.Mission.AssignCat(ctx, cat.ID, missionID); err != nil {
		logError(err, "failed to assign cat to mission")
		c.JSON(http.StatusInternalServerError, newResponse("Could not assign mission"))
		return
	}
	c.JSON(http.StatusOK, newResponse("Mission assigned", cat))
}

func AddMissionTarget(c *gin.Context) {
	var target store.Target
	if err := c.ShouldBindJSON(&target); err != nil {
		logError(err, "failed to parse target data")
		c.JSON(http.StatusUnprocessableEntity, newResponse("Could not parse request data"))
		return
	}

	ctx := c.Request.Context()
	missionID := c.GetInt64("missionID")

	count, err := application.App.Store.Mission.GetTargetsQuantity(ctx, missionID)
	if err != nil {
		logError(err, "failed to get targets count")
		c.JSON(http.StatusInternalServerError, newResponse("Could not get targets count"))
		return
	}

	if count >= 3 {
		c.JSON(http.StatusBadRequest, newResponse("Maximum number of targets (3) reached"))
		return
	}

	if err := application.App.Store.Mission.AddTarget(ctx, missionID, &target); err != nil {
		logError(err, "failed to add target to mission")
		c.JSON(http.StatusInternalServerError, newResponse("Could not add target"))
		return
	}
	c.JSON(http.StatusOK, newResponse("Target added"))
}

func DeleteMissionTarget(c *gin.Context) {
	ctx := c.Request.Context()
	missionID := c.GetInt64("missionID")
	targetID := c.GetInt64("targetID")

	target, err := application.App.Store.Mission.GetTargetByID(ctx, targetID)
	if err != nil {
		logError(err, "failed to get target for deletion")
		status := http.StatusInternalServerError
		message := "Could not get target"

		if errors.Is(err, store.ErrorNotFound) {
			status = http.StatusNotFound
			message = "Target not found"
		}

		c.JSON(status, newResponse(message))
		return
	}

	if target.IsComplete {
		c.JSON(http.StatusBadRequest, newResponse("Cannot delete completed target"))
		return
	}

	count, err := application.App.Store.Mission.GetTargetsQuantity(ctx, missionID)
	if err != nil {
		logError(err, "failed to get targets count")
		c.JSON(http.StatusInternalServerError, newResponse("Could not get targets count"))
		return
	}

	if count <= 1 {
		c.JSON(http.StatusBadRequest, newResponse("Cannot delete last target"))
		return
	}

	if err := application.App.Store.Mission.RemoveTarget(ctx, targetID); err != nil {
		logError(err, "failed to delete target")
		c.JSON(http.StatusInternalServerError, newResponse("Could not delete target"))
		return
	}
	c.JSON(http.StatusOK, newResponse("Target deleted"))
}

func UpdateMissionTarget(c *gin.Context) {
	ctx := c.Request.Context()
	targetID := c.GetInt64("targetID")

	target, err := application.App.Store.Mission.GetTargetByID(ctx, targetID)
	if err != nil {
		logError(err, "failed to get target for update")
		status := http.StatusInternalServerError
		message := "Could not get target"

		if errors.Is(err, store.ErrorNotFound) {
			status = http.StatusNotFound
			message = "Target not found"
		}

		c.JSON(status, newResponse(message))
		return
	}

	if target.IsComplete {
		c.JSON(http.StatusBadRequest, newResponse("Cannot update completed target"))
		return
	}

	var req requestMissionComplete
	if err := c.ShouldBindJSON(&req); err != nil {
		logError(err, "failed to parse target data")
		c.JSON(http.StatusUnprocessableEntity, newResponse("Could not parse request data"))
		return
	}

	target.IsComplete = req.IsComplete

	mission, err := application.App.Store.Mission.GetByID(ctx, target.MissionID)
	if err != nil {
		logError(err, "failed to get mission for target update")
		c.JSON(http.StatusInternalServerError, newResponse("Could not verify mission status"))
		return
	}

	if mission.CatID == nil {
		c.JSON(http.StatusBadRequest, newResponse("Cannot update target: no spy assigned to mission"))
		return
	}

	if err := application.App.Store.Mission.UpdateTarget(ctx, target); err != nil {
		logError(err, "failed to update target")
		c.JSON(http.StatusInternalServerError, newResponse("Could not update target"))
		return
	}
	c.JSON(http.StatusOK, newResponse("Target updated"))
}

func AddNoteOnTarget(c *gin.Context) {
	ctx := c.Request.Context()
	targetID := c.GetInt64("targetID")

	target, err := application.App.Store.Mission.GetTargetByID(ctx, targetID)
	if err != nil {
		logError(err, "failed to get target for adding note")
		status := http.StatusInternalServerError
		message := "Could not get target"

		if errors.Is(err, store.ErrorNotFound) {
			status = http.StatusNotFound
			message = "Target not found"
		}

		c.JSON(status, newResponse(message))
		return
	}

	if target.IsComplete {
		c.JSON(http.StatusBadRequest, newResponse("Cannot add note to completed target"))
		return
	}

	var note store.Note
	if err := c.ShouldBindJSON(&note); err != nil {
		logError(err, "failed to parse note data")
		c.JSON(http.StatusUnprocessableEntity, newResponse("Could not parse request data"))
		return
	}

	note.TargetID = target.ID

	if err := application.App.Store.Mission.AddNote(ctx, &note); err != nil {
		logError(err, "failed to add note to target")
		c.JSON(http.StatusInternalServerError, newResponse("Could not add note"))
		return
	}
	c.JSON(http.StatusOK, newResponse("Note added"))
}
