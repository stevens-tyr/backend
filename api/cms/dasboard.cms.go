package cms

import (
	"fmt"
	"os"

	jwt "github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
	"github.com/stevens-tyr/tyr-gin"
	bson "gopkg.in/mgo.v2/bson"
)

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

	col, err := tyrgin.GetMongoCollectionCreate("users", db)
	if !tyrgin.ErrorHandler(err, c, 500, gin.H{
		"status_code": 500,
		"message":     "Failed to get collection.",
		"error":       err,
	}) {
		return
	}

	fmt.Println("uid", claims["uid"].(string))

	query := []bson.M{
		{"$match": bson.M{"_id": bson.ObjectIdHex(claims["uid"].(string))}},
		{"$unwind": "$enrolledCourses"},
		{
			"$lookup": bson.M{
				"from":         "courses",
				"localField":   "enrolledCourses.courseID",
				"foreignField": "_id",
				"as":           "course",
			},
		},
		//{"$project": bson.M{"course": bson.M{"$arrayElemAt": bson.D{{"$course", 0}}}, "_id": 0}},
	}
	fmt.Println(query)
	pipe := col.Pipe(query)

	var courses []map[string]interface{}
	err = pipe.All(&courses)
	fmt.Println(err, courses)

	c.JSON(200, gin.H{
		"status_code": 200,
		"msg":         "courses bitch",
		"courses":     courses[0]["course"],
	})
}
