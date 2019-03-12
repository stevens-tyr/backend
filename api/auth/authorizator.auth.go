package auth

import (
	"fmt"
	"strings"

	"github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
	"github.com/mongodb/mongo-go-driver/bson/primitive"
)

func allowed(levels []string, claims map[string]interface{}, c *gin.Context) bool {
	enrolledCourses := claims["courses"].(map[string]interface{})
	uid := claims["uid"]
	cid, _ := c.Get("cids")
	sid, exists := c.Get("sid")
	allowed := false

	val, found := enrolledCourses[cid.(string)]
	c.Set("role", val)
	fmt.Println("cid:", cid, levels, val)
	if found && (in(levels, "any") || in(levels, val.(string))) {
		return true
	}

	if in(levels, "student") && exists {
		sub, err := sm.GetUsersSubmission(sid, uid)
		if err != nil {
			fmt.Println("err", err)
		}
		if sub != nil {
			allowed = true
		}
	}

	return allowed
}

func in(terms []string, term string) bool {
	for _, val := range terms {
		if val == term {
			return true
		}
	}

	return false
}

func determineLevel(route string) []string {
	var allowed []string
	if _, found := routeLevels["admin"][route]; found {
		allowed = append(allowed, "admin")
	}

	if _, found := routeLevels["any"][route]; found {
		allowed = append(allowed, "any")
	}

	if _, found := routeLevels["assistant"][route]; found {
		allowed = append(allowed, "assistant")
	}

	if _, found := routeLevels["teacher"][route]; found {
		allowed = append(allowed, "teacher")
	}

	if _, found := routeLevels["student"][route]; found {
		allowed = append(allowed, "student")
	}

	if len(allowed) == 0 {
		allowed = append(allowed, "whitelisted")
	}

	return allowed
}

// Authorizator a default function for a gin jwt, that authorizes a user.
func Authorizator(d interface{}, c *gin.Context) bool {
	route := strings.TrimPrefix(c.Request.URL.String(), "/api/v1/plague_doctor/")
	for _, p := range c.Params {
		route = strings.Replace(route, p.Value, ":"+p.Key, 1)
	}

	claims := jwt.ExtractClaims(c)
	uids := claims["uid"].(string)
	val, _ := primitive.ObjectIDFromHex(uids)
	c.Set("uid", val)

	userLevelForRouteShouldBe := determineLevel(route)
	fmt.Println("user level:", userLevelForRouteShouldBe, route)
	if in(userLevelForRouteShouldBe, "whitelisted") {
		return true
	}

	admin := claims["admin"].(bool)
	if in(userLevelForRouteShouldBe, "admin") && admin {
		return true
	} else if in(userLevelForRouteShouldBe, "admin") && !admin {
		return false
	}

	return allowed(userLevelForRouteShouldBe, claims, c)
}
