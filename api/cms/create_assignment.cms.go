package cms

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
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
		c.Set("error", errors.ErrorInvalidJSON)
		return
	}
	versionCheck(&capre)

	var tests []cmsforms.CreateAssignmentTest
	for _, test := range capre.Tests {
		var toAdd cmsforms.CreateAssignmentTest
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
	aid, supportingFilesID, err := am.Create(capost, cids.(string))
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

	err = gfs.Upload(*aid, supportingFilesID, capre.Name, bytes.NewReader(supportingFiles))
	if err != nil {
		c.Set("error", err)
		am.Delete(*aid)
		return
	}

	c.JSON(200, gin.H{
		"message": "Assignment Created.",
	})
}

func CreateAssignmentFromFile(c *gin.Context) {
	afs, errs := c.FormFile("assignment")
	if errs != nil && errs != http.ErrMissingFile {
		c.Set("error", errors.ErrorUploadingFile)
		return
	}

	af, errs := afs.Open()
	if errs != nil {
		c.Set("error", errors.ErrorFailedToOpenFile)
		return
	}

	byteAF, errs := ioutil.ReadAll(af)
	if errs != nil {
		c.Set("error", errors.ErrorFailedToReadFile)
		return
	}

	var ca forms.CreateAssignmentPostForm
	json.Unmarshal(byteAF, &ca)

	cids, _ := c.Get("cids")
	aid, supportingFilesID, err := am.Create(ca, cids.(string))
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

	err = gfs.Upload(*aid, supportingFilesID, ca.Name, bytes.NewReader(supportingFiles))
	if err != nil {
		c.Set("error", err)
		return
	}

	c.JSON(200, gin.H{
		"message": "Assignment Created.",
	})
}
