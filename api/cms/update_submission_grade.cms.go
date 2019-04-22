package cms

import (
	submodels "backend/models/cmsmodels/submissionmodels"
	"fmt"

	"github.com/gin-gonic/gin"
)

func UpdateGrade(c *gin.Context) {
	sid, _ := c.Get("sid")

	var testResults []submodels.WorkerResult
	c.BindJSON(&testResults)
	fmt.Println(testResults[0])

	err := sm.UpdateGrade(
		sid,
		testResults,
	)
	if err != nil {
		c.Set("error", err)
		return
	}

	c.JSON(200, gin.H{
		"message": "Submission Grade Updated.",
	})
}
