package cms

import (
	"github.com/gin-gonic/gin"
	"github.com/mongodb/mongo-go-driver/bson/primitive"

	"backend/forms"

	"github.com/stevens-tyr/tyr-gin"
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
		tyrgin.ErrorHandler(err, c, 400, gin.H{
			"status_code": 400,
			"message":     "Invalid json.",
			"error":       err.Error(),
		})
		return
	}
	versionCheck(&ca)

	aid, err := am.Create(ca)
	if err != nil {
		tyrgin.ErrorHandler(err, c, 400, gin.H{
			"status_code": 400,
			"error":       err.Error(),
		})
		return
	} 

	cid, err := primitive.ObjectIDFromHex(c.Param("cid"))
	if err != nil {
		tyrgin.ErrorHandler(err, c, 500, gin.H{
			"status_code": 500,
			"error":       err.Error(),
		})
		return
	}

	err = cm.AddAssignment(*aid, cid)
	if err != nil {
		tyrgin.ErrorHandler(err, c, 400, gin.H{
			"status_code": 400,
			"error":       err.Error(),
		})
		return
	} 
	// sft, err := c.FormFile("studentFacingTests")
	// if err != nil && err != http.ErrMissingFile {
	// 	tyrgin.ErrorHandler(err, c, 400, gin.H{
	// 		"status_code": 400,
	// 		"message":     "Problem uploading student facing tests.",
	// 		"error":       err,
	// 	})
	// 	return
	// }

	// studentTestFiles, err := utils.CheckFileType(sft)
	// if err != nil && err != utils.ErrorFileDNE {
	// 	tyrgin.ErrorHandler(err, c, 400, gin.H{
	// 		"status_code":   400,
	// 		"message":       "Incorrect file type for student facing tests.",
	// 		"allowed_types": []string{".zip", ".tar.gz"},
	// 		"error":         err,
	// 	})
	// 	return
	// }

	// aft, err := c.FormFile("adminFacingTests")
	// if err != nil {
	// 	msg := "Problem uploading admin facing tests."
	// 	if err == http.ErrMissingFile {
	// 		msg = "Please upload an admin facing test."
	// 	}

	// 	tyrgin.ErrorHandler(err, c, 400, gin.H{
	// 		"status_code": 400,
	// 		"message":     msg,
	// 		"error":       err,
	// 	})
	// 	return
	// }

	// adminTestFiles, err := utils.CheckFileType(aft)
	// if !tyrgin.ErrorHandler(err, c, 400, gin.H{
	// 	"status_code":   400,
	// 	"message":       "Incorrect file type for student facing tests.",
	// 	"allowed_types": []string{".zip", ".tar.gz"},
	// 	"error":         err,
	// }) {
	// 	return
	// }

	// sf, err := c.FormFile("supportingFiles")
	// if err != nil {
	// 	msg := "Problem uploading Supporting Files for assignmnet."
	// 	if err == http.ErrMissingFile {
	// 		msg = "Please upload supporting files."
	// 	}

	// 	tyrgin.ErrorHandler(err, c, 400, gin.H{
	// 		"status_code": 400,
	// 		"message":     msg,
	// 		"error":       err,
	// 	})
	// 	return
	// }

	// supportingFiles, err := utils.CheckFileType(sf)
	// if !tyrgin.ErrorHandler(err, c, 400, gin.H{
	// 	"status_code":   400,
	// 	"message":       "Incorrect file type for student facing tests.",
	// 	"allowed_types": []string{".zip", ".tar.gz"},
	// 	"error":         err,
	// }) {
	// 	return
	// }

	// // mongo
	// db, err := tyrgin.GetMongoDB(os.Getenv("DB_NAME"))
	// if !tyrgin.ErrorHandler(err, c, 500, gin.H{
	// 	"status_code": 500,
	// 	"message":     "Failed to get Mongo Session/DB.",
	// 	"error":       err,
	// }) {
	// 	return
	// }

	// assignCol := tyrgin.GetMongoCollection("assignments", db)

	// aid := primitive.NewObjectID()
	// assign := models.Assignment{
	// 	ID:              aid,
	// 	Language:        ca.Language,
	// 	Version:         ca.Version,
	// 	Name:            ca.Name,
	// 	NumAttempts:     ca.NumAttempts,
	// 	Description:     ca.Description,
	// 	SupportingFiles: fmt.Sprintf("%s.%s.supportingFiles.tar.gz", c.Param("cid"), aid),
	// 	DueDate:         ca.DueDate,
	// 	Published:       false,
	// 	TestScripts: models.TestScripts{
	// 		StudentFacing: "",
	// 		AdminFacing:   fmt.Sprintf("%s.%s.adminTestFiles.tar.gz", c.Param("cid"), aid),
	// 	},
	// 	Submissions: make([]models.AssignmentSubmission, 0),
	// }

	// if len(studentTestFiles) > 0 {
	// 	assign.TestScripts.StudentFacing = fmt.Sprintf("%s.%s.studentTestFiles.tar.gz", c.Param("cid"), assign.ID)
	// }

	// _, err = assignCol.InsertOne(ctx.Background(), &assign, options.InsertOne())
	// if !tyrgin.ErrorHandler(err, c, 500, gin.H{
	// 	"status_code": 500,
	// 	"message":     "Failed to create assignment.",
	// 	"error":       err,
	// }) {
	// 	return
	// }

	// courseCol := tyrgin.GetMongoCollection("courses", db)

	// cid, err := primitive.ObjectIDFromHex(c.Param("cid"))
	// if !tyrgin.ErrorHandler(err, c, 500, gin.H{
	// 	"staus_code": 500,
	// 	"message":    "Inavlid course id of course.",
	// 	"error":      err,
	// }) {
	// 	return
	// }

	// sec := c.Param("section")

	// _, err = courseCol.UpdateOne(
	// 	ctx.Background(),
	// 	bson.M{"_id": cid, "section": sec},
	// 	bson.M{"$push": bson.M{"assignments": assign.ID}},
	// 	options.Update(),
	// )
	// if !tyrgin.ErrorHandler(err, c, 500, gin.H{
	// 	"staus_code": 500,
	// 	"message":    "Failed to update course.",
	// 	"error":      err,
	// }) {
	// 	return
	// }

	// // upload files
	// fdb, err := tyrgin.GetMongoDB(os.Getenv("GRIDFS_DB_NAME"))
	// if !tyrgin.ErrorHandler(err, c, 500, gin.H{
	// 	"status_code": 500,
	// 	"message":     "Failed to get Mongo Session/DB.",
	// 	"error":       err,
	// }) {
	// 	return
	// }

	// bucketSize, err := strconv.Atoi(os.Getenv("UPLOAD_SIZE"))
	// if !tyrgin.ErrorHandler(err, c, 500, gin.H{
	// 	"staus_code": 500,
	// 	"message":    "Failed to get gridfs bucket chunk size.",
	// 	"error":      err,
	// }) {
	// 	return
	// }

	// bucket, err := tyrgin.GetGridFSBucket(fdb, fmt.Sprintf("%s%s", c.Param("cid"), aid), int32(bucketSize))
	// if !tyrgin.ErrorHandler(err, c, 500, gin.H{
	// 	"staus_code": 500,
	// 	"message":    "Failed to get assignments bucket.",
	// 	"error":      err,
	// }) {
	// 	return
	// }

	// err = bucket.GridFSUploadFile(primitive.NewObjectID(), assign.TestScripts.AdminFacing, bytes.NewReader(adminTestFiles))
	// if !tyrgin.ErrorHandler(err, c, 500, gin.H{
	// 	"staus_code": 500,
	// 	"message":    "Failed to upload admin facing tests.",
	// 	"error":      err,
	// }) {
	// 	return
	// }

	// err = bucket.GridFSUploadFile(primitive.NewObjectID(), assign.SupportingFiles, bytes.NewReader(supportingFiles))
	// if !tyrgin.ErrorHandler(err, c, 500, gin.H{
	// 	"staus_code": 500,
	// 	"message":    "Failed to upload supporting files.",
	// 	"error":      err,
	// }) {
	// 	return
	// }

	// if len(studentTestFiles) > 0 {
	// 	err = bucket.GridFSUploadFile(primitive.NewObjectID(), assign.TestScripts.StudentFacing, bytes.NewReader(studentTestFiles))
	// 	if !tyrgin.ErrorHandler(err, c, 500, gin.H{
	// 		"staus_code": 500,
	// 		"message":    "Failed to upload student facing tests.",
	// 		"error":      err,
	// 	}) {
	// 		return
	// 	}
	// }

	c.JSON(201, gin.H{
		"status_code":    201,
		"message":        "Assignment Created.",
		"assignmentLink": "url",
	})
}
