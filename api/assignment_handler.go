package api

import (
	"github.com/gin-gonic/gin"
)

// CreateAssignment will create an assignment.
// TBD
func CreateAssignment(c *gin.Context) {
	var ca createAssignment
	if err := c.ShouldBindJson(&ca); err != nil {
		return
	}

	return
}

// SubmitAssignment will create an assignment.
// TBD
func SubmitAssignment(c *gin.Context) {
	var sa submitAssignment
	if err := c.ShouldBindJson(&sa); err != nil {
		return
	}
}
