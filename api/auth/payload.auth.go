package auth

import (
	"fmt"
	jwt "github.com/appleboy/gin-jwt"

	models "backend/models/usermodels"
)

// PayloadFunc uses the User's courses as jwt claims.
func PayloadFunc(data interface{}) jwt.MapClaims {
	fmt.Println("login data sent type", data)
	switch data.(type) {
	case *models.MongoUser:
		user := data.(*models.MongoUser)
		return jwt.MapClaims{
			"uid":         user.ID,
			"courses":     user.EnrolledCourses,
			"assignments": make([]string, 0),
		}
	default:
		fmt.Println("yo i am technically not a models user")
		return jwt.MapClaims{}
	}
}