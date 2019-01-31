package models

import (
	"errors"

	"github.com/mongodb/mongo-go-driver/bson/primitive"
)

var (
	// ErrorEmailNotValid an error to throw when an email format is not valid
	ErrorEmailNotValid = errors.New("EMAIL NOT VALID")
	// ErrorUnresolvedEmailHost an error to throw when the email host is unresolvable
	ErrorUnresolvableEmailHost = errors.New("EMAIL HOST UNRESOLVABLE")
	// UserNotFoundError an error to throw for when a User is not found.
	ErrorUserNotFound = errors.New("USER DOES NOT EXIST")
	// IncorrectPasswordError an error to throw for when an inccorect passowrd is entered.
	ErrorIncorrectPassword = errors.New("INCORRECT PASSWORD")
)

// Mongo Schemas
type (
	// Email struct to get the email of a User.
	Email struct {
		Email string `bson:"email" json:"email" binding:"required"`
	}

	// Login struct a form to login a Tyr User.
	Login struct {
		Email    string `bson:"email" json:"email" binding:"required"`
		Password string `bson:"password" json:"password" binding:"required"`
	}

	// RegisterForm struct a form for register a Tyr User.
	RegisterForm struct {
		Email                string `bson:"email" json:"email" binding:"required"`
		Password             string `bson:"password" json:"password" binding:"required"`
		PasswordConfirmation string `bson:"passwordConfirmation" json:"passwordConfirmation" binding:"required"`
		First                string `bson:"firstName" json:"firstName" binding:"required"`
		Last                 string `bson:"lastName" json:"lastName" binding:"required"`
	}

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
	Submission struct {
		ID             primitive.ObjectID  `bson:"_id" json:"id" binding:"required"`
		UserID         primitive.ObjectID  `bson:"userID" json:"userID" binding:"required"`
		AttemptNumber  int                `bson:"attemptNumber" json:"attemptNumber" binding:"required"`
		SubmissionDate primitive.DateTime `bson:"submissionDate" json:"submissionDate" binding:"required"`
		File           string             `bson:"file" json:"file" binding:"required"`
		ErrorTesting   bool               `bson:"errorTesting" json:"errorTesting" binding:"required"`
		Cases          Cases              `bson:"cases" json:"cases" binding:"required"`
	}

	// TestScripts struct to represent filenames on gcp storage of scripts.
	TestScripts struct {
		StudentFacing string `bson:"studentFacing" json:"studentFacing" binding:"required"`
		AdminFacing   string `bson:"adminFacing" json:"adminFacing" binding:"required"`
	}

	AssignmentSubmission struct {
		UserID        primitive.ObjectID `bson:"userID" json:"userID" binding:"required"`
		SubmissionID  primitive.ObjectID `bson:"submissionID" json:"submissionID" binding:"required"`
		AttemptNumber int               `bson:"attemptNumber" json:"attemptNumber" binding:"required"`
	}

	// Assignment struct to store information about an assignment.
	Assignment struct {
		ID              primitive.ObjectID      `bson:"_id" json:"id" binding:"required"`
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

	// Course struct ot store information about a course.
	Course struct {
		ID          primitive.ObjectID   `bson:"_id,omitempty" json:"id" binding:"required"`
		Department  string              `bson:"department" json:"department" binding:"required"`
		Number      int                 `bson:"number" json:"number" binding:"required"`
		Section     string              `bson:"section" json:"section" binding:"required"`
		Professors  []primitive.ObjectID `bson:"professors" json:"professors" binding:"required"`
		Assistants  []primitive.ObjectID `bson:"assistants" json:"assitants" binding:"required"`
		Students    []primitive.ObjectID `bson:"students" json:"students" binding:"required"`
		Assignments []primitive.ObjectID `bson:"assignments" json:"assignments" binding:"required"`
	}

	// EnrolledCourse struct keeps track of a user's course and enrollment type
	EnrolledCourse struct {
		CourseID       primitive.ObjectID `bson:"courseID" json:"courseID" binding:"required"`
		EnrollmentType string            `bson:"enrollmentType" json:"enrollmentType" binding:"required"`
	}

	// User a default User struct to represent a User in Tyr.
	User struct {
		ID              primitive.ObjectID `bson:"_id,omitempty" json:"id" biding:"required"`
		Email           string            `bson:"email" json:"email" binding:"required"`
		Password        []byte            `bson:"password" json:"password" binding:"required"`
		First           string            `bson:"firstName" json:"first_name" binding:"required"`
		Last            string            `bson:"lastName" json:"last_name" binding:"required"`
		EnrolledCourses []EnrolledCourse  `bson:"enrolledCourses" json:"enrolledCourses" binding:"required"`
	}
)
