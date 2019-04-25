package cms

import (
	submodels "backend/models/cmsmodels/submissionmodels"

	"github.com/gin-gonic/gin"
)

// Update Grade will be called by court_herald to update the grade from brian
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
