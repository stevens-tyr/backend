package cms

import (
	"github.com/gin-gonic/gin"
)

func DeleteAssignment(c *gin.Context) {
	aid, _ := c.Get("aid")

	err := am.Delete(aid)
	if err != nil {
		c.Set("error", err)
		return
	}

	c.JSON(200, gin.H{
		"message": "Assignment Deleted.",
	})
}
