package cmsforms

import (
	"github.com/mongodb/mongo-go-driver/bson/primitive"
)

// Assignment Creation/Submission types/structs
type (
	AssignmentAgg struct {
		ID        primitive.ObjectID `bson:"_id,omitempty" json:"id" binding:"required"`
		DueDate   primitive.DateTime `bson:"dueDate" json:"dueDate" binding:"required"`
		Name      string             `bson:"name" json:"name" binding:"required"`
		Published bool               `bson:"published" json:"-" binding:"required"`
		CourseID  primitive.ObjectID `bson:"courseID" json:"courseID" binding:",omitempty"`
	}

	CourseAddUser struct {
		Level string `json:"level" binding:"required"`
		Email string `json:"email" binding:"required"`
	}

	CourseBulkAddUser struct {
		Level  string   `json:"level" binding:"required"`
		Emails []string `json:"emails" binding:"required"`
	}

	CreateAssignmentTest struct {
		Name           string `json:"name"`
		ExpectedOutput string `json:"expectedOutput"`
		StudentFacing  bool   `json:"studentFacing"`
		TestCMD        string `json:"testCMD"`
	}

	CreateAssignmentPreParse struct {
		Language     string             `form:"language" binding:"required"`
		Version      string             `form:"version"`
		Name         string             `form:"name" binding:"required"`
		NumAttempts  int                `form:"numAttempts" binding:"required"`
		Description  string             `form:"description" binding:"required"`
		DueDate      primitive.DateTime `form:"dueDate" binding:"required"`
		TestBuildCMD string             `form:"TestBuildCMD"`
		Tests        []string           `form:"tests" binding:"required"`
	}

	CreateAssignmentPostParse struct {
		Language     string
		Version      string
		Name         string
		NumAttempts  int
		Description  string
		DueDate      primitive.DateTime
		TestBuildCMD string
		Tests        []CreateAssignmentTest
	}

	CreateCourse struct {
		Department string `json:"department" binding:"required"`
		Number     int    `json:"number" binding:"required"`
		Section    string `json:"section" binding:"required"`
		Semester   string `json:"semester" binding:"required"`
	}

	// Course Aggregaton struct ot store information about a course.
	CourseAgg struct {
		ID         primitive.ObjectID `bson:"_id,omitempty" json:"id" binding:"required"`
		Department string             `bson:"department" json:"department" binding:"required"`
		Number     int                `bson:"number" json:"number" binding:"required"`
		Section    string             `bson:"section" json:"section" binding:"required"`
		LongName   string             `bson:"longName" json:"longName" binding:"required"`
		Role       string             `json:"role" binding:"required"`
	}

	sub struct {
		Time    primitive.DateTime `form:"submissionTime""`
		Attempt int                `bson:"attemptNumber"`
	}

	student struct {
		ID    primitive.ObjectID `bson:"_id,omitempty"`
		First string             `bson:"firstName"`
		Last  string             `bson:"lastName"`
		Subs  []sub              `bson:"submissions"`
	}

	GradeAgg struct {
		Students []student `bson:"students"`
	}
)
