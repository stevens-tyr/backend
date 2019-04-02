package cms

import (
	"github.com/mongodb/mongo-go-driver/bson/primitive"
	"github.com/gin-gonic/gin"

	"backend/errors"
	"backend/models/cmsmodels/assignmentmodels"
)

func UpdateAssignment(c *gin.Context) {
	aid, _ := c.Get("aid")

	assign, err := am.Get(aid)
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

	if val, ok := up["language"]; ok {
		assign.Language = val.(string)
	}
	if val, ok := up["version"]; ok {
		assign.Version = val.(string)
	}
	if val, ok := up["name"]; ok {
		assign.Name = val.(string)
	}
	if val, ok := up["description"]; ok {
		assign.Description = val.(string)
	}
	if val, ok := up["dueDate"]; ok {
		assign.DueDate = val.(primitive.DateTime)
	}
	if val, ok := up["published"]; ok {
		assign.Published = val.(bool)
	}
	if val, ok := up["testBuildCMD"]; ok {
		assign.TestBuildCMD = val.(string)
	}
	if val, ok := up["tests"]; ok {
		assign.Tests = val.([]assignmentmodels.Test)
	}
	if val, ok := up["numAttempts"]; ok {
		assign.NumAttempts = val.(int)
	}

	err = am.Update(*assign)
	if err != nil{
		c.Set("error", err)
		return
	}
	
	c.JSON(200, gin.H{
		"message": "Assignment Updated.",
	})
}
