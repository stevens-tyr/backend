package cms

import (
	"github.com/gin-gonic/gin"

	"backend/errors"
)

func UpdateGrade(c *gin.Context) {
	sid, _ := c.Get("sid")

	sub, err := sm.Get(sid)
	if err != nil {
		c.Set("error", err)
		return
	}
	
	var up map[string]interface{}
	errs := c.ShouldBind(&up)
	if errs != nil {
		c.Set("error", errors.ErrorInvalidJSON)
		return
	}
	
	if val, ok := up["studentFacingPass"]; ok {
		sub.Cases.StudentFacing.Pass = val.(int)
	}
	if val, ok := up["studentFacingFail"]; ok {
		sub.Cases.StudentFacing.Fail = val.(int)
	}
	if val, ok := up["adminFacingPass"]; ok {
		sub.Cases.AdminFacing.Pass = val.(int)
	}
	if val, ok := up["adminFacingFail"]; ok {
		sub.Cases.AdminFacing.Fail = val.(int)
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
