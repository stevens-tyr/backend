package cms

import (

	jwt "github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
	"github.com/mongodb/mongo-go-driver/bson/primitive"
 
	"github.com/stevens-tyr/tyr-gin"
)

// Dashboard is the function for a route to display all course a user has.
func Dashboard(c *gin.Context) {
	claims := jwt.ExtractClaims(c)

	uid, err := primitive.ObjectIDFromHex(claims["uid"].(string))
	if err != nil {
		tyrgin.ErrorHandler(err, c, 500, gin.H{
			"status_code": 500,
			"error":       err.Error(),
		})
		return
	}

	courses, err := um.GetCourses(uid)
	if err != nil {
		tyrgin.ErrorHandler(err, c, 400, gin.H{
			"status_code": 400,
			"error": err.Error(),
		})
		return
	} 
	c.JSON(200, gin.H{
		"status_code": 200,
		"msg":         "User's courses.",
		"courses":     courses,
	})
}
