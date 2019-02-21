package auth

import (
	"time"

	"github.com/gin-gonic/gin"
)

func TokenResponse(c *gin.Context, code int, token string, expire time.Time) {
	c.JSON(200, gin.H{
		"token":  token,
		"expire": expire.Format(time.RFC3339),
	})
}
