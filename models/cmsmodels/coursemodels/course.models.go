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

// Course struct ot store information about a course.
type MongoCourse struct {
	ID          primitive.ObjectID   `bson:"_id,omitempty" json:"id" binding:"required"`
	Department  string               `bson:"department" json:"department" binding:"required"`
	Number      int                  `bson:"number" json:"number" binding:"required"`
	Section     string               `bson:"section" json:"section" binding:"required"`
	Semester    string               `bson:"semester" json:"semester" binding:"required"`
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

func (c *CourseInterface) FindOne(department, section, semester string, number int) (*MongoCourse, errors.APIError) {
	var course *MongoCourse

	res := c.col.FindOne(
		c.ctx,
		bson.M{
			"department": department,
			"number":     number,
			"section":    section,
			"semester":   semester,
		},
		options.FindOne(),
	)
	res.Decode(&course)

	if course == nil {
		return nil, errors.ErrorResourceNotFound
	}

	return course, nil
}

func (c *CourseInterface) Create(uid interface{}, form forms.CreateCourseForm) (*primitive.ObjectID, errors.APIError) {
	course, err := c.FindOne(
		form.Department,
		form.Section,
		form.Semester,
		form.Number,
	)
	if err != nil && err != errors.ErrorResourceNotFound {
		return nil, err
	}

	if course != nil {
		return nil, errors.ErrorCannotCreateDuplicateData
	}

	uids := uid.(string)
	val, errs := primitive.ObjectIDFromHex(uids)
	if errs != nil {
		return nil, errors.ErrorInvalidObjectID
	}
	professors := []primitive.ObjectID{val}

	course = &MongoCourse{
		Department:  form.Department,
		Number:      form.Number,
		Section:     form.Section,
		Semester:    form.Semester,
		Professors:  professors,
		Assistants:  make([]primitive.ObjectID, 0),
		Students:    make([]primitive.ObjectID, 0),
		Assignments: make([]primitive.ObjectID, 0),
	}

	res, errs := c.col.InsertOne(c.ctx, course, options.InsertOne())
	if errs != nil {
		return nil, errors.ErrorDatabaseFailedCreate
	}

	cid := res.InsertedID.(primitive.ObjectID)
	return &cid, nil
}

func (c *CourseInterface) UserExists(cid, uid interface{}) (bool, errors.APIError) {
	filter := bson.D{
		{"_id", cid},
		{
			"$or", bson.A{
				bson.M{"assitants": bson.M{"$elemMatch": uid}},
				bson.M{"professors": bson.M{"$elemMatch": uid}},
				bson.M{"students": bson.M{"$elemMatch": uid}},
			},
		},
	}

	res := c.col.FindOne(
		c.ctx,
		filter,
		options.FindOne(),
	)
	var course *MongoCourse
	err := res.Decode(&course)
	if err != nil {
		return false, errors.ErrorResourceNotFound
	}

	if course != nil {
		return true, nil
	}

	return false, nil
}

func (c *CourseInterface) AddUser(level string, uid, cid interface{}) errors.APIError {
	userAlreadyInCourse, err := c.UserExists(cid, uid)
	if err != nil {
		return err
	}

	if userAlreadyInCourse {
		return errors.ErrorUserAlreadyEnrolled
	}

	var tag string
	switch level {
	case "student":
		tag = "students"
		break
	case "assitant":
		tag = "assitants"
		break
	case "professor":
		tag = "professors"
		break
	}

	_, errs := c.col.UpdateOne(
		c.ctx,
		bson.M{"_id": cid},
		bson.M{"$push": bson.M{tag: uid}},
		options.Update(),
	)
	if errs != nil {
		return errors.ErrorDatabaseFailedUpdate
	}

	return nil
}

func (c *CourseInterface) AddAssignment(aid, cid interface{}) errors.APIError {
	_, err := c.col.UpdateOne(
		c.ctx,
		bson.M{"_id": cid},
		bson.M{"$push": bson.M{"assignments": aid}},
		options.Update(),
	)
	if err != nil {
		return errors.ErrorDatabaseFailedUpdate
	}

	return nil
}

func (c *CourseInterface) GetAssignments(cid interface{}) ([]forms.AssignmentAggQuery, errors.APIError) {
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
		return assignments, errors.ErrorDatabaseFailedQuery
	}

	for cur.Next(c.ctx) {
		var assignment map[string]forms.AssignmentAggQuery
		err = cur.Decode(&assignment)
		if err != nil {
			return assignments, errors.ErrorInvlaidBSON
		}
		if assignment["assignment"].Published {
			assignments = append(assignments, assignment["assignment"])
		}
	}

	return assignments, nil
}
