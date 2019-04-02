package cms

import (
	"github.com/gin-gonic/gin"

	"backend/errors"
)

func UpdateCourse(c *gin.Context) {
	cid, _ := c.Get("cid")

	course, err := cm.GetByID(cid)
	if err != nil {
		c.Set("error", err)
		return
	}

	var up map[string]interface{}
	errs := c.ShouldBind(&up)
	if errs != nil {
		c.Set("error", errors.ErrorInvlaidJSON)
		return
	}
	
	if val, ok := up["department"]; ok { 
		course.Department = val.(string)
	}
	if val, ok := up["department"]; ok {
		course.LongName = val.(string)
	}
	if val, ok := up["department"]; ok {
		course.Number = val.(int)
	}
	if val, ok := up["department"]; ok {
		course.Section = val.(string)
	}
	if val, ok := up["department"]; ok {
		course.Semester = val.(string)
	}

	err = cm.Update(*course)
	if err != nil {
		c.Set("error", err)
		return
	}
	
	c.JSON(200, gin.H{
		"message": "Course Updated.",
	})
}
