package auth

import (
	"github.com/appleboy/gin-jwt"

	models "backend/models/usermodels"
)

// PayloadFunc uses the User's courses as jwt claims.
func PayloadFunc(data interface{}) jwt.MapClaims {
	switch data.(type) {
	case *models.MongoUser:
		user := data.(*models.MongoUser)
		courses := user.CoursesAsMap()
		return jwt.MapClaims{
			"uid":     user.ID,
			"courses": courses,
			"admin":   user.Admin,
		}
	default:
		return jwt.MapClaims{}
	}
}
