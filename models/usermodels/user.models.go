package usermodels

import (
	"context"
	"fmt"
	"os"

	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/primitive"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/mongo/options"
	bcrypt "golang.org/x/crypto/bcrypt"

	"backend/errors"
	"backend/forms"

	"github.com/stevens-tyr/tyr-gin"
)

type (
	// EnrolledCourse struct keeps track of a user's course and enrollment type
	EnrolledCourse struct {
		CourseID       primitive.ObjectID `bson:"courseID" json:"courseID" binding:"required"`
		EnrollmentType string             `bson:"enrollmentType" json:"enrollmentType" binding:"required"`
	}

	// User a default User struct to represent a User in Tyr.
	MongoUser struct {
		ID              primitive.ObjectID `bson:"_id,omitempty" json:"id" biding:"required"`
		Admin           bool               `bson:"admin" json:"admin`
		Email           string             `bson:"email" json:"email" binding:"required"`
		Password        []byte             `bson:"password" json:"password" binding:"required"`
		First           string             `bson:"firstName" json:"first_name" binding:"required"`
		Last            string             `bson:"lastName" json:"last_name" binding:"required"`
		EnrolledCourses []EnrolledCourse   `bson:"enrolledCourses" json:"enrolledCourses" binding:"required"`
	}

	// A struct to represent a bunch of User functions.
	UserInterface struct {
		ctx context.Context
		col *mongo.Collection
	}
)

func (m *MongoUser) CoursesAsMap() map[string]string {
	courses := make(map[string]string)

	for _, course := range m.EnrolledCourses {
		courses[course.CourseID.Hex()] = course.EnrollmentType
	}

	return courses
}

func New() *UserInterface {
	db, _ := tyrgin.GetMongoDB(os.Getenv("DB_NAME"))
	col := tyrgin.GetMongoCollection("users", db)

	return &UserInterface{
		context.Background(),
		col,
	}
}

func (u *UserInterface) FindOne(email string) (*MongoUser, errors.APIError) {
	var user *MongoUser

	res := u.col.FindOne(u.ctx, bson.M{"email": email}, options.FindOne())
	res.Decode(&user)
	fmt.Println("user", user)

	if user == nil {
		return nil, errors.ErrorResourceNotFound
	}

	return user, nil
}

func (u *UserInterface) Login(form forms.UserLoginForm) (interface{}, errors.APIError) {
	user, err := u.FindOne(form.Email)
	if err != nil {
		return "User not found.", err
	}

	if er := bcrypt.CompareHashAndPassword(user.Password, []byte(form.Password)); er != nil {
		return "Incorrect password", errors.ErrorIncorrectCredentials
	}

	return user, nil
}

func (u *UserInterface) Register(form forms.UserRegisterForm) errors.APIError {
	user, err := u.FindOne(form.Email)
	if err != nil && err != errors.ErrorResourceNotFound {
		return err
	}

	if user != nil {
		return errors.ErrorCannotCreateDuplicateData
	}

	if form.Password != form.PasswordConfirmation {
		return errors.ErrorIncorrectCredentials
	}

	hash, errs := bcrypt.GenerateFromPassword([]byte(form.Password), bcrypt.DefaultCost)
	if errs != nil {
		return errors.ErrorHashFailure
	}

	user = &MongoUser{
		Email:           form.Email,
		Admin:           false,
		Password:        hash,
		First:           form.First,
		Last:            form.Last,
		EnrolledCourses: make([]EnrolledCourse, 0),
	}

	_, errs = u.col.InsertOne(u.ctx, user, options.InsertOne())
	if errs != nil {
		return errors.ErrorDatabaseFailedCreate
	}

	return nil
}

func (u *UserInterface) GetCourses(uid interface{}) ([]forms.CourseAggQuery, errors.APIError) {
	query := []interface{}{
		bson.M{"$match": bson.M{"_id": uid}},
		bson.M{"$unwind": "$enrolledCourses"},
		bson.M{
			"$lookup": bson.M{
				"from":         "courses",
				"localField":   "enrolledCourses.courseID",
				"foreignField": "_id",
				"as":           "course",
			},
		},
		bson.M{"$project": bson.M{
			"_id": 0,
			"course": bson.M{
				"$arrayElemAt": bson.A{"$course", 0},
			},
		},
		},
	}

	var courses []forms.CourseAggQuery
	cur, err := u.col.Aggregate(u.ctx, query, options.Aggregate())
	if err != nil {
		return courses, errors.ErrorDatabaseFailedQuery
	}

	for cur.Next(u.ctx) {
		var course map[string]forms.CourseAggQuery
		err = cur.Decode(&course)
		if err != nil {
			return courses, errors.ErrorDatabaseFailedExtract
		}
		courses = append(courses, course["course"])
	}
	return courses, nil
}

func (u *UserInterface) CourseExists(cid, uid interface{}) (bool, errors.APIError) {
	filter := bson.D{
		{"_id", uid},
		{
			"enrolledCourses", bson.D{
				{
					"$elemMatch", bson.D{{"courseID", cid}},
				},
			},
		},
	}

	res := u.col.FindOne(
		u.ctx,
		filter,
		options.FindOne(),
	)
	var user *MongoUser
	err := res.Decode(&user)
	if err != nil {
		return false, errors.ErrorDatabaseFailedExtract
	}
	if user != nil {
		return true, nil
	}

	return false, nil
}

func (u *UserInterface) AddCourse(level string, cid, uid interface{}) errors.APIError {
	alreadyEnrolled, _ := u.CourseExists(cid, uid)
	if alreadyEnrolled {
		return errors.ErrorUserAlreadyEnrolled
	}

	_, err := u.col.UpdateOne(
		u.ctx,
		bson.M{"_id": uid},
		bson.M{"$push": bson.M{"enrolledCourses": bson.M{"courseID": cid, "enrollmentType": level}}},
		options.Update(),
	)
	if err != nil {
		return errors.ErrorDatabaseFailedUpdate
	}

	return nil
}
