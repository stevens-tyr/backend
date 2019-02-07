package assignmentmodels

import (
 	"context"
 	"os"

	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/bson/primitive"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo/options"

	"backend/forms"

	"github.com/stevens-tyr/tyr-gin"
)

// Course struct ot store information about a course.
type MongoCourse struct {
	ID          primitive.ObjectID   `bson:"_id,omitempty" json:"id" binding:"required"`
	Department  string               `bson:"department" json:"department" binding:"required"`
	Number      int                  `bson:"number" json:"number" binding:"required"`
	Section     string               `bson:"section" json:"section" binding:"required"`
	Professors  []primitive.ObjectID `bson:"professors" json:"professors" binding:"required"`
	Assistants  []primitive.ObjectID `bson:"assistants" json:"assitants" binding:"required"`
	Students    []primitive.ObjectID `bson:"students" json:"students" binding:"required"`
	Assignments []primitive.ObjectID `bson:"assignments" json:"assignments" binding:"required"`
}

type CourseInterface struct {
	ctx context.Context
	col *mongo.Collection
}

func New() *CourseInterface {
	db, _ := tyrgin.GetMongoDB(os.Getenv("DB_NAME"))
	col := tyrgin.GetMongoCollection("courses", db) 

	return &CourseInterface{
		context.Background(),
		col,
	}
}

func (c *CourseInterface) GetAssignments(cid primitive.ObjectID) ([]forms.AssignmentAggQuery, error) {
	var assignments []forms.AssignmentAggQuery

	query := []interface{}{
		bson.M{"$match": bson.M{"_id": cid}},
		bson.M{"$unwind": "$assignments"},
		bson.M{
			"$lookup": bson.M{
				"as":           "assignment",
				"from":         "assignments",
				"localField":   "assignments",
				"foreignField": "_id",
			},
		},
		bson.M{"$project": bson.M{
			"_id": 0,
			"assignment": bson.M{
				"$arrayElemAt": bson.A{"$assignment", 0},
			},
		},
		},
	}

	cur, err := c.col.Aggregate(c.ctx, query, options.Aggregate())
	if err != nil {
		return assignments, err
	}

	for cur.Next(c.ctx) {
		var assignment map[string]forms.AssignmentAggQuery
		err = cur.Decode(&assignment)
		if err != nil {
			return assignments, err
		}
		if assignment["assignment"].Published {
			 assignments = append(assignments, assignment["assignment"])
		}
	}

	return assignments, nil
}