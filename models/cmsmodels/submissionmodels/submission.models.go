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
		Results string `bson:"results" json:"results" binding:"required"`
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
		AssignmentID        primitive.ObjectID `bson:"assignmentID" json:"assignmentID" binding:"required"`
		AttemptNumber int                `bson:"attemptNumber" json:"attemptNumber" binding:"required"`
		SubmissionDate primitive.DateTime `bson:"submissionDate" json:"submissionDate" binding:"required"`
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

// GetUsersRecentSubmissions grabs the most recent submissions up until limit
func (s *SubmissionInterface) GetUsersRecentSubmissions(uid interface{}, limit int64) ([]map[string]interface{}, errors.APIError) {
	query := []interface{}{
		bson.M{"$match": bson.M{"userID": uid}},
		bson.M{
			"$lookup": bson.M{
				"from": "courses",
				"let": bson.M{ "assID": "$assignmentID" },
				"pipeline": bson.A{
					bson.M{ "$match": bson.M{ "$expr": bson.M{ "$in": bson.A{"$$assID", "$assignments"} } } },
					bson.M{
						"$project": bson.M{
							"professors": 0,
							"assistants": 0,
							"students": 0,
							"assignments": 0,
						},
					},
				},
				"as": "course",
			},
		},
		bson.M{
			"$lookup": bson.M{
				"from": "assignments",
				"let": bson.M{ "assID": "$assignmentID" },
				"pipeline": bson.A{
					bson.M{
						"$match": bson.M{
							"$expr": bson.M{ "$eq": bson.A{"$_id", "$$assID"} } } },
					bson.M{
						"$project": bson.M{
							"tests": 0,
							"submissions": 0,
							"testBuildCMD": 0,
						},
					},
				},
				"as": "assignment",
			},
		},
		bson.M{
			"$project": bson.M{
				"course": bson.M{ "$arrayElemAt": bson.A{"$course", 0} },
				"assignmentID": 1,
				"submissionDate": 1,
				"file": 1,
				"errorTesting": 1,
				"cases.studentFacing": 1,
				"attemptNumber": 1,
				"assignment": bson.M{ "$arrayElemAt": bson.A{"$assignment", 0} },
			},
		},
		bson.M{ "$sort": bson.M{ "submissionDate": -1 } },
		bson.M{ "$limit": limit },
	}

	var recentSubmissions []map[string]interface{}
	cur, err := s.col.Aggregate(
		s.ctx,
		query,
		options.Aggregate(),
	)
	if err != nil {
		return nil, errors.ErrorInvlaidBSON
	}

	for cur.Next(s.ctx) {
		var submission map[string]interface{}
		err = cur.Decode(&submission)
		if submission == nil {
			return nil, errors.ErrorResourceNotFound
		}
		recentSubmissions = append(recentSubmissions, submission)
	}

	return recentSubmissions, nil
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
