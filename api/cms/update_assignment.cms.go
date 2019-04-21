package cms

import (
	"encoding/json"
	
	"github.com/gin-gonic/gin"

	"backend/errors"
	"backend/models/cmsmodels/assignmentmodels"
	"backend/forms"
)

func UpdateAssignment(c *gin.Context) {
	aid, _ := c.Get("aid")

	assign, err := am.Get(aid)
	if err != nil {
		c.Set("error", err)
		return
	}

	var up forms.UpdateAssignmentForm
	errs := c.ShouldBind(&up)
	if errs != nil {
		c.Set("error", errors.ErrorInvalidJSON)
		return
	}

	if up.Language != nil {
		assign.Language = *up.Language
	}
	if up.Version != nil {
		assign.Version = *up.Version
	}
	if up.Name != nil {
		assign.Name = *up.Name
	}
	if up.Description != nil {
		assign.Description = *up.Description
	}
	if up.DueDate != nil {
		assign.DueDate = *up.DueDate
	}
	if up.Published != nil {
		assign.Published = *up.Published
	}
	if up.TestBuildCMD != nil {
		assign.TestBuildCMD = *up.TestBuildCMD
	}
	if up.Tests != nil {
		var tests []assignmentmodels.Test
		for _, test := range *up.Tests {
			var toAdd assignmentmodels.Test
			json.Unmarshal([]byte(test), &toAdd)
			tests = append(tests, toAdd)
		}
		assign.Tests = tests
	}
	if up.NumAttempts != nil {
		assign.NumAttempts = *up.NumAttempts
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
