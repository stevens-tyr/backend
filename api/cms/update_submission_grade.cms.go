package cms

import (
	submodels "backend/models/cmsmodels/submissionmodels"

	"github.com/gin-gonic/gin"
)

// UpdateGrade will be called by court_herald to update the grade from brian
func UpdateGrade(c *gin.Context) {
	sid, _ := c.Get("sid")

	var testResults []submodels.WorkerResult
	c.BindJSON(&testResults)

	err := sm.UpdateGrade(sid, testResults)
	if err != nil {
		c.Set("error", err)
		return
	}
	c.JSON(200, gin.H{
		"message": "Submission Grade Updated.",
	})
}

// UpdateGradeError will edit the grade if an error is encountered while grading
func UpdateGradeError(c *gin.Context) {
	sid, _ := c.Get("sid")

	err := sm.UpdateError(sid)
	if err != nil {
		c.Set("error", err)
		return
	}
  
	c.JSON(200, gin.H{
		"message": "Submission Error Update.",
	})
}
