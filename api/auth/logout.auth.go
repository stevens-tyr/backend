package auth

import "github.com/gin-gonic/gin"

func Logout(c *gin.Context) {
	c.SetCookie(
		AuthMiddleware.CookieName,
		"",
		-1,
		"/",
		AuthMiddleware.CookieDomain,
		AuthMiddleware.SecureCookie,
		AuthMiddleware.CookieHTTPOnly,
	)

	c.JSON(200, gin.H{
		"message": "Logged Out.",
	})
}
