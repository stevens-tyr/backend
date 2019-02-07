package usermodels

import (
	"errors"
	"context"
	"os"

	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/bson/primitive"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo/options"
	bcrypt "golang.org/x/crypto/bcrypt"

	"github.com/stevens-tyr/tyr-gin"
	"backend/forms"
)

var (
	// UserNotFoundError an error to throw for when a User is not found.
	ErrorUserNotFound = errors.New("USER DOES NOT EXIST")
	// IncorrectPasswordError an error to throw for when an inccorect passowrd is entered.
	ErrorIncorrectPassword = errors.New("INCORRECT PASSWORD")
	// ErrorNonMatchingPassword an error to throw when a password cofirmation does not match the password.
	ErrorNonMatchingPassword = errors.New("CONFIRMATION MUST MATCH")
	// ErrorFailedToCreateUser an error for when you fail to create a user.
	ErrorFailedToCreateUser = errors.New("FAILED TO CREATE USER")
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
		Email           string             `bson:"email" json:"email" binding:"required"`
		Password        []byte             `bson:"password" json:"password" binding:"required"`
		First           string             `bson:"firstName" json:"first_name" binding:"required"`
		Last            string             `bson:"lastName" json:"last_name" binding:"required"`
		EnrolledCourses []EnrolledCourse   `bson:"enrolledCourses" json:"enrolledCourses" binding:"required"`
	}

	// A struct to represent a bunch of User functions.
	UserInterface struct{
		ctx context.Context
		col *mongo.Collection
	}
)

func New() *UserInterface {
	db, _ := tyrgin.GetMongoDB(os.Getenv("DB_NAME"))
	col := tyrgin.GetMongoCollection("users", db) 

	return &UserInterface{
		context.Background(),
		col,
	}
}

func (u *UserInterface) FindOne(email string) (*MongoUser, error){
	var user *MongoUser

	res := u.col.FindOne(u.ctx, bson.M{"email": email}, options.FindOne())
	res.Decode(&user)

	if user == nil {
		return nil, ErrorUserNotFound
	}

	return user, nil
}

func (u *UserInterface) Login(form forms.UserLoginForm) (interface{}, error) {
	user, err := u.FindOne(form.Email)
	if err != nil{
		return "User not found.", err
	}

	if err = bcrypt.CompareHashAndPassword(user.Password, []byte(form.Password)); err != nil {
		return "Incorrect password", ErrorIncorrectPassword
	}

	return user, nil
}

func (u *UserInterface) Register(form forms.UserRegisterForm) (error) {
	user, err := u.FindOne(form.Email)
	if err != nil && err != ErrorUserNotFound {
		return err
	}

	if user != nil {
		return errors.New("user with email already exists")
	}

	if form.Password != form.PasswordConfirmation {
		return ErrorNonMatchingPassword
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(form.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil
	}

	user = &MongoUser{
		Email:           form.Email,
		Password:        hash,
		First:           form.First,
		Last:            form.Last,
		EnrolledCourses: make([]EnrolledCourse, 0),
	}

	_, err = u.col.InsertOne(u.ctx, user, options.InsertOne())
	if err != nil {
		return ErrorFailedToCreateUser
	}

	return nil
}

func (u *UserInterface) GetCourses(uid primitive.ObjectID) ([]forms.CourseAggQuery, error) {
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
		return courses, err
	}

	for cur.Next(u.ctx) {
		var course map[string]forms.CourseAggQuery
		err = cur.Decode(&course)
		if err != nil {
			return courses, err
		}
		courses = append(courses, course["course"])
	}
	return courses, nil
}
