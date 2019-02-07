package assignmentmodels

import (
	"context"
	"os"

	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/bson/primitive"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo/options"

	"github.com/stevens-tyr/tyr-gin"
)

type (
	// TestScripts struct to represent filenames on gcp storage of scripts.
	TestScripts struct {
		StudentFacing string `bson:"studentFacing" json:"studentFacing" binding:"required"`
		AdminFacing   string `bson:"adminFacing" json:"adminFacing" binding:"required"`
	}

	AssignmentSubmission struct {
		UserID        primitive.ObjectID `bson:"userID" json:"userID" binding:"required"`
		SubmissionID  primitive.ObjectID `bson:"submissionID" json:"submissionID" binding:"required"`
		AttemptNumber int                `bson:"attemptNumber" json:"attemptNumber" binding:"required"`
	}

	// Assignment struct to store information about an assignment.
	MongoAssignment struct {
		ID              primitive.ObjectID     `bson:"_id" json:"id" binding:"required"`
		Language        string                 `bson:"language" json:"lanaguage" binding:"required"`
		Version         string                 `bson:"version" json:"version" binding:"required"`
		Name            string                 `bson:"name" json:"name" binding:"required"`
		NumAttempts     int                    `bson:"numAttempts" json:"numAttempts" binding:"required"`
		Description     string                 `bson:"description" json:"description" binding:"required"`
		SupportingFiles string                 `bson:"supportingFiles" json:"supportingFiles" binding:"required"`
		DueDate         primitive.DateTime     `bson:"dueDate" json:"dueDate" binding:"required"`
		Published       bool                   `bson:"published" json:"published" binding:"required"`
		TestScripts     TestScripts            `bson:"testScripts" json:"testScripts" binding:"required"`
		Submissions     []AssignmentSubmission `bson:"submissions" json:"submissions" binding:"required"`
	}

	AssignmentInterface struct {
		ctx context.Context
		col *mongo.Collection
	}
)

func New() *AssignmentInterface {
	db, _ := tyrgin.GetMongoDB(os.Getenv("DB_NAME"))
	col := tyrgin.GetMongoCollection("assignments", db) 

	return &AssignmentInterface{
		context.Background(),
		col,
	}
}

func (a *AssignmentInterface) Create() (error) {
	return nil
}

func (a *AssignmentInterface) Get(aid primitive.ObjectID) (*MongoAssignment, error) {
	var assign *MongoAssignment
	res := a.col.FindOne(a.ctx, bson.M{"_id": aid}, options.FindOne())

	err := res.Decode(&assign)
	if err != nil {
		return nil, err
	}

	return assign, nil
}

func (a *AssignmentInterface) LatestUserSubmission(aid, uid primitive.ObjectID) (*MongoAssignment, int, error) {
	assignment, err := a.Get(aid)
	if err != nil {
		return nil, 0, err
	}

	attempt := 0
	for _, assignSub := range assignment.Submissions {
		if assignSub.UserID == uid && assignSub.AttemptNumber > attempt {
			attempt = assignSub.AttemptNumber
		}
	}

	return assignment, attempt, nil
}

func (a *AssignmentInterface) InsertSubmission(aid, sid, uid primitive.ObjectID, attempt int) (error) {
	insert := AssignmentSubmission{
		UserID:        uid,
		SubmissionID:  sid,
		AttemptNumber: attempt,
	}

	_, err := a.col.UpdateOne(
		a.ctx,
		bson.M{"_id": aid},
		bson.M{"$push": bson.M{"submissions": &insert}},
		options.Update(),
	)
	if err != nil {
		return err
	}

	return nil
}