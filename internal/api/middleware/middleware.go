package middleware

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		log.Printf("Incoming request: %s %s", c.Request.Method, c.Request.URL.Path)

		c.Next()

		duration := time.Since(start)
		log.Printf("Response status: %d, Duration: %v", c.Writer.Status(), duration)
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
