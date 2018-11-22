package api

import (
	"backend/api/auth"
	"backend/api/cms"

	"github.com/gin-gonic/gin"
	"github.com/stevens-tyr/tyr-gin"
)

func SetUp() *gin.Engine {
	server := tyrgin.SetupRouter()

	tyrgin.ServeReact(server)

	server.MaxMultipartMemory = 50 << 20

	server.Use(tyrgin.Logger())
	server.Use(gin.Recovery())

	var authEndpoints = []tyrgin.APIAction{
		tyrgin.NewRoute(auth.AuthMiddleware.LoginHandler, "login", false, tyrgin.POST),
		tyrgin.NewRoute(auth.AuthMiddleware.RefreshHandler, "refresh_token", false, tyrgin.GET),
		tyrgin.NewRoute(auth.Register, "register", false, tyrgin.POST),
		tyrgin.NewRoute(func(c *gin.Context) {
			c.JSON(200, gin.H{
				"status_code": 200,
				"msg":         "Success.",
			})
		}, "logged_in", true, tyrgin.GET),
	}

	tyrgin.AddRoutes(server, auth.AuthMiddleware, "1", "auth", authEndpoints)

	var cmsEndpoints = []tyrgin.APIAction{
		tyrgin.NewRoute(cms.CreateAssignment, "course/:cid/:section/assignment/create", true, tyrgin.POST),
		tyrgin.NewRoute(cms.SubmitAssignment, "course/:cid/:section/assignment/submit/:aid", true, tyrgin.POST),
		tyrgin.NewRoute(cms.Dashboard, "dashboard", true, tyrgin.GET),
	}

	tyrgin.AddRoutes(server, auth.AuthMiddleware, "1", "plague_doctor", cmsEndpoints)

	server.NoRoute(tyrgin.NotFound)

	return server
}
