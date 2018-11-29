package cms

import (
	ctx "context"
	"os"
	"fmt"

	// "backend/models"

	// jwt "github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
	"github.com/mongodb/mongo-go-driver/bson"

	"github.com/mongodb/mongo-go-driver/bson/objectid"
	"github.com/mongodb/mongo-go-driver/options"

	"github.com/stevens-tyr/tyr-gin"
)

type (
	// Course struct
	Course struct {
		ID          objectid.ObjectID `bson:"_id,omitempty" json:"id" binding:"required"`
		Department  string            `bson:"department" json:"department" binding:"required"`
		Number      int               `bson:"number" json:"number" binding:"required"`
		Section     string            `bson:"section" json:"section" binding:"required"`
		Name        string            `bson:"name" json:"name" binding:"required"`
		Assignments []objectid.ObjectID     `bson:"assignments" json:"assignments"`
	}
	//Assignment struct
	Assignment struct {
		ID 					objectid.ObjectID `bson:"_id,omitempty" json:"id" binding:"required"`
		Name        string            `bson:"name" json:"name" binding:"required"`
		Description string `bson:"description" json:"description" binding:"required"`
		DueDate 		string `bson:"dueDate" json:"dueDate" binding:"required"`
		Submissions []objectid.ObjectID     `bson:"assignments" json:"assignments"`
	}
	// //Submission struct
	// Assignment struct {
	// 	ID 					objectid.ObjectID `bson:"_id,omitempty" json:"id" binding:"required"`
	// 	Name        string            `bson:"name" json:"name" binding:"required"`
	// 	Description string `bson:"description" json:"description" binding:"required"`
	// 	DueDate 		string `bson:"dueDate" json:"dueDate" binding:"required"`
	// 	Submissions []objectid.ObjectID     `bson:"assignments" json:"assignments"`
	// }
)

func GetCourse(c *gin.Context) {

	db, err := tyrgin.GetMongoDB(os.Getenv("DB_NAME"))
	if !tyrgin.ErrorHandler(err, c, 500, gin.H{
		"status_code": 500,
		"message":     "Failed to get Mongo Session.",
		"error":       err,
	}) {
		return
	}

	col := tyrgin.GetMongoCollection("courses", db)

	objid, err := objectid.FromHex(c.Param("cid"))
	res := col.FindOne(ctx.Background(), bson.M{"_id": objid}, options.FindOne())
	
	var course Course
	err = res.Decode(&course)
	if err != nil {
		fmt.Println(err)
	}

	c.JSON(200, gin.H{
		"status_code": 200,
		"msg":         "Course",
		"course":      course,
	})
}

func GetAssignment(c *gin.Context) {

	db, err := tyrgin.GetMongoDB(os.Getenv("DB_NAME"))
	if !tyrgin.ErrorHandler(err, c, 500, gin.H{
		"status_code": 500,
		"message":     "Failed to get Mongo Session.",
		"error":       err,
	}) {
		return
	}

	col := tyrgin.GetMongoCollection("assignments", db)

	objid, err := objectid.FromHex(c.Param("aid"))
	res := col.FindOne(ctx.Background(), bson.M{"_id": objid}, options.FindOne())
	
	var assign Assignment
	err = res.Decode(&assign)
	if err != nil {
		fmt.Println(err)
	}

	c.JSON(200, gin.H{
		"status_code": 200,
		"msg":         "Assignment",
		"assign":      assign,
	})
}

// func GetAssignment(c *gin.Context) {

// 	db, err := tyrgin.GetMongoDB(os.Getenv("DB_NAME"))
// 	if !tyrgin.ErrorHandler(err, c, 500, gin.H{
// 		"status_code": 500,
// 		"message":     "Failed to get Mongo Session.",
// 		"error":       err,
// 	}) {
// 		return
// 	}

// 	col := tyrgin.GetMongoCollection("assignments", db)

// 	objid, err := objectid.FromHex(c.Param("aid"))
// 	res := col.FindOne(ctx.Background(), bson.M{"_id": objid}, options.FindOne())
	
// 	var assign Assignment
// 	err = res.Decode(&assign)
// 	if err != nil {
// 		fmt.Println(err)
// 	}

// 	c.JSON(200, gin.H{
// 		"status_code": 200,
// 		"msg":         "Assignment",
// 		"assign":      assign,
// 	})
// }
