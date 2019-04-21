package cms

import (
	"github.com/gin-gonic/gin"

	"backend/errors"
	"backend/forms"
)

func UpdateGrade(c *gin.Context) {
	sid, _ := c.Get("sid")

	sub, err := sm.Get(sid)
	if err != nil {
		c.Set("error", err)
		return
	}
	
	var up forms.UpdateSubmissionGradeForm
	errs := c.ShouldBind(&up)
	if errs != nil {
		c.Set("error", errors.ErrorInvalidJSON)
		return
	}

	if up.StudentFacingPass != nil {
		sub.Cases.StudentFacing.Pass = *up.StudentFacingPass
	}
	if up.StudentFacingFail != nil {
		sub.Cases.StudentFacing.Fail = *up.StudentFacingFail
	}
	if up.AdminFacingPass != nil {
		sub.Cases.AdminFacing.Pass = *up.AdminFacingPass
	}
	if up.AdminFacingFail != nil {
		sub.Cases.AdminFacing.Fail = *up.AdminFacingFail
	}
	
	err = sm.UpdateGrade(
		sid,
		sub.Cases.StudentFacing.Pass,
		sub.Cases.StudentFacing.Fail,
		sub.Cases.AdminFacing.Pass,
		sub.Cases.AdminFacing.Fail,
	)
	if err != nil{
		c.Set("error", err)
		return
	}
	
	c.JSON(200, gin.H{
		"message": "Submission Grade Updated.",
	})
}
