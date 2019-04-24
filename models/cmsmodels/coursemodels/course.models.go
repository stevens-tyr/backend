package coursemodels

import (
	"bytes"
	"context"
	"encoding/csv"
	"os"
	"strconv"

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
	LongName    string               `bson:"longName" json:"longName" binding:"required"`
	Number      int                  `bson:"number" json:"number" binding:"required"`
	Section     string               `bson:"section" json:"section" binding:"required"`
	Semester    string               `bson:"semester" json:"semester" binding:"required"`
	Professors  []primitive.ObjectID `bson:"professors" json:"professors" binding:"required"`
	Assistants  []primitive.ObjectID `bson:"assistants" json:"assistants" binding:"required"`
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

func (c *CourseInterface) GetByID(cid interface{}) (*MongoCourse, errors.APIError) {
	var course *MongoCourse

	res := c.col.FindOne(c.ctx, bson.M{"_id": cid}, options.FindOne())
	res.Decode(&course)

	if course == nil {
		return nil, errors.ErrorResourceNotFound
	}

	return course, nil
}

func (c *CourseInterface) Delete(cid interface{}) errors.APIError {
	_, err := c.col.DeleteOne(c.ctx, bson.M{"_id": cid}, options.Delete())
	if err != nil {
		return errors.ErrorDatabaseFailedDelete
	}

	return nil
}

func (c *CourseInterface) RemoveAssignment(aid, cid interface{}) errors.APIError {
	_, err := c.col.UpdateOne(
		c.ctx,
		bson.M{
			"_id": cid,
		},
		bson.M{
			"$pull": bson.M{
				"assignments": aid,
			},
		},
	)
	if err != nil {
		return errors.ErrorDatabaseFailedUpdate
	}

	return nil
}

func (c *CourseInterface) Update(course MongoCourse) errors.APIError {
	_, err := c.col.UpdateOne(
		c.ctx,
		bson.M{
			"_id": course.ID,
		},
		bson.M{
			"$set": bson.M{
				"department": course.Department,
				"longName":   course.LongName,
				"section":    course.Section,
				"semester":   course.Semester,
			},
		},
	)
	if err != nil {
		return errors.ErrorDatabaseFailedUpdate
	}

	return nil
}

func (c *CourseInterface) Get(cid, uid interface{}, role string) (map[string]interface{}, errors.APIError) {
	userLookup := func(userType string) bson.M {
		return bson.M{
			"$lookup": bson.M{
				"from": "users",
				"let":  bson.M{"userType": "$" + userType},
				"pipeline": bson.A{
					bson.M{"$match": bson.M{"$expr": bson.M{"$in": bson.A{"$_id", "$$userType"}}}},
					bson.M{"$project": bson.M{"admin": 0, "enrolledCourses": 0, "password": 0}},
				},
				"as": userType,
			},
		}
	}

	query := []interface{}{
		bson.M{"$match": bson.M{"_id": cid}},
		userLookup("students"),
		userLookup("professors"),
		userLookup("assistants"),
	}

	if role == "student" {
		query = append(query, bson.M{
			"$lookup": bson.M{
				"from": "assignments",
				"let":  bson.M{"ass": "$assignments"},
				"pipeline": bson.A{
					bson.M{
						"$match": bson.M{
							"$expr": bson.M{"$and": bson.A{bson.M{"$in": bson.A{"$_id", "$$ass"}}, "$published"}},
						},
					},
					bson.M{
						"$lookup": bson.M{
							"from": "submissions",
							"let":  bson.M{"assID": "$_id"},
							"pipeline": bson.A{
								bson.M{
									"$match": bson.M{
										"$expr": bson.M{
											"$and": bson.A{
												bson.M{"$eq": bson.A{"$userID", uid}},
												bson.M{"$eq": bson.A{"$$assID", "$assignmentID"}},
											},
										},
									},
								},
								bson.M{
									"$project": bson.M{
										"assignmentID":   1,
										"submissionDate": 1,
										"file":           1,
										"errorTesting":   1,
										"results":        bson.M{"$filter": bson.M{"input": "$results", "as": "result", "cond": bson.M{"$eq": bson.A{"$$result.studentFacing", true}}}},
										"attemptNumber":  1,
										"inProgress":     1,
									},
								},
								bson.M{"$sort": bson.M{"submissionDate": -1}},
							},
							"as": "submissions",
						},
					},
				},
				"as": "assignments",
			},
		})
	} else {
		query = append(query, bson.M{
			"$lookup": bson.M{
				"from": "assignments",
				"let":  bson.M{"ass": "$assignments"},
				"pipeline": bson.A{
					bson.M{"$match": bson.M{"$expr": bson.M{"$in": bson.A{"$_id", "$$ass"}}}},
					bson.M{
						"$lookup": bson.M{
							"from":         "submissions",
							"localField":   "submissions.submissionID",
							"foreignField": "_id",
							"as":           "submissions",
						},
					},
				},
				"as": "assignments",
			},
		})
	}

	var course map[string]interface{}
	cur, err := c.col.Aggregate(
		c.ctx,
		query,
		options.Aggregate(),
	)
	if err != nil {
		return nil, errors.ErrorInvalidBSON
	}

	for cur.Next(c.ctx) {
		err = cur.Decode(&course)
		if course == nil {
			return nil, errors.ErrorResourceNotFound
		}
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

	uidpo := uid.(primitive.ObjectID)
	professors := []primitive.ObjectID{uidpo}

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
				bson.M{"assitants": uid},
				bson.M{"professors": uid},
				bson.M{"students": uid},
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
	case "assistant":
		tag = "assistants"
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

func (c *CourseInterface) GetAssignments(cid interface{}, role string) ([]forms.AssignmentAggQuery, errors.APIError) {
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
			return assignments, errors.ErrorInvalidBSON
		}
		if role != "student" || assignment["assignment"].Published {
			assignments = append(assignments, assignment["assignment"])
		}
	}

	return assignments, nil
}

func (c *CourseInterface) GetGradesAsCSV(aid, cid interface{}) (*bytes.Buffer, string, int64, errors.APIError) {
	userLookup := func(userType string) bson.M {
		return bson.M{
			"$lookup": bson.M{
				"from": "users",
				"let":  bson.M{"userType": "$" + userType},
				"pipeline": bson.A{
					bson.M{"$match": bson.M{"$expr": bson.M{"$in": bson.A{"$_id", "$$userType"}}}},
					bson.M{"$project": bson.M{"admin": 0, "email": 0, "enrolledCourses": 0, "password": 0}},
					bson.M{
						"$lookup": bson.M{
							"from": "submissions",
							"let":  bson.M{"uid": "$_id"},
							"pipeline": bson.A{
								bson.M{
									"$match": bson.M{
										"$expr": bson.M{
											"$and": bson.A{
												bson.M{"$eq": bson.A{"$assignmentID", aid}},
												bson.M{"$eq": bson.A{"$userID", "$$uid"}},
											},
										},
									},
								},
								bson.M{"$sort": bson.M{"submissionDate": -1}},
								bson.M{"$project": bson.M{"_id": 0, "assignmentID": 0, "userID": 0, "file": 0}},
								bson.M{"$limit": 1},
							},
							"as": "submissions",
						},
					},
				},
				"as": userType,
			},
		}
	}

	query := []interface{}{
		bson.M{"$match": bson.M{"_id": cid}},
		userLookup("students"),
		bson.M{
			"$project": bson.M{
				"_id":         0,
				"assignments": 0,
				"assistants":  0,
				"department":  0,
				"longName":    0,
				"number":      0,
				"professors":  0,
				"section":     0,
				"semester":    0,
			},
		},
	}

	var results forms.GradeAggQuery
	cur, err := c.col.Aggregate(
		c.ctx,
		query,
		options.Aggregate(),
	)
	if err != nil {
		return nil, "", 0, errors.ErrorInvalidBSON
	}

	for cur.Next(c.ctx) {
		err = cur.Decode(&results)
		if err != nil {
			return nil, "", 0, errors.ErrorResourceNotFound
		}
	}

	records := [][]string{
		{"First Name", "Last Name", "Grade", "TestCases", "Attempt Number", "Submission Time"},
	}

	for _, student := range results.Students {
		var grade, attempt string
		if len(student.Subs) > 0 {
			sub := student.Subs[0]
			grade = "100"
			attempt = strconv.Itoa(sub.Attempt)
		} else {
			grade = "0"
			attempt = "0"
		}
		records = append(records, []string{student.First, student.Last, grade, attempt, "0"})
	}

	csvBytes := &bytes.Buffer{}
	writer := csv.NewWriter(csvBytes)
	err = writer.WriteAll(records)
	if err != nil {
		return nil, "", 0, errors.ErrorFailedToWriteCSV
	}
	writer.Flush()

	return csvBytes, "filename", int64(csvBytes.Len()), nil
}
