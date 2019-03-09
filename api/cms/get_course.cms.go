package cms

import (
	"github.com/gin-gonic/gin"
)

func GetCourse(c *gin.Context) {
	cid, _ := c.Get("cid")
	uid, _ := c.Get("uid")

	course, err := cm.Get(cid, uid)
	if err != nil {
		c.Set("error", err)
		return
	}

	c.JSON(200, gin.H{
		"status_code": 200,
		"msg":         "Course Info.",
		"course":  course,
	})
}