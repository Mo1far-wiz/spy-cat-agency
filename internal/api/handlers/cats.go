package handlers

import (
	"errors"
	"net/http"
	"spy-cat-agency/internal/application"
	"spy-cat-agency/internal/breed"
	"spy-cat-agency/internal/store"

	"github.com/gin-gonic/gin"
)

type requestChangeSalary struct {
	Salary float64 `json:"salary" validate:"required"`
}

type errorResponse struct {
	Message string `json:"message"`
}

func newErrorResponse(message string) errorResponse {
	return errorResponse{Message: message}
}

func GetAllCats(c *gin.Context) {
	cats, err := application.App.Store.Cat.GetAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, newErrorResponse("Could not get all cats"))
		return
	}
	c.JSON(http.StatusOK, cats)
}

func GetCatByID(c *gin.Context) {
	cat, err := application.App.Store.Cat.GetByID(c.Request.Context(), c.GetInt64("catID"))
	if err != nil {
		status := http.StatusInternalServerError
		message := "Could not get cat"

		if errors.Is(err, store.ErrorNotFound) {
			status = http.StatusNotFound
			message = "Cat not found"
		}

		c.JSON(status, newErrorResponse(message))
		return
	}
	c.JSON(http.StatusOK, cat)
}

func CreateCat(c *gin.Context) {
	var cat store.Cat
	if err := c.ShouldBindJSON(&cat); err != nil {
		c.JSON(http.StatusUnprocessableEntity, newErrorResponse("Could not parse request data"))
		return
	}

	exists, err := breed.ValidateCatBreed(cat.Breed)
	if err != nil {
		c.JSON(http.StatusInternalServerError, newErrorResponse("Could not validate breed"))
		return
	}

	if !exists {
		c.JSON(http.StatusBadRequest, newErrorResponse("Invalid breed"))
		return
	}

	if err := application.App.Store.Cat.Create(c.Request.Context(), &cat); err != nil {
		c.JSON(http.StatusInternalServerError, newErrorResponse("Could not create cat"))
		return
	}
	c.JSON(http.StatusCreated, cat)
}

func UpdateCat(c *gin.Context) {
	var request requestChangeSalary
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusUnprocessableEntity, newErrorResponse("Could not parse request data"))
		return
	}

	cat := store.Cat{
		ID:     c.GetInt64("catID"),
		Salary: request.Salary,
	}

	if err := application.App.Store.Cat.Update(c.Request.Context(), &cat); err != nil {
		status := http.StatusInternalServerError
		message := "Could not update cat"

		if errors.Is(err, store.ErrorNotFound) {
			status = http.StatusNotFound
			message = "Cat not found"
		}

		c.JSON(status, newErrorResponse(message))
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "updated"})
}

func DeleteCat(c *gin.Context) {
	if err := application.App.Store.Cat.Delete(c.Request.Context(), c.GetInt64("catID")); err != nil {
		status := http.StatusInternalServerError
		message := "Could not delete cat"

		if errors.Is(err, store.ErrorNotFound) {
			status = http.StatusNotFound
			message = "Cat not found"
		}

		c.JSON(status, newErrorResponse(message))
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}
