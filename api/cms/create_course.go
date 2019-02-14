package cms

import (
	jwt "github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
	"github.com/mongodb/mongo-go-driver/bson/primitive"
 
	"backend/forms"

	"github.com/stevens-tyr/tyr-gin"
)

func CreateCourse(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	uid, err := primitive.ObjectIDFromHex(claims["uid"].(string))
	if err != nil {
		tyrgin.ErrorHandler(err, c, 500, gin.H{
			"status_code": 500,
			"error":       err.Error(),
		})
		return
	}

	var createCourse forms.CreateCourseForm
	if err = c.ShouldBindJSON(&createCourse); err != nil {
		tyrgin.ErrorHandler(err, c, 500, gin.H{
			"status_code": 500,
			"error":       err.Error(),
		})
		return
	}
	
	cid, err := cm.Create(uid, createCourse)
	if err != nil {
		tyrgin.ErrorHandler(err, c, 500, gin.H{
			"status_code": 500,
			"error":       err.Error(),
		})
		return
	}

	err = um.AddCourse("professor", *cid, uid)

	c.JSON(200, gin.H{
		"status_code": 200,
		"msg":         "Course created.",
	})
}