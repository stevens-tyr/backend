package cms

import (
	"github.com/gin-gonic/gin"
)

func GetSubmission(c *gin.Context) {
	sid, _ := c.Get("sid")

	submission, err := sm.Get(sid)
	if err != nil {
		c.Set("error", err)
		return
	}

	c.JSON(200, gin.H{
		"status_code": 200,
		"msg":         "submission.",
		"submission":  submission,
	})
}
