package cms

import (
	"github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"

	"backend/forms"
)

// Dashboard is the function for a route to display all course a user has.
func Dashboard(c *gin.Context) {
	uid, _ := c.Get("uid")

	claims := jwt.ExtractClaims(c)
	courses, err := um.GetCourses(uid, claims["courses"].(map[string]interface{}))
	if err != nil {
		c.Set("error", err)
		return
	}

	var assignments []forms.AssignmentAggQuery
	for _, course := range courses {
		courseAssignments, err := cm.GetAssignments(course.ID, course.Role)
		if err != nil {
			c.Set("error", err)
			return
		}

		assignments = append(assignments, courseAssignments...)
	}

	submissions, err := sm.GetUsersSubmissionsLimited(uid, 5)
	if err != nil {
		c.Set("error", err)
		return
	}

	c.JSON(200, gin.H{
		"status_code": 200,
		"msg":         "User's Info.",
		"courses":     courses,
		"assignments": assignments,
		"mostRecentSubmissions": submissions,
	})
}