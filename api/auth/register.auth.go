package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/goware/emailx"

	"backend/models"
	forms "backend/forms"

	"github.com/stevens-tyr/tyr-gin"
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
			tyrgin.ErrorHandler(err, c, 400, gin.H{
			"status_code": 400,
			"message":     err.Error(),
		})
	}

	err = um.Register(register)
	if err != nil {
		tyrgin.ErrorHandler(err, c, 400, gin.H{
		"status_code": 400,
		"message":     err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"status_code": 200,
		"message":     "User created.",
	})
}
