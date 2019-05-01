package cms

import (
	"bytes"
	"fmt"

	jwt "github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
	"github.com/mongodb/mongo-go-driver/bson/primitive"

	"backend/errors"
	"backend/utils"
)

// SubmitAssignment will submit and grade the submission. Also updates the assignment.
func SubmitAssignment(c *gin.Context) {
	sub, err := c.FormFile("submission")
	if err != nil {
		c.Set("error", err)
		return
	}

	submissionFiles, err := utils.CheckFileType(sub)
	if err != nil {
		c.Set("error", err)
		return
	}

	// Upload
	claims := jwt.ExtractClaims(c)

	sid := primitive.NewObjectID()
	fid := primitive.NewObjectID()
	submittedFilesName := fmt.Sprintf("sub-%s-%s.tar.gz", c.Param("aid"), claims["uid"])
	reader := bytes.NewReader(submissionFiles)
	err = gfs.Upload(&fid, submittedFilesName, reader)
	if err != nil {
		c.Set("error", err)
		return
	}

	uid, _ := c.Get("uid")
	aid, _ := c.Get("aid")

	// See if previous submission exists
	assign, attempt, err := am.LatestUserSubmission(aid, uid)
	if err != nil {
		c.Set("error", err)
		return
	}

	if attempt+1 > assign.NumAttempts {
		c.Set("error", errors.ErrorSubmissionAttemptsExceeded)
		return
	}

	err = am.InsertSubmission(aid, uid, sid, attempt+1)
	if err != nil {
		c.Set("error", err)
		return
	}

	job, err := sm.Submit(aid, fid, uid, sid, attempt+1, submittedFilesName, assign.Tests, assign.TestBuildCMD, assign.Language)
	if err != nil {
		am.DeleteSubmission(aid, sid)
		c.Set("error", err)
		return
	}

	c.JSON(201, gin.H{
		"status_code": 201,
		"message":     "Submission Grader Started.",
		"job":         job,
	})
}
