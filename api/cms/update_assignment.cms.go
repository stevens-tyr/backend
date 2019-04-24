package cms

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"backend/errors"
	"backend/forms"
	"backend/models/cmsmodels/assignmentmodels"
	"backend/utils"
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
		fmt.Println("ERROR:", errs)
		c.Set("error", errors.ErrorInvalidJSON)
		return
	}

	sf, errs := c.FormFile("supportingFiles")
	if errs != nil && errs != http.ErrMissingFile {
		c.Set("error", errors.ErrorUploadingFile)
		return
	}

	if sf != nil {
		supportingFiles, err := utils.CheckFileType(sf)
		if err != nil {
			c.Set("error", err)
			return
		}

		err = gfs.Delete(assign.SupportingFiles)
		if err != nil {
			c.Set("error", err)
			return
		}

		err = gfs.Upload(aid, &assign.SupportingFiles, assign.Name, bytes.NewReader(supportingFiles))
		if err != nil {
			c.Set("error", err)
			am.Delete(aid)
			return
		}
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
	if len(up.Tests) > 0 {
		var tests []assignmentmodels.Test
		for _, test := range up.Tests {
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
	if err != nil {
		c.Set("error", err)
		return
	}

	c.JSON(200, gin.H{
		"message": "Assignment Updated.",
	})
}
