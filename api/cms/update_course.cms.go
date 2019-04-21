package cms

import (
	"github.com/gin-gonic/gin"

	"backend/errors"
	"backend/forms"
)

func UpdateCourse(c *gin.Context) {
	cid, _ := c.Get("cid")

	course, err := cm.GetByID(cid)
	if err != nil {
		c.Set("error", err)
		return
	}

	var up forms.UpdateCourseForm
	errs := c.ShouldBind(&up)
	if errs != nil {
		c.Set("error", errors.ErrorInvalidJSON)
		return
	}

	if up.Department != nil {
		course.Department = *up.Department
	}
	if up.LongName != nil {
		course.LongName = *up.LongName
	}
	if up.Number != nil {
		course.Number = *up.Number
	}
	if up.Section != nil {
		course.Section = *up.Section
	}
	if up.Semester != nil {
		course.Semester = *up.Semester
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
