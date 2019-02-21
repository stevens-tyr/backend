package cms

import (
	"github.com/gin-gonic/gin"
)

func GetAssignment(c *gin.Context) {
	aid, _ := c.Get("aid")

	assignment, err := am.Get(aid)
	if err != nil {
		c.Set("error", err)
		return
	}

	c.JSON(200, gin.H{
		"status_code": 200,
		"msg":         "assignment.",
		"assignment":  assignment,
	})
}
