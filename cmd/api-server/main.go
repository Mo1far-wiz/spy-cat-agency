package main

import (
	"spy-cat-agency/internal/api"

	"github.com/gin-gonic/gin"
)

func main() {
	server := gin.Default()

	api.Mount(server)

	server.Run(":8080")
}
