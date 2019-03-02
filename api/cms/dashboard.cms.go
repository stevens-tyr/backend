package cms

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

// Dashboard is the function for a route to display all course a user has.
func Dashboard(c *gin.Context) {
	uid, _ := c.Get("uid")
	fmt.Println("uid", uid)

	courses, err := um.GetCourses(uid)
	if err != nil {
		c.Set("error", err)
		return
	}

	c.JSON(200, gin.H{
		"status_code": 200,
		"msg":         "User's courses.",
		"courses":     courses,
	})
}
