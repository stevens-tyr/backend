package main

import (
	"github.com/gin-gonic/gin"
	"github.com/stevens-tyr/tyr-gin"
)

func main() {
	server := tyrgin.SetupRouter()

	server.Use(tyrgin.Logger())
	server.Use(gin.Recovery())

	server.Run(":5555")
}
