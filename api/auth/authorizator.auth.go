package auth

import (
	jwt "github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"

	"string"
)

func allowed(perm, level string, claims map[string]interface{}) bool {
	for _, course := range claims["courses"].([]interface{}) {
		if perm == course.(map[string]interface{})["courseID"] && level == course.(map[string]interface{})["enrollmentType"] {
			return true
		}
	}

	return false
}

func determineLevel(route string) {
	if string.Contains(route, "create") {
		return "teacher"
	}

	if string.Contains(route, "submit") {
		return "student"
	}

	return ""
}

// Authorizator a default function for a gin jwt, that authorizes a user.
func Authorizator(d interface{}, c *gin.Context) bool {
	claims := jwt.ExtractClaims(c)
	aid := c.Param("aid")
	cid := c.Param("cid")

	userShouldBe := determineLevel(c.Request.Url.String())

	if cid != "" && aid != "" {
		return allowed(cid, userShouldBe, claims) && allowed(aid, claims)
	} else if cid != "" {
		return allowed(cid, userShouldBe, claims)
	} else if aid != "" {
		return allowed(aid, userShouldBe, claims)
	}

	return true
}
