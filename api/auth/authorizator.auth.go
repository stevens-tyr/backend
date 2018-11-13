package auth

import (
	jwt "github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
)

func allowed(perm string, claims map[string]interface{}) bool {
	for _, course := range claims["courses"].([]interface{}) {
		if perm == course.(map[string]interface{})["courseID"] {
			return true
		}
	}

	return false
}

// Authorizator a default function for a gin jwt, that authorizes a user.
func Authorizator(d interface{}, c *gin.Context) bool {
	claims := jwt.ExtractClaims(c)
	aid := c.Param("aid")
	cid := c.Param("cid")

	if cid != "" && aid != "" {
		return allowed(cid, claims) && allowed(aid, claims)
	} else if cid != "" {
		return allowed(cid, claims)
	} else if aid != "" {
		return allowed(aid, claims)
	}

	return true
}
