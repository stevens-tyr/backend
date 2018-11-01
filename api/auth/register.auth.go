package auth

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/goware/emailx"
	"github.com/stevens-tyr/tyr-gin"
	"golang.org/x/crypto/bcrypt"
	mgo "gopkg.in/mgo.v2"
	bson "gopkg.in/mgo.v2/bson"

	"backend/models"
)

// IsValidEmail checks an email string to be valid and with resolvable host.
func IsValidEmail(email string) error {

	err := emailx.Validate(email)
	if err != nil {
		if err == emailx.ErrInvalidFormat {
			return models.ErrorEmailNotValid
		}
		if err == emailx.ErrUnresolvableHost {
			return models.ErrorUnresolvableEmailHost
		}
		return err
	}
	return nil
}

// Register a function that registers a User.
func Register(c *gin.Context) {
	var register models.RegisterForm
	if err := c.ShouldBindJSON(&register); err != nil {
		c.JSON(400, gin.H{
			"status_code": 400,
			"message":     "Incorrect json format.",
		})
		return
	}

	db, err := tyrgin.GetMongoDB(os.Getenv("DB_NAME"))
	if err != nil {
		c.JSON(500, gin.H{
			"status_code": 500,
			"message":     "Failed to get Mongo Session.",
		})
		return
	}

	col, err := tyrgin.GetMongoCollectionCreate("users", db)
	if err != nil {
		c.JSON(500, gin.H{
			"status_code": 500,
			"message":     "Failed to get collection.",
		})
		return
	}

	if err = IsValidEmail(register.Email); err != nil {
		msg := "Email is invalid"
		if err == models.ErrorUnresolvableEmailHost {
			msg = "Unable to resolve email host"
		}
		c.JSON(400, gin.H{
			"status_code": 400,
			"message":     msg,
		})
		return
	}

	var user models.User
	if err = col.Find(bson.M{"email": register.Email}).One(&user); err != mgo.ErrNotFound {
		c.JSON(400, gin.H{
			"status_code": 400,
			"message":     "Email is taken.",
		})
		return
	}

	if register.Password != register.PasswordConfirmation {
		c.JSON(400, gin.H{
			"status_code": 400,
			"message":     "Your password and password confirmation do not match.",
		})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(register.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(500, gin.H{
			"status_code": 500,
			"message":     "Failed to generate hash",
		})
		return
	}

	user = models.User{
		Email:           register.Email,
		Password:        hash,
		First:           register.First,
		Last:            register.Last,
		EnrolledCourses: make([]models.EnrolledCourse, 0),
	}

	err = col.Insert(&user)
	if err != nil {
		c.JSON(500, gin.H{
			"status_code": 500,
			"message":     "Failed to create user.",
		})
		return
	}

	c.JSON(200, gin.H{
		"status_code": 200,
		"message":     "User created.",
	})

}
