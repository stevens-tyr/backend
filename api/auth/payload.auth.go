package auth

import (
	jwt "github.com/appleboy/gin-jwt"

	"backend/models"
)

// PayloadFunc uses the User's courses as jwt claims.
func PayloadFunc(data interface{}) jwt.MapClaims {
	switch data.(type) {
	case models.User:
		user := data.(models.User)
		return jwt.MapClaims{
			"uid":         user.ID,
			"courses":     user.EnrolledCourses,
			"assignments": make([]string, 0),
		}
	default:
		return jwt.MapClaims{}
	}
}
