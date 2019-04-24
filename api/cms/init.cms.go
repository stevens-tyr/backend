package cms

import (
	"backend/models"
)

var am = models.NewMongoAssignmentInterface()
var cm = models.NewMongoCourseInterface()
var gfs = models.NewGridFSInterface()
var um = models.NewMongoUserInterface()
var sm = models.NewMongoSubmissionInterface()
