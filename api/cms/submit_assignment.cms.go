package cms

import (
	"bytes"
	"fmt"

	// "net/http"
	"os"
	"strconv"

	jwt "github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
	"github.com/mongodb/mongo-go-driver/bson/primitive"

	"backend/utils"

	"github.com/stevens-tyr/tyr-gin"
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

	fdb, err := tyrgin.GetMongoDB(os.Getenv("GRIDFS_DB_NAME"))
	if err != nil {
		c.Set("error", err)
		return
	}

	bucketSize, err := strconv.Atoi(os.Getenv("UPLOAD_SIZE"))
	if err != nil {
		c.Set("error", err)
		return
	}

	bucket, err := tyrgin.GetGridFSBucket(fdb, "fs", int32(bucketSize))
	if err != nil {
		c.Set("error", err)
		return
	}

	sid := primitive.NewObjectID()
	fid := primitive.NewObjectID()
	submittedFilesName := fmt.Sprintf("sub-%s-%s.tar.gz", c.Param("aid"), claims["uid"])
	reader := bytes.NewReader(submissionFiles)
	err = bucket.GridFSUploadFile(fid, submittedFilesName, reader)
	if err != nil {
		c.Set("error", err)
		return
	}

	uid, _ := c.Get("uid")
	aid, _ := c.Get("aid")

	// See if previous submission exists
	_, attempt, err := am.LatestUserSubmission(aid, uid)
	if err != nil {
		c.Set("error", err)
		return
	}

	err = am.InsertSubmission(aid, uid, sid, attempt+1)
	if err != nil {
		c.Set("error", err)
		return
	}

	job, err := sm.Submit(aid, fid, uid, sid, attempt+1, submittedFilesName)
	if err != nil {
		c.Set("error", err)
		return
	}

	c.JSON(201, gin.H{
		"status_code": 201,
		"message":     "Submission Grader Started.",
		"job":         job,
	})
}
