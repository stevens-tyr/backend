package cms

import (
	ctx "context"
	"os"

	"backend/models"

	jwt "github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/objectid"
	"github.com/mongodb/mongo-go-driver/options"

	"github.com/stevens-tyr/tyr-gin"
)

// Dashboard is the function for a route to display all course a user has.
func Dashboard(c *gin.Context) {
	claims := jwt.ExtractClaims(c)

	db, err := tyrgin.GetMongoDB(os.Getenv("DB_NAME"))
	if !tyrgin.ErrorHandler(err, c, 500, gin.H{
		"status_code": 500,
		"message":     "Failed to get Mongo Session.",
		"error":       err,
	}) {
		return
	}

	col := tyrgin.GetMongoCollection("users", db)

	uid, err := objectid.FromHex(claims["uid"].(string))
	if !tyrgin.ErrorHandler(err, c, 500, gin.H{
		"status_code": 500,
		"message":     "Failed to extract userid as valid mongo objectid.",
		"error":       err,
	}) {
		return
	}

	query := []interface{}{
		bson.M{"$match": bson.M{"_id": uid}},
		bson.M{"$unwind": "$enrolledCourses"},
		bson.M{
			"$lookup": bson.M{
				"from":         "courses",
				"localField":   "enrolledCourses.courseID",
				"foreignField": "_id",
				"as":           "course",
			},
		},
		bson.M{"$project": bson.M{
			"_id": 0,
			"course": bson.M{
				"$arrayElemAt": bson.A{"$course", 0},
			},
		},
		},
	}

	cur, err := col.Aggregate(ctx.Background(), query, options.Aggregate())
	if !tyrgin.ErrorHandler(err, c, 500, gin.H{
		"status_code": 500,
		"message":     "Failed to query mongo.",
		"error":       err,
	}) {
		return
	}

	var courses []models.CourseAgg
	for cur.Next(ctx.Background()) {
		var course map[string]models.CourseAgg
		err = cur.Decode(&course)
		if !tyrgin.ErrorHandler(err, c, 500, gin.H{
			"status_code": 500,
			"message":     "Failed to decode course.",
			"error":       err,
		}) {
			return
		}
		courses = append(courses, course["course"])
	}

	c.JSON(200, gin.H{
		"status_code": 200,
		"msg":         "User's courses.",
		"courses":     courses,
	})
}
