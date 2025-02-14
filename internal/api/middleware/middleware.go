package middleware

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Printf("Incoming request: %s %s FROM %s", c.Request.Method, c.Request.URL.Path, c.ClientIP())
		c.Next()
	}
}

func ExtractID(key string) gin.HandlerFunc {
	return func(c *gin.Context) {
		idParam := c.Param(key)

		if idParam == "" {
			c.Next()
			return
		}

		id, err := strconv.ParseInt(idParam, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
			c.Abort()
			return
		}

		c.Set(key, id)

		c.Next()
	}
}
