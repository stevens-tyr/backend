package cms

import (
	"bytes"
	"fmt"
	"net/http"
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
		msg := "Problem uploading submission."
		if err == http.ErrMissingFile {
			msg = "Please upload your submission."
		}

		tyrgin.ErrorHandler(err, c, 400, gin.H{
			"status_code": 400,
			"message":     msg,
			"error":       err.Error() ,
		})
		return
	}

	submissionFiles, err := utils.CheckFileType(sub)
	if err != nil {
		tyrgin.ErrorHandler(err, c, 400, gin.H{
			"status_code":   400,
			"message":       "Incorrect file type for submission.",
			"allowed_types": []string{".zip", ".tar.gz"},
			"error":         err.Error(),
		}) 
		return
	}

	// Upload
	claims := jwt.ExtractClaims(c)

	fdb, err := tyrgin.GetMongoDB(os.Getenv("GRIDFS_DB_NAME"))
	if err != nil {
		tyrgin.ErrorHandler(err, c, 500, gin.H{
			"status_code": 500,
			"message":     "Failed to get Mongo Session/DB.",
			"error":       err.Error(),
		})
		return
	}

	bucketSize, err := strconv.Atoi(os.Getenv("UPLOAD_SIZE"))
	if err != nil {
		tyrgin.ErrorHandler(err, c, 500, gin.H{
			"staus_code": 500,
			"message":    "Failed to get gridfs bucket chunk size.",
			"error":      err.Error(),
		})
		return
	}

	bucket, err := tyrgin.GetGridFSBucket(fdb, "fs", int32(bucketSize))
	if err != nil {
		tyrgin.ErrorHandler(err, c, 500, gin.H{
			"staus_code": 500,
			"message":    "Failed to get files bucket.",
			"error":      err.Error(),
		})
		return
	}

	sid := primitive.NewObjectID()
	submittedFilesName := fmt.Sprintf("name%s%s%s%s.tar.gz", c.Param("cid"), c.Param("aid"), claims["uid"].(string), sid.String())
	reader := bytes.NewReader(submissionFiles)
	err = bucket.GridFSUploadFile(sid, submittedFilesName, reader)
	if err != nil {
		tyrgin.ErrorHandler(err, c, 500, gin.H{
			"staus_code": 500,
			"message":    "Failed to upload supporting files.",
			"error":      err.Error(),
		})
		return
	}

	uid, _ := primitive.ObjectIDFromHex(claims["uid"].(string))

	// See if previous submission exists
	//cid := c.Param("cid")
	aid, _ := primitive.ObjectIDFromHex(c.Param("aid"))

	_, attempt, err := am.LatestUserSubmission(aid, uid)
	if err != nil {
		tyrgin.ErrorHandler(err, c, 500, gin.H{
			"staus_code": 500,
			"message":    err.Error(),
		})
		return
	}

	err = am.InsertSubmission(aid, uid, sid, attempt + 1)
	if err != nil {
		tyrgin.ErrorHandler(err, c, 500, gin.H{
			"staus_code": 500,
			"message":    "Failed to update assignment.",
			"error":      err.Error(),
		})
		return
	}

	err = sm.Submit(aid, uid, sid, attempt + 1, submittedFilesName)
	if err != nil {
		tyrgin.ErrorHandler(err, c, 500, gin.H{
			"staus_code": 500,
			"message":    "Failed to create submission.",
			"error": err.Error(),
		})
		return
	}
	

	c.JSON(201, gin.H{
		"status_code": 201,
		"file_name":   submittedFilesName,
		"message":     "Submission Graded.",
	})
}
