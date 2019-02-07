package cms

import (
	"github.com/gin-gonic/gin"
	"github.com/mongodb/mongo-go-driver/bson/primitive"
 
	"github.com/stevens-tyr/tyr-gin"
)

// CourseAssignments is the function for a route to display all assignments a course has.
func CourseAssignments(c *gin.Context) {
	cid, err := primitive.ObjectIDFromHex(c.Param("cid"))
	if err != nil {
		tyrgin.ErrorHandler(err, c, 500, gin.H{
			"status_code": 500,
			"error":       err.Error(),
		})
		return
	}

	courses, err := cm.GetAssignments(cid)
	if err != nil {
		tyrgin.ErrorHandler(err, c, 400, gin.H{
			"status_code": 400,
			"error": err.Error(),
		})
		return
	} 
	c.JSON(200, gin.H{
		"status_code": 200,
		"msg":         "User's courses.",
		"courses":     courses,
	})
}
