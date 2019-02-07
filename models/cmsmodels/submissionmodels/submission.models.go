package submissionmodels

import (
	"context"
	"os"

	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/bson/primitive"
	// "github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo/options"

	"github.com/stevens-tyr/tyr-gin"
)

type (
	// FacingTests struct to store the number of tests for private and public facing tests.
	FacingTests struct {
		Pass int `bson:"pass" json:"pass" binding:"required"`
		Fail int `bson:"fail" json:"fail" binding:"required"`
	}

	// Cases struct to store tests/failed passed for admin/student
	Cases struct {
		StudentFacing FacingTests `bson:"studentFacing" json:"studentFacing" binding:"required"`
		AdminFacing   FacingTests `bson:"adminFacing" json:"adminFacing" binding:"required"`
	}

	// Submission struct the struct to represent a submission to an page.
	MongoSubmission struct {
		ID             primitive.ObjectID `bson:"_id" json:"id" binding:"required"`
		UserID         primitive.ObjectID `bson:"userID" json:"userID" binding:"required"`
		AttemptNumber  int                `bson:"attemptNumber" json:"attemptNumber" binding:"required"`
		SubmissionDate primitive.DateTime `bson:"submissionDate" json:"submissionDate" binding:"required"`
		File           string             `bson:"file" json:"file" binding:"required"`
		ErrorTesting   bool               `bson:"errorTesting" json:"errorTesting" binding:"required"`
		Cases          Cases              `bson:"cases" json:"cases" binding:"required"`
	}

	SubmissionInterface struct {
		ctx context.Context
		col *mongo.Collection
	}
)

func New() *SubmissionInterface {
	db, _ := tyrgin.GetMongoDB(os.Getenv("DB_NAME"))
	col := tyrgin.GetMongoCollection("submissions", db) 

	return &SubmissionInterface{
		context.Background(),
		col,
	}
}

func (s *SubmissionInterface) Submit(aid, uid, sid primitive.ObjectID, attempt int, filename string, ) (error) {
	submission := MongoSubmission{
		ID:            sid,
		UserID:        uid,
		AttemptNumber: attempt,
		File:          filename,
		ErrorTesting:  false,
		Cases: Cases{
			StudentFacing: FacingTests{
				Pass: 10,
				Fail: 0,
			},
			AdminFacing: FacingTests{
				Pass: 12,
				Fail: 3,
			},
		},
	}

	_, err := s.col.InsertOne(s.ctx, &submission, options.InsertOne())
	if err != nil {
		return err
	}
	
	return nil
}