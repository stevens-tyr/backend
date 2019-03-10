package cms

import (
	"github.com/gin-gonic/gin"
)

func GetAssignment(c *gin.Context) {
	// lets verify assignment is in course in future?
	aid, _ := c.Get("aid")
	uid, _ := c.Get("uid")
	role, _ := c.Get("role")
 
	assignment, err := am.GetFull(aid, uid, role.(string))
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
