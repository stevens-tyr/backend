package cms

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"backend/errors"
	"backend/forms"
	"backend/forms/cmsforms"
	"backend/utils"
)

// versionCheck a function to check the language version(Docker) of a assignment creation form.
func versionCheck(a *forms.CreateAssignmentPreForm) {
	if a.Version == "" {
		a.Version = "latest"
	}
}

// CreateAssignment will create an assignment and add its id to a course.
func CreateAssignment(c *gin.Context) {
	var capre forms.CreateAssignmentPreForm
	err := c.ShouldBind(&capre)
	if err != nil {
		fmt.Println("fuck", err)
		c.Set("error", errors.ErrorInvlaidJSON)
		return
	}
	versionCheck(&capre)

	var tests []cmsforms.CreateAssginmentTest
	for _, test := range capre.Tests {
		var toAdd cmsforms.CreateAssginmentTest
		json.Unmarshal([]byte(test), &toAdd)
		tests = append(tests, toAdd)
	}
	
	capost := forms.CreateAssignmentPostForm{
		capre.Language,
		capre.Version,
		capre.Name,
		capre.NumAttempts,
		capre.Description,
		capre.DueDate,
		capre.TestBuildCMD,
		tests,
	}	

	cids, _ := c.Get("cids")
	aid, supportingFilesName, err := am.Create(capost, cids.(string))
	if err != nil {
		c.Set("error", err)
		return
	}

	cid, _ := c.Get("cid")
	err = cm.AddAssignment(*aid, cid)
	if err != nil {
		c.Set("error", err)
		return
	}

	sf, errs := c.FormFile("supportingFiles")
	if errs != nil && errs != http.ErrMissingFile {
		c.Set("error", errors.ErrorUploadingFile)
		return
	}

	supportingFiles, err := utils.CheckFileType(sf)
	if err != nil {
		c.Set("error", err)
		return
	}

	err = gfs.Upload(*aid, supportingFilesName, bytes.NewReader(supportingFiles))
	if err != nil {
		c.Set("error", err)
		return
	}

	c.JSON(200, gin.H{
		"message": "Assignment Created.",
	})
}
