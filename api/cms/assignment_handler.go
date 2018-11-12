package cms

import (
	"errors"
	"net/http"
	"os"

	jwt "github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
	"github.com/stevens-tyr/tyr-gin"
	bson "gopkg.in/mgo.v2/bson"

	"backend/models"
	"backend/utils"
)

// versionCheck a function to check the language version(Docker) of a assignment creation form.
func versionCheck(a *models.CreateAssignment) {
	if a.Version == "" {
		a.Version = "latest"
	}
}

// CreateAssignment will create an assignment and add its id to a course.
func CreateAssignment(c *gin.Context) {
	var ca models.CreateAssignment
	err := c.ShouldBind(&ca)
	if !tyrgin.ErrorHandler(err, c, 400, gin.H{
		"status_code": 400,
		"message":     "Invalid json.",
		"error":       err,
	}) {
		return
	}
	versionCheck(&ca)

	sft, err := c.FormFile("studentFacingTests")
	if err != nil && err != http.ErrMissingFile {
		tyrgin.ErrorHandler(err, c, 400, gin.H{
			"status_code": 400,
			"message":     "Problem uploading student facing tests.",
			"error":       err,
		})
		return
	}

	_, err = utils.CheckFileType(sft)
	if err != nil && err != utils.ErrorFileDNE {
		tyrgin.ErrorHandler(err, c, 400, gin.H{
			"status_code":   400,
			"message":       "Incorrect file type for student facing tests.",
			"allowed_types": []string{".zip", ".tar.gz"},
			"error":         err,
		})
		return
	}

	aft, err := c.FormFile("adminFacingTests")
	if err != nil {
		msg := "Problem uploading student facing tests."
		if err == http.ErrMissingFile {
			msg = "Please upload an admin facing test."
		}

		tyrgin.ErrorHandler(err, c, 400, gin.H{
			"status_code": 400,
			"message":     msg,
			"error":       err,
		})
		return
	}

	_, err = utils.CheckFileType(aft)
	if !tyrgin.ErrorHandler(err, c, 400, gin.H{
		"status_code":   400,
		"message":       "Incorrect file type for student facing tests.",
		"allowed_types": []string{".zip", ".tar.gz"},
		"error":         err,
	}) {
		return
	}

	// mongo
	db, err := tyrgin.GetMongoDB(os.Getenv("DB_NAME"))
	if !tyrgin.ErrorHandler(err, c, 500, gin.H{
		"status_code": 500,
		"message":     "Failed to get Mongo Session.",
		"error":       err,
	}) {
		return
	}

	assignCol, err := tyrgin.GetMongoCollectionCreate("assignments", db)
	if !tyrgin.ErrorHandler(err, c, 500, gin.H{
		"status_code": 500,
		"message":     "Failed to get collection.",
	}) {
		return
	}

	assign := models.Assignment{
		ID:              bson.NewObjectId(),
		Language:        ca.Language,
		Version:         ca.Version,
		Name:            ca.Name,
		NumAttempts:     ca.NumAttempts,
		Description:     ca.Description,
		SupportingFiles: "url",
		DueDate:         ca.DueDate,
		Published:       false,
		TestScripts: models.TestScripts{
			StudentFacing: "url",
			AdminFacing:   "url",
		},
		Submissions: make([]models.AssignmentSubmission, 0),
	}

	err = assignCol.Insert(&assign)
	if !tyrgin.ErrorHandler(err, c, 500, gin.H{
		"status_code": 500,
		"message":     "Failed to create assignment.",
		"error":       err,
	}) {
		return
	}

	courseCol, err := tyrgin.GetMongoCollectionCreate("courses", db)
	if !tyrgin.ErrorHandler(err, c, 500, gin.H{
		"status_code": 500,
		"message":     "Failed to get collection.",
		"error":       err,
	}) {
		return
	}

	cid := bson.ObjectIdHex(c.Param("cid"))
	sec := c.Param("section")

	err = courseCol.Update(bson.M{"_id": cid, "section": sec}, bson.M{"$push": bson.M{"assignments": assign.ID}})
	if !tyrgin.ErrorHandler(err, c, 500, gin.H{
		"staus_code": 500,
		"message":    "Failed to update course.",
		"error":      err,
	}) {
		return
	}

	// upload files to google cloud

	c.JSON(201, gin.H{
		"status_code":              201,
		"message":                  "Assignment Created.",
		"assignmentSubmissionLink": "...",
	})
}

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

	_, err = utils.CheckFileType(sub)
	if !tyrgin.ErrorHandler(err, c, 400, gin.H{
		"status_code":   400,
		"message":       "Incorrect file type for submission.",
		"allowed_types": []string{".zip", ".tar.gz"},
		"error":         err,
	}) {
		return
	}

	// Upload

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

	subCol, err := tyrgin.GetMongoCollectionCreate("submissions", db)
	if !tyrgin.ErrorHandler(err, c, 500, gin.H{
		"status_code": 500,
		"message":     "Failed to get sub collection.",
		"error":       err,
	}) {
		return
	}

	claims := jwt.ExtractClaims(c)
	uid := bson.ObjectIdHex(claims["uid"].(string))
	//uid := bson.ObjectIdHex("5bd7a96091895e864db1ab7b")

	// See if previous submission exists
	//cid := c.Param("cid")
	aid := bson.ObjectIdHex(c.Param("aid"))

	assignCol, err := tyrgin.GetMongoCollectionCreate("assignments", db)
	if !tyrgin.ErrorHandler(err, c, 500, gin.H{
		"status_code": 500,
		"message":     "Failed to get assign collection.",
		"error":       err,
	}) {
		return
	}

	var assign models.Assignment
	err = assignCol.Find(bson.M{"_id": aid}).One(&assign)
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
		ID:            bson.NewObjectId(),
		UserID:        uid,
		AttemptNumber: previousSub.AttemptNumber + 1,
		File:          "url",
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
		AttemptNumber: 1,
	}

	if err = subCol.Insert(&msub); err != nil {
		c.JSON(500, gin.H{
			"staus_code": 500,
			"message":    "Failed to create submission.",
		})
		return
	}

	err = assignCol.Update(bson.M{"_id": aid}, bson.M{"$push": bson.M{"submissions": &up}})
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
