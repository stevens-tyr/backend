package cms

import (
	"github.com/gin-gonic/gin"
	"github.com/mongodb/mongo-go-driver/bson/primitive"
 
	"backend/forms"

	"github.com/stevens-tyr/tyr-gin"
)

func CourseAddUser(c *gin.Context) {
	cid, err := primitive.ObjectIDFromHex(c.Param("cid"))
	if err != nil {
		tyrgin.ErrorHandler(err, c, 500, gin.H{
			"status_code": 500,
			"error":       err.Error(),
		})
		return
	}

	var addUser forms.CourseAddUserForm
	if err = c.ShouldBindJSON(&addUser); err != nil {
		tyrgin.ErrorHandler(err, c, 500, gin.H{
			"status_code": 500,
			"error":       err.Error(),
		})
		return
	}

	user, err := um.FindOne(addUser.Email)
	if err != nil {
		tyrgin.ErrorHandler(err, c, 500, gin.H{
			"status_code": 500,
			"error":       err.Error(),
		})
		return
	}
	
	err = um.AddCourse(addUser.Level, cid, user.ID)
	if err != nil {
		tyrgin.ErrorHandler(err, c, 500, gin.H{
			"status_code": 500,
			"error":       err.Error(),
		})
		return
	}

	err = cm.AddUser(addUser.Level, user.ID, cid)
	if err != nil {
		tyrgin.ErrorHandler(err, c, 500, gin.H{
			"status_code": 500,
			"error":       err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"status_code": 200,
		"msg":         "User added.",
	})
}