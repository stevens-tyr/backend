package cms

import (
	"github.com/gin-gonic/gin"
)

// CourseAssignments is the function for a route to display all assignments a course has.
func CourseAssignments(c *gin.Context) {
	cid, _ := c.Get("cid")

	assignments, err := cm.GetAssignments(cid)
	if err != nil {
		c.Set("error", err)
		return
	}
	c.JSON(200, gin.H{
		"message": "Course assignments.",
		"assignments": assignments,
	})
}
