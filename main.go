package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	server := gin.Default()

	server.Use(gin.Recovery())

	server.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": 200,
		})
	})

	server.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{
			"status":  404,
			"message": "page not found",
		})
	})

	server.Run(":5000")
}
