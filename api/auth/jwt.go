package auth

import (
	"os"
	"time"

	jwt "github.com/appleboy/gin-jwt"

	"backend/models"
)

var um = models.NewMongoUserInterface()
var sm = models.NewMongoSubmissionInterface()

// AuthMiddleware is a jwt middleware for auth requests
var AuthMiddleware, _ = jwt.New(&jwt.GinJWTMiddleware{
	Realm:           os.Getenv("JWT_REALM"),
	Key:             []byte(os.Getenv("JWT_SECRET")),
	Timeout:         time.Hour,
	MaxRefresh:      time.Hour * 24,
	Authenticator:   Authenticator,
	Authorizator:    Authorizator,
	LoginResponse:   TokenResponse,
	PayloadFunc:     PayloadFunc,
	RefreshResponse: TokenResponse,
	Unauthorized:    Unauthorized,
	TokenLookup:     "header:Authorization, cookie:JWTToken",
	TokenHeadName:   "Bearer",
	TimeFunc:        time.Now,
	SendCookie:      true,
	SecureCookie:    false,
	//CookieHTTPOnly: true,
	//CookieDomain: "localhost:5555",
	CookieName: "JWTToken",
})

var routeLevels = map[string]map[string]string{
	"admin": {
		"create/course": "CreateCourse",
	},
	"any": {
		"course/:cid": "GetCourse",
		"course/:cid/assignments":                                   "CourseAssignments",
		"course/:cid/assignment/:aid/submission/download/:sid/:num": "DownloadSubmission",
		"course/:cid/assignment/:aid/details":                       "GetAssignment",
	},
	"assistant": map[string]string{
		"course/:cid/add/user":          "CourseAddUser",
		"course/:cid/assignment/create": "CreateAssignment",
		"course/:cid/assignment/fromfile": "CreateAssignmentFromFile",
		"course/:cid/assignment/:aid/csv": "GradesAsCSV",
		"course/:cid/assignment/:aid/update": "UpdateAssignment",
		"course/:cid/update": "UpdateCourse",
		"course/:cid/submission/:sid/update": "UpdateGrade",
	},
	"teacher": {
		"course/:cid/add/user":             "CourseAddUser",
		"course/:cid/add/users":            "CourseAddUsers",
		"course/:cid/assignment/create":    "CreateAssignment",
		"course/:cid/assignment/fromfile": "CreateAssignmentFromFile",
		"course/:cid/assignment/:aid/file": "AssignmentAsFile",
		"course/:cid/assignment/:aid/csv": "GradesAsCSV",
		"course/:cid/assignment/:aid/update": "UpdateAssignment",
		"course/:cid/update": "UpdateCourse",
		"course/:cid/submission/:sid/update": "UpdateGrade",
	},
	"student": {
		"course/:cid/:section/assignment/submit/:aid": "SubmitAssignment",
	},
}
