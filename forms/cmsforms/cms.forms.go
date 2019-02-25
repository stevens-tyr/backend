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
		Published bool               `bson:"published" json:"published" binding:"required"`
	}

	CourseAddUser struct {
		Level string `json:"level" binding:"required"`
		Email string `json:"email" binding:"required"`
	}

	test struct {
		Name           string `json:"name" binding:"required"`
		ExpectedOutput string `json:"expectedOutput" binding:"required"`
		StudentFacing  bool   `json:"studentFacing" binding:"exists"`
		TestCMD        string `json:"testCMD" binding:"required"`
	}

	CreateAssignment struct {
		Language     string             `form:"language" binding:"required"`
		Version      string             `form:"version"`
		Name         string             `form:"name" binding:"required"`
		NumAttempts  int                `form:"numAttempts" binding:"required"`
		Description  string             `form:"description" binding:"required"`
		DueDate      primitive.DateTime `form:"dueDate" binding:"required"`
		TestBuildCMD string             `form:"TestBuildCMD"`
		Tests        []test             `form:"tests" binding:"required"`
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
	}
)
