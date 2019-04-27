package api

import (
	"backend/api/auth"
	"backend/api/cms"
	"backend/middleware"

	"github.com/gin-gonic/gin"
	"github.com/stevens-tyr/tyr-gin"
)

// SetUp is a function to set up the routes for plague doctor microservice.
func SetUp() *gin.Engine {
	server := tyrgin.SetupRouter()

	tyrgin.ServeReact(server)

	server.MaxMultipartMemory = 50 << 20

	server.Use(middleware.ObjectIDs())
	server.Use(middleware.ErrorHandler())

	var authEndpoints = []tyrgin.APIAction{
		tyrgin.NewRoute(auth.AuthMiddleware.LoginHandler, "login", tyrgin.POST),
		tyrgin.NewRoute(auth.AuthMiddleware.RefreshHandler, "refresh_token", tyrgin.GET),
		tyrgin.NewRoute(auth.Register, "register", tyrgin.POST),
	}

	var secureAuthEndpoints = []tyrgin.APIAction{
		tyrgin.NewRoute(auth.Check, "logged_in", tyrgin.GET),
		tyrgin.NewRoute(auth.Logout, "logout", tyrgin.GET),
	}
	tyrgin.AddRoutes(server, false, auth.AuthMiddleware, "1", "auth", authEndpoints)
	tyrgin.AddRoutes(server, true, auth.AuthMiddleware, "1", "auth", secureAuthEndpoints)

	var secureCmsEndpoints = []tyrgin.APIAction{
		tyrgin.NewRoute(cms.AssignmentAsFile, "course/:cid/assignment/:aid/file", tyrgin.GET),
		tyrgin.NewRoute(cms.CourseAssignments, "course/:cid/assignments", tyrgin.GET),
		tyrgin.NewRoute(cms.CourseAddUser, "course/:cid/add/user", tyrgin.POST),
		tyrgin.NewRoute(cms.CourseAddUsers, "course/:cid/add/users", tyrgin.POST),
		tyrgin.NewRoute(cms.CreateAssignment, "course/:cid/assignment/create", tyrgin.POST),
		tyrgin.NewRoute(cms.CreateAssignmentFromFile, "course/:cid/assignment/create/file", tyrgin.POST),
		tyrgin.NewRoute(cms.CreateCourse, "create/course", tyrgin.POST),
		tyrgin.NewRoute(cms.Dashboard, "dashboard", tyrgin.GET),
		tyrgin.NewRoute(cms.DeleteAssignment, "course/:cid/assignment/:aid/delete", tyrgin.DELETE),
		tyrgin.NewRoute(cms.DeleteCourse, "course/:cid/delete", tyrgin.DELETE),
		tyrgin.NewRoute(cms.DownloadSubmission, "course/:cid/assignment/:aid/submission/:sid/:num", tyrgin.GET),
		tyrgin.NewRoute(cms.GetAssignment, "course/:cid/assignment/:aid/details", tyrgin.GET),
		tyrgin.NewRoute(cms.GetCourse, "course/:cid", tyrgin.GET),
		tyrgin.NewRoute(cms.GradesAsCSV, "course/:cid/assignment/:aid/csv", tyrgin.GET),
		tyrgin.NewRoute(cms.SubmitAssignment, "course/:cid/assignment/submit/:aid", tyrgin.POST),
		tyrgin.NewRoute(cms.UpdateAssignment, "course/:cid/assignment/:aid/update", tyrgin.PATCH),
		tyrgin.NewRoute(cms.UpdateCourse, "course/:cid/update", tyrgin.PATCH),
	}

	var cmsEndpoints = []tyrgin.APIAction{
		tyrgin.NewRoute(cms.UpdateGrade, "job/:secret/submission/:sid/update", tyrgin.PATCH),
		tyrgin.NewRoute(cms.UpdateGradeError, "job/:secret/submission/:sid/error", tyrgin.PATCH),
		tyrgin.NewRoute(cms.JobDownloadSubmission, "job/:secret/submission/:sid/download", tyrgin.GET),
		tyrgin.NewRoute(cms.JobDownloadSupportingFiles, "job/:secret/assignment/:aid/supportingfiles/download", tyrgin.GET),
		tyrgin.NewRoute(auth.Register, "register", tyrgin.POST),
	}

	tyrgin.AddRoutes(server, true, auth.AuthMiddleware, "1", "plague_doctor", secureCmsEndpoints)
	tyrgin.AddRoutes(server, false, auth.AuthMiddleware, "1", "plague_doctor", cmsEndpoints)

	server.NoRoute(tyrgin.NotFound)

	return server
}
