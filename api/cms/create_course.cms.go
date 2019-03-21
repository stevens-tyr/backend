package cms

import (
//	"fmt"
	"time"
	
	"github.com/gin-gonic/gin"
	"github.com/dgrijalva/jwt-go"
	

	"backend/api/auth"
	"backend/errors"
	"backend/forms"
)

func CreateCourse(c *gin.Context) {
	uid, _ := c.Get("uid")

	var createCourse forms.CreateCourseForm
	if errs := c.ShouldBindJSON(&createCourse); errs != nil {
		c.Set("error", errors.ErrorInvlaidJSON)
		return
	}

	cid, err := cm.Create(uid, createCourse)
	if err != nil {
		c.Set("error", err)
		return
	}

	err = um.AddCourse("teacher", cid, uid)
	if err != nil {
		c.Set("error", err)
		return
	}

	user, err := um.FindOneById(uid)
	if err != nil {
		c.Set("error", err)
		return
	}
	
	token := jwt.New(jwt.GetSigningMethod(auth.AuthMiddleware.SigningAlgorithm))
	claims := token.Claims.(jwt.MapClaims)
	for key, val := range auth.AuthMiddleware.PayloadFunc(user) {
		claims[key] = val
	}
	expire := auth.AuthMiddleware.TimeFunc().Add(auth.AuthMiddleware.Timeout)
	claims["exp"] = expire.Unix()
	claims["orig_iat"] = auth.AuthMiddleware.TimeFunc().Unix()
	tokenString, _ := token.SignedString(auth.AuthMiddleware.Key)
	
	c.SetCookie(
		auth.AuthMiddleware.CookieName,
		tokenString,
		int(expire.Unix() - time.Now().Unix()),
		"/",
		auth.AuthMiddleware.CookieDomain,
		auth.AuthMiddleware.SecureCookie,
		auth.AuthMiddleware.CookieHTTPOnly,
	)
	
	c.JSON(200, gin.H{
		"message": "Course created.",
		"token":   tokenString,
		"expire":  expire,
	})
}
