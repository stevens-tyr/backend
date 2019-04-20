package cms

import (
	"github.com/gin-gonic/gin"

	"backend/errors"
	"backend/forms"
)

func CourseAddUser(c *gin.Context) {
	cid, _ := c.Get("cid")

	var addUser forms.CourseAddUserForm
	if err := c.ShouldBindJSON(&addUser); err != nil {
		c.Set("error", errors.ErrorInvalidJSON)
		return
	}

	user, err := um.FindOne(addUser.Email)
	if err != nil {
		c.Set("error", err)
		return
	}

	err = um.AddCourse(addUser.Level, cid, user.ID)
	if err != nil {
		c.Set("error", err)
		return
	}

	err = cm.AddUser(addUser.Level, user.ID, cid)
	if err != nil {
		c.Set("error", err)

		return
	}

	c.JSON(200, gin.H{
		"msg": "User added.",
	})
}
