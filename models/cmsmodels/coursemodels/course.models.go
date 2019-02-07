package assignmentmodels

import (
 	"context"
 	"os"

	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/bson/primitive"
	// "github.com/mongodb/mongo-go-driver/bson"
	// "github.com/mongodb/mongo-go-driver/mongo/options"

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