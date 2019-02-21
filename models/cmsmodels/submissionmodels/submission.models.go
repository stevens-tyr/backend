package submissionmodels

import (
	"context"
	"os"

	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/primitive"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/mongo/options"

	"backend/errors"

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
		ID            primitive.ObjectID `bson:"_id" json:"id" binding:"required"`
		UserID        primitive.ObjectID `bson:"userID" json:"userID" binding:"required"`
		AttemptNumber int                `bson:"attemptNumber" json:"attemptNumber" binding:"required"`
		// SubmissionDate primitive.DateTime `bson:"submissionDate" json:"submissionDate" binding:"required"`
		File         string `bson:"file" json:"file" binding:"required"`
		ErrorTesting bool   `bson:"errorTesting" json:"errorTesting" binding:"required"`
		Cases        Cases  `bson:"cases" json:"cases" binding:"required"`
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

func (s *SubmissionInterface) GetUsersSubmissions(uid interface{}) ([]MongoSubmission, errors.APIError) {
	var submissions []MongoSubmission
	cur, err := s.col.Find(
		s.ctx,
		bson.M{
			"userID": uid,
		},
		options.Find(),
	)

	for cur.Next(s.ctx) {
		var submission MongoSubmission
		err = cur.Decode(&submission)
		if err != nil {
			return submissions, errors.ErrorInvlaidBSON
		}

		submissions = append(submissions, submission)
	}

	return submissions, nil
}

func (s *SubmissionInterface) GetUsersSubmission(sid, uid interface{}) (*MongoSubmission, errors.APIError) {
	var submission *MongoSubmission
	res := s.col.FindOne(
		s.ctx,
		bson.M{
			"_id":    sid,
			"userID": uid,
		},
		options.FindOne(),
	)

	res.Decode(&submission)
	if submission == nil {
		return nil, errors.ErrorResourceNotFound
	}

	return submission, nil
}

func (s *SubmissionInterface) Submit(aid, uid, sid interface{}, attempt int, filename string) errors.APIError {
	submission := MongoSubmission{
		ID:            sid.(primitive.ObjectID),
		UserID:        uid.(primitive.ObjectID),
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
		return errors.ErrorDatabaseFailedCreate
	}

	return nil
}
