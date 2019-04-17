package cms

import (
	"github.com/gin-gonic/gin"
)

func DeleteCourse(c *gin.Context) {
	cid, _ := c.Get("cid")

	err := cm.Delete(cid)
	if err != nil {
		c.Set("error", err)
		return
	}

	c.JSON(200, gin.H{
		"message": "Course Deleted.",
	})
}
