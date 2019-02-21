package auth

import (
	"github.com/gin-gonic/gin"

	"github.com/appleboy/gin-jwt"

	"backend/forms"
)

// Authenticator a default function for a gin jwt, that authenticates a user.
func Authenticator(c *gin.Context) (interface{}, error) {
	var login forms.UserLoginForm
	if errs := c.ShouldBindJSON(&login); errs != nil {
		return "Missing login values.", jwt.ErrMissingLoginValues
	}
	val, err := um.Login(login)
	return val, err
}
