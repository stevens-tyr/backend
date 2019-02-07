package api

import (
	"backend/api/auth"
	"backend/api/cms"

	"github.com/gin-gonic/gin"
	"github.com/stevens-tyr/tyr-gin"
)

// SetUp is a function to set up the routes for plague doctor microservice.
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
		tyrgin.NewRoute(auth.Check, "logged_in", true, tyrgin.GET),
	}

	tyrgin.AddRoutes(server, auth.AuthMiddleware, "1", "auth", authEndpoints)

	var cmsEndpoints = []tyrgin.APIAction{
		tyrgin.NewRoute(cms.CourseAssignments, "course/:cid/:section/assignments", true, tyrgin.GET),
		tyrgin.NewRoute(cms.CreateAssignment, "course/:cid/:section/assignment/create", true, tyrgin.POST),
		tyrgin.NewRoute(cms.Dashboard, "dashboard", true, tyrgin.GET),
		tyrgin.NewRoute(cms.DownloadSubmission, "course/:cid/:section/assignment/:aid/submission/download/:sid/:num", true, tyrgin.GET),
		tyrgin.NewRoute(cms.GetAssignment, "course/:cid/:section/assignment/:aid/details", true, tyrgin.GET),
		tyrgin.NewRoute(cms.SubmitAssignment, "course/:cid/:section/assignment/submit/:aid", true, tyrgin.POST),
	}

	tyrgin.AddRoutes(server, auth.AuthMiddleware, "1", "plague_doctor", cmsEndpoints)

	server.NoRoute(tyrgin.NotFound)

	return server
}
