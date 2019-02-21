package cms

import (
	"github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"

	"backend/api/auth"
	"backend/errors"
	"backend/forms"
)

func CreateCourse(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	uid := claims["uid"]

	var createCourse forms.CreateCourseForm
	if errs := c.ShouldBindJSON(&createCourse); errs != nil {
		c.Set("error", errors.ErrorInvlaidJSON)
		return
	}

	cid, err := cm.Create(uid, createCourse)
	if err != nil {
		c.Set("error", err)
		return
	}

	err = um.AddCourse("professor", *cid, uid)
	if err != nil {
		c.Set("error", err)
		return
	}

	token, expire, errs := auth.AuthMiddleware.RefreshToken(c)
	if errs != nil {
		c.Set("error", errors.ErrorGenerateTokenFailure)
		return
	}

	c.JSON(200, gin.H{
		"status_code": 200,
		"msg":         "Course created.",
		"token":       token,
		"expire":      expire,
	})
}
