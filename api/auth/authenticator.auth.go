package auth

import (
	"os"

	jwt "github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
	"github.com/stevens-tyr/tyr-gin"
	bcrypt "golang.org/x/crypto/bcrypt"
	bson "gopkg.in/mgo.v2/bson"

	"backend/models"
)

// Authenticator a default function for a gin jwt, that authenticates a user.
func Authenticator(c *gin.Context) (interface{}, error) {
	var login models.Login
	if err := c.ShouldBindJSON(&login); err != nil {
		return "Missing login values.", jwt.ErrMissingLoginValues
	}

	db, _ := tyrgin.GetMongoDB(os.Getenv("DB_NAME"))
	col, err := tyrgin.GetMongoCollectionCreate("users", db)

	var user models.User
	if err = col.Find(bson.M{"email": login.Email}).One(&user); err != nil {
		return "User not found.", models.ErrorUserNotFound
	}

	if err = bcrypt.CompareHashAndPassword(user.Password, []byte(login.Password)); err != nil {
		return "Incorrect password", models.ErrorIncorrectPassword
	}

	return user, nil
}
