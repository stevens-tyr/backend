package auth

import "github.com/gin-gonic/gin"

// Authorizator a default function for a gin jwt, that authorizes a user.
func Authorizator(d interface{}, c *gin.Context) bool {
	//claims := jwt.ExtractClaims(c)

	// todo in future check
	// look at route see if what part of course/assignment they are accessing
	// check if they have permission claims["courses"] slice of enrolledCourse

	return true
}
