package handlers

import (
	"errors"
	"net/http"
	"spy-cat-agency/internal/application"
	"spy-cat-agency/internal/store"

	"github.com/gin-gonic/gin"
)

type request struct {
	Salary float64 `json:"salary" validate:"required"`
}

func GetAllCats(context *gin.Context) {
	ctx := context.Request.Context()

	cats, err := application.App.Store.Cat.GetAll(ctx)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not get all cats."})
		return
	}

	context.JSON(http.StatusOK, cats)
}

func GetCatByID(context *gin.Context) {
	ctx := context.Request.Context()

	id := context.GetInt64("catID")
	cat, err := application.App.Store.Cat.GetByID(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrorNotFound):
			context.JSON(http.StatusNotFound, gin.H{"message": "Cat not found."})
		default:
			context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not get cat."})
		}
		return
	}

	context.JSON(http.StatusOK, cat)
}

func CreateCat(context *gin.Context) {
	ctx := context.Request.Context()
	var cat store.Cat
	err := context.ShouldBindBodyWithJSON(&cat)
	if err != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{"message": "Could not parse request data."})
		return
	}

	err = application.App.Store.Cat.Create(ctx, &cat)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not create cat."})
	}

	context.JSON(http.StatusCreated, cat)
}

func UpdateCat(context *gin.Context) {
	ctx := context.Request.Context()

	var request request
	err := context.ShouldBindBodyWithJSON(&request)
	if err != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{"message": "Could not parse request data."})
		return
	}

	id := context.GetInt64("catID")

	cat := store.Cat{
		ID:     id,
		Salary: request.Salary,
	}

	err = application.App.Store.Cat.Update(ctx, &cat)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrorNotFound):
			context.JSON(http.StatusNotFound, gin.H{"message": "Cat not found."})
		default:
			context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not update cat."})
		}
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "updated"})
}

func DeleteCat(context *gin.Context) {
	ctx := context.Request.Context()
	id := context.GetInt64("catID")

	err := application.App.Store.Cat.Delete(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrorNotFound):
			context.JSON(http.StatusNotFound, gin.H{"message": "Cat not found."})
		default:
			context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not delete cat."})
		}
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "deleted"})
}
