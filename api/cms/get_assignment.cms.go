package cms

import (
	"github.com/gin-gonic/gin"
	"github.com/mongodb/mongo-go-driver/bson/primitive"
 
	"github.com/stevens-tyr/tyr-gin"
)

func GetAssignment(c *gin.Context) {
	aid, err := primitive.ObjectIDFromHex(c.Param("aid"))
	if err != nil {
		tyrgin.ErrorHandler(err, c, 500, gin.H{
			"status_code": 500,
			"error":       err.Error(),
		})
		return
	}

	assignment, err := am.Get(aid)
	if err != nil {
		tyrgin.ErrorHandler(err, c, 500, gin.H{
			"status_code": 500,
			"error":       err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"status_code": 200,
		"msg":         "assignment.",
		"assignment":   assignment,
	})
}