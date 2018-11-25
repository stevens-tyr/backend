package cms

import (
	"bytes"
	ctx "context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"

	jwt "github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/objectid"
	"github.com/mongodb/mongo-go-driver/options"

	"backend/models"
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
			"error":       err,
		})
		return
	}

	submissionFiles, err := utils.CheckFileType(sub)
	if !tyrgin.ErrorHandler(err, c, 400, gin.H{
		"status_code":   400,
		"message":       "Incorrect file type for submission.",
		"allowed_types": []string{".zip", ".tar.gz"},
		"error":         err,
	}) {
		return
	}

	// Upload
	claims := jwt.ExtractClaims(c)

	fdb, err := tyrgin.GetMongoDB(os.Getenv("GRIDFS_DB_NAME"))
	if !tyrgin.ErrorHandler(err, c, 500, gin.H{
		"status_code": 500,
		"message":     "Failed to get Mongo Session/DB.",
		"error":       err,
	}) {
		return
	}

	bucketSize, err := strconv.Atoi(os.Getenv("UPLOAD_SIZE"))
	if !tyrgin.ErrorHandler(err, c, 500, gin.H{
		"staus_code": 500,
		"message":    "Failed to get gridfs bucket chunk size.",
		"error":      err,
	}) {
		return
	}

	bucket, err := tyrgin.GetGridFSBucket(fdb, fmt.Sprintf("%s%s", c.Param("cid"), c.Param("aid")), int32(bucketSize))
	if !tyrgin.ErrorHandler(err, c, 500, gin.H{
		"staus_code": 500,
		"message":    "Failed to get assignments bucket.",
		"error":      err,
	}) {
		return
	}

	sid := objectid.New()
	submittedFilesName := fmt.Sprintf("%s%s%s%s.tar.gz", c.Param("cid"), c.Param("aid"), claims["uid"].(string), sid)
	_, err = bucket.GridFSUploadFile(submittedFilesName, bytes.NewReader(submissionFiles))
	if !tyrgin.ErrorHandler(err, c, 500, gin.H{
		"staus_code": 500,
		"message":    "Failed to upload supporting files.",
		"error":      err,
	}) {
		return
	}

	// Run tests

	// Mongo
	db, err := tyrgin.GetMongoDB(os.Getenv("DB_NAME"))
	if !tyrgin.ErrorHandler(err, c, 500, gin.H{
		"status_code": 500,
		"message":     "Failed to get Mongo Session.",
		"error":       err,
	}) {
		return
	}

	subCol := tyrgin.GetMongoCollection("submissions", db)

	uid, _ := objectid.FromHex(claims["uid"].(string))

	// See if previous submission exists
	//cid := c.Param("cid")
	aid, _ := objectid.FromHex(c.Param("aid"))

	assignCol := tyrgin.GetMongoCollection("assignments", db)

	var assign models.Assignment
	res := assignCol.FindOne(ctx.Background(), bson.M{"_id": aid}, options.FindOne())

	err = res.Decode(&assign)
	if !tyrgin.ErrorHandler(err, c, 500, gin.H{
		"staus_code": 500,
		"message":    "Failed to find assignment.",
		"error":      err,
	}) {
		return
	}

	var previousSub models.AssignmentSubmission
	for _, assignSub := range assign.Submissions {
		if assignSub.UserID == uid && assignSub.AttemptNumber > previousSub.AttemptNumber {
			previousSub = assignSub
		}
	}

	if previousSub.AttemptNumber+1 > assign.NumAttempts && assign.NumAttempts != 0 {
		err = errors.New("Number of attemtps exceeded")
		tyrgin.ErrorHandler(err, c, 400, gin.H{
			"status_code": 400,
			"message":     "Number of attempts exceeded.",
			"error":       err,
		})
		return
	}

	// Otherwise (hardcoded values for now)
	msub := models.Submission{
		ID:            sid,
		UserID:        uid,
		AttemptNumber: previousSub.AttemptNumber + 1,
		File:          submittedFilesName,
		ErrorTesting:  true,
		Cases: models.Cases{
			StudentFacing: models.FacingTests{
				Pass: 10,
				Fail: 0,
			},
			AdminFacing: models.FacingTests{
				Pass: 12,
				Fail: 3,
			},
		},
	}

	up := models.AssignmentSubmission{
		UserID:        uid,
		SubmissionID:  msub.ID,
		AttemptNumber: previousSub.AttemptNumber + 1,
	}

	_, err = subCol.InsertOne(ctx.Background(), &msub, options.InsertOne())
	if err != nil {
		c.JSON(500, gin.H{
			"staus_code": 500,
			"message":    "Failed to create submission.",
		})
		return
	}

	_, err = assignCol.UpdateOne(
		ctx.Background(),
		bson.M{"_id": aid},
		bson.M{"$push": bson.M{"submissions": &up}},
		options.Update(),
	)
	if !tyrgin.ErrorHandler(err, c, 500, gin.H{
		"staus_code": 500,
		"message":    "Failed to update assignment.",
		"error":      err,
	}) {
		return
	}

	c.JSON(201, gin.H{
		"status_code": 201,
		"message":     "Submission Graded.",
	})
}
