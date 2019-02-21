package assignmentmodels

import (
	"context"
	"os"

	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/primitive"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/mongo/options"

	"backend/errors"
	"backend/forms"

	"github.com/stevens-tyr/tyr-gin"
)

type (
	AssignmentSubmission struct {
		UserID        primitive.ObjectID `bson:"userID" json:"userID" binding:"required"`
		SubmissionID  primitive.ObjectID `bson:"submissionID" json:"submissionID" binding:"required"`
		AttemptNumber int                `bson:"attemptNumber" json:"attemptNumber" binding:"required"`
	}

	Test struct {
		Name           string `bson:"name" json:"name" binding:"required"`
		ExpectedOutput string `bson:"expectedOutput" json:"expectedOutput" binding:"required"`
		StudentFacing  bool   `bson:"studentFacing" json:"studentFacing" binding:"exists"`
		TestCMD        string `bson:"testCMD" json:"testCMD" binding:"required"`
	}

	// Assignment struct to store information about an assignment.
	MongoAssignment struct {
		ID              primitive.ObjectID     `bson:"_id" json:"id"`
		Language        string                 `bson:"language" json:"lanaguage" binding:"required"`
		Version         string                 `bson:"version" json:"version" binding:"required"`
		Name            string                 `bson:"name" json:"name" binding:"required"`
		NumAttempts     int                    `bson:"numAttempts" json:"numAttempts" binding:"required"`
		Description     string                 `bson:"description" json:"description" binding:"required"`
		DueDate         primitive.DateTime     `bson:"dueDate" json:"dueDate" binding:"required"`
		Published       bool                   `bson:"published" json:"published" binding:"required"`
		SupportingFiles string                 `bson:"supportingFiles" json:"supportingFiles"`
		TestBuildCMD    string                 `bson:"TestBuildCMD" json:"testBuildCMD"`
		Tests           []Test                 `bson:"tests" json:"tests" binding:"required"`
		Submissions     []AssignmentSubmission `bson:"submissions" json:"submissions"`
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

func (a *AssignmentInterface) Create(form forms.CreateAssignmentForm) (*primitive.ObjectID, errors.APIError) {
	tests := make([]Test, len(form.Tests))
	for index := range form.Tests {
		tests[index] = Test(form.Tests[index])
	}

	aid := primitive.NewObjectID()
	assign := MongoAssignment{
		ID:          aid,
		Language:    form.Language,
		Version:     form.Version,
		Name:        form.Name,
		NumAttempts: form.NumAttempts,
		Description: form.Description,
		// SupportingFiles: fmt.Sprintf("%s.%s.supportingFiles.tar.gz", c.Param("cid"), aid),
		DueDate:      form.DueDate,
		Published:    false,
		TestBuildCMD: form.TestBuildCMD,
		Tests:        tests,
		Submissions:  make([]AssignmentSubmission, 0),
	}

	_, err := a.col.InsertOne(a.ctx, assign, options.InsertOne())
	if err != nil {
		return nil, errors.ErrorDatabaseFailedCreate
	}

	return &aid, nil
}

func (a *AssignmentInterface) Get(aid interface{}) (*MongoAssignment, errors.APIError) {
	var assign *MongoAssignment
	res := a.col.FindOne(a.ctx, bson.M{"_id": aid}, options.FindOne())

	err := res.Decode(&assign)
	if err != nil {
		return nil, errors.ErrorInvlaidBSON
	}

	return assign, nil
}

func (a *AssignmentInterface) LatestUserSubmission(aid, uid interface{}) (*MongoAssignment, int, errors.APIError) {
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

func (a *AssignmentInterface) InsertSubmission(aid, sid, uid interface{}, attempt int) errors.APIError {
	insert := AssignmentSubmission{
		UserID:        uid.(primitive.ObjectID),
		SubmissionID:  sid.(primitive.ObjectID),
		AttemptNumber: attempt,
	}

	_, err := a.col.UpdateOne(
		a.ctx,
		bson.M{"_id": aid},
		bson.M{"$push": bson.M{"submissions": &insert}},
		options.Update(),
	)
	if err != nil {
		return errors.ErrorDatabaseFailedUpdate
	}

	return nil
}
