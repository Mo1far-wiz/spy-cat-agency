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

func GetAllMissions(context *gin.Context) {
	ctx := context.Request.Context()

	missions, err := application.App.Store.Mission.GetAllWithTargets(ctx)
	if err != nil {
		log.Println("error: ", err.Error())
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not get all missions."})
		return
	}

	context.JSON(http.StatusOK, missions)
}

func GetMissionByID(context *gin.Context) {
	ctx := context.Request.Context()

	id := context.GetInt64("missionID")

	mission, err := application.App.Store.Mission.GetByIDWithTargets(ctx, id)
	if err != nil {
		log.Println("error: ", err.Error())
		switch {
		case errors.Is(err, store.ErrorNotFound):
			context.JSON(http.StatusNotFound, gin.H{"message": "Mission not found."})
		default:
			context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not get mission."})
		}
		return
	}

	context.JSON(http.StatusOK, mission)
}

func CreateMission(context *gin.Context) {
	ctx := context.Request.Context()
	var mission store.Mission
	err := context.ShouldBindBodyWithJSON(&mission)
	if err != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{"message": "Could not parse request data."})
		return
	}

	err = application.App.Store.Mission.Create(ctx, &mission)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not create mission."})
	}

	context.JSON(http.StatusCreated, mission)
}

func UpdateMission(context *gin.Context) {
	ctx := context.Request.Context()

	var request requestMissionComplete
	err := context.ShouldBindBodyWithJSON(&request)
	if err != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{"message": "Could not parse request data."})
		return
	}

	id := context.GetInt64("missionID")

	targets, err := application.App.Store.Mission.GetAllMissionTargets(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrorNotFound):
			context.JSON(http.StatusNotFound, gin.H{"message": "Targets not found."})
		default:
			context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not get target."})
		}
		return
	}

	for _, t := range targets {
		if !t.IsComplete {
			context.JSON(http.StatusBadRequest, gin.H{"message": "spy have to complete all targets."})
			return
		}
	}

	mission := store.Mission{
		ID:         id,
		IsComplete: request.IsComplete,
	}

	err = application.App.Store.Mission.Update(ctx, &mission)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrorNotFound):
			context.JSON(http.StatusNotFound, gin.H{"message": "Mission not found."})
		default:
			context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not update mission."})
		}
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "updated"})
}

func DeleteMission(context *gin.Context) {
	ctx := context.Request.Context()
	id := context.GetInt64("missionID")

	mission, err := application.App.Store.Mission.GetByID(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrorNotFound):
			context.JSON(http.StatusNotFound, gin.H{"message": "Mission not found."})
		default:
			context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not get mission."})
		}
		return
	}

	if mission.CatID != nil {
		context.JSON(http.StatusNotFound, gin.H{"message": "Mission can not be deleted: spy already assigned."})
	}

	err = application.App.Store.Mission.Delete(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrorNotFound):
			context.JSON(http.StatusNotFound, gin.H{"message": "Mission not found."})
		default:
			context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not delete mission."})
		}
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

func AssignCatForMission(context *gin.Context) {
	ctx := context.Request.Context()
	missionID := context.GetInt64("missionID")
	catID := context.GetInt64("catID")

	cat, err := application.App.Store.Cat.GetByID(ctx, catID)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrorNotFound):
			context.JSON(http.StatusNotFound, gin.H{"message": "Cat not found."})
		default:
			context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not get cat."})
		}
		return
	}

	hasMissions, err := application.App.Store.Cat.HasIncompleteMission(ctx, catID)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrorNotFound):
			context.JSON(http.StatusNotFound, gin.H{"message": "Cat not found."})
		default:
			context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not get cat."})
		}
		return
	}

	if hasMissions {
		context.JSON(http.StatusNotFound, gin.H{"message": "Mission can not be assigned: spy has unfinished business."})
		return
	}

	hasSpy, err := application.App.Store.Mission.HasAssignedSpy(ctx, missionID)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrorNotFound):
			context.JSON(http.StatusNotFound, gin.H{"message": "Mission not found."})
		default:
			context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not get mission."})
		}
		return
	}

	if hasSpy {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Mission already has assigned spy."})
		return
	}

	err = application.App.Store.Mission.AssignCat(ctx, cat.ID, missionID)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrorNotFound):
			context.JSON(http.StatusNotFound, gin.H{"message": "Mission not found."})
		default:
			context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not assign mission."})
		}
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "assigned", "cat": cat})
}

func AddMissionTarget(context *gin.Context) {
	ctx := context.Request.Context()
	missionID := context.GetInt64("missionID")

	var target store.Target
	err := context.ShouldBindBodyWithJSON(&target)
	if err != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{"message": "Could not parse request data."})
		return
	}

	num, err := application.App.Store.Mission.GetTargetsQuantity(ctx, missionID)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrorNotFound):
			context.JSON(http.StatusNotFound, gin.H{"message": "Mission not found."})
		default:
			context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not get mission."})
		}
		return
	}

	if num == 3 {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Maximum number of targets reached."})
		return
	}

	err = application.App.Store.Mission.AddTarget(ctx, missionID, &target)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not add target."})
	}

	context.JSON(http.StatusOK, gin.H{"message": "added"})
}

func DeleteMissionTarget(context *gin.Context) {
	ctx := context.Request.Context()
	missionID := context.GetInt64("missionID")
	targetID := context.GetInt64("targetID")

	target, err := application.App.Store.Mission.GetTargetByID(ctx, targetID)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrorNotFound):
			context.JSON(http.StatusNotFound, gin.H{"message": "Target not found."})
		default:
			context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not get target."})
		}
		return
	}

	if target.IsComplete {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Target already completed."})
		return
	}

	num, err := application.App.Store.Mission.GetTargetsQuantity(ctx, missionID)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrorNotFound):
			context.JSON(http.StatusNotFound, gin.H{"message": "Mission not found."})
		default:
			context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not get mission."})
		}
		return
	}

	if num == 1 {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Minimum number of targets reached."})
		return
	}

	err = application.App.Store.Mission.RemoveTarget(ctx, targetID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not add target."})
	}

	context.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

func UpdateMissionTarget(context *gin.Context) {
	ctx := context.Request.Context()
	targetID := context.GetInt64("targetID")

	target, err := application.App.Store.Mission.GetTargetByID(ctx, targetID)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrorNotFound):
			context.JSON(http.StatusNotFound, gin.H{"message": "Target not found."})
		default:
			context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not get target."})
		}
		return
	}

	if target.IsComplete {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Target already completed."})
		return
	}

	err = application.App.Store.Mission.UpdateTarget(ctx, target)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update target."})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "updated"})
}

func AddNoteOnTarget(context *gin.Context) {
	ctx := context.Request.Context()
	targetID := context.GetInt64("targetID")

	target, err := application.App.Store.Mission.GetTargetByID(ctx, targetID)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrorNotFound):
			context.JSON(http.StatusNotFound, gin.H{"message": "Target not found."})
		default:
			context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not get target."})
		}
		return
	}

	if target.IsComplete {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Target already completed."})
		return
	}

	var note store.Note
	err = context.ShouldBindBodyWithJSON(&note)
	if err != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{"message": "Could not parse request data."})
		return
	}

	note.TargetID = target.ID
	note.MissionID = target.MissionID

	err = application.App.Store.Mission.AddNote(ctx, &note)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to add note."})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "note added"})
}
