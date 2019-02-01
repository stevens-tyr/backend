package auth

import (
	ctx "context"
	"fmt"
	"os"

	jwt "github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo/options"
	bcrypt "golang.org/x/crypto/bcrypt"

	"backend/models"

	"github.com/stevens-tyr/tyr-gin"
)

// Authenticator a default function for a gin jwt, that authenticates a user.
func Authenticator(c *gin.Context) (interface{}, error) {
	var login models.Login
	if err := c.ShouldBindJSON(&login); err != nil {
		return "Missing login values.", jwt.ErrMissingLoginValues
	}

	db, _ := tyrgin.GetMongoDB(os.Getenv("DB_NAME"))
	col := tyrgin.GetMongoCollection("users", db)

	var user models.User
	res := col.FindOne(ctx.Background(), bson.M{"email": login.Email}, options.FindOne())

	err := res.Decode(&user)
	fmt.Println(user)
	if user.Email == "" {
		return "User not found.", models.ErrorUserNotFound
	}

	if err = bcrypt.CompareHashAndPassword(user.Password, []byte(login.Password)); err != nil {
		return "Incorrect password", models.ErrorIncorrectPassword
	}

	return user, nil
}
