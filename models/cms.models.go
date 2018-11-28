package models

import (
	"github.com/mongodb/mongo-go-driver/bson/objectid"
	"github.com/mongodb/mongo-go-driver/bson/primitive"
)

// Assignment Creation/Submission types/structs
type (
	CreateAssignmentForm struct {
		Language    string             `form:"lang" binding:"required"`
		Version     string             `form:"version"`
		Name        string             `form:"name" json:"name" binding:"required"`
		NumAttempts int                `form:"numAttempts" binding:"required"`
		Description string             `form:"description" binding:"required"`
		DueDate     primitive.DateTime `form:"dueDate" binding:"required"`
	}

	// Course Aggregaton struct ot store information about a course.
	CourseAgg struct {
		ID         objectid.ObjectID `bson:"_id,omitempty" json:"id" binding:"required"`
		Department string            `bson:"department" json:"department" binding:"required"`
		Number     int               `bson:"number" json:"number" binding:"required"`
		Section    string            `bson:"section" json:"section" binding:"required"`
		Name       string            `bson:"name" json:"name" binding:"required"`
	}
)
