package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/goware/emailx"

	"backend/errors"
	"backend/forms"
	"backend/models"
)

// isValidEmail checks an email string to be valid and with resolvable host.
func isValidEmail(email string) error {
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
	var register forms.UserRegisterForm
	err := c.ShouldBindJSON(&register)
	if err != nil {
		c.Set("error", errors.ErrorInvlaidJSON)
		return
	}

	err = um.Register(register)
	if err != nil {
		c.Set("error", err)
		return
	}

	c.JSON(200, gin.H{
		"message": "User created.",
	})
}
