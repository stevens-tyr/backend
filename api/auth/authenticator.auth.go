package auth

import (
	"github.com/gin-gonic/gin"

	jwt "github.com/appleboy/gin-jwt"

	"backend/forms"
)

// Authenticator a default function for a gin jwt, that authenticates a user.
func Authenticator(c *gin.Context) (interface{}, error) {
	var login forms.UserLoginForm
	if err := c.ShouldBindJSON(&login); err != nil {
		return "Missing login values.", jwt.ErrMissingLoginValues
	}

	return um.Login(login)
}
