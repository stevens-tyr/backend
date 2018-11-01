package models

import (
	bson "gopkg.in/mgo.v2/bson"
)

// Assignment Creation/Submission types/structs
type (
	CreateAssignment struct {
		Language    string              `form:"lang" binding:"required"`
		Version     string              `form:"version"`
		Name        string              `form:"name" json:"name" binding:"required"`
		NumAttempts int                 `form:"numAttempts" binding:"required"`
		Description string              `form:"description" binding:"required"`
		DueDate     bson.MongoTimestamp `form:"dueDate" binding:"required"`
	}
)
