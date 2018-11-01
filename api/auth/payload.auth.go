package auth

import (
	jwt "github.com/appleboy/gin-jwt"

	"backend/models"
)

// PayloadFunc uses the User's courses as jwt claims.
func PayloadFunc(data interface{}) jwt.MapClaims {
	switch data.(type) {
	case models.User:
		return jwt.MapClaims{"courses": data.(models.User).EnrolledCourses}
	default:
		return jwt.MapClaims{}
	}
}
