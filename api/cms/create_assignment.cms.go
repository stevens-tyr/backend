package cms

import (
	"bytes"
	"net/http"

	"github.com/gin-gonic/gin"

	"backend/errors"
	"backend/forms"
	"backend/utils"
)

// versionCheck a function to check the language version(Docker) of a assignment creation form.
func versionCheck(a *forms.CreateAssignmentForm) {
	if a.Version == "" {
		a.Version = "latest"
	}
}

// CreateAssignment will create an assignment and add its id to a course.
func CreateAssignment(c *gin.Context) {
	var ca forms.CreateAssignmentForm
	err := c.ShouldBind(&ca)
	if err != nil {
		c.Set("error", errors.ErrorInvlaidJSON)
		return
	}
	versionCheck(&ca)

	cids, _ := c.Get("cids")
	aid, supportingFilesName, err := am.Create(ca, cids.(string))
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
