package auth

import (
	"github.com/gin-gonic/gin"
)

func Check(c *gin.Context) {
	c.JSON(200, gin.H{
		"status_code": 200,
		"msg":         "Success.",
	})
}