package cms

import (
	"github.com/gin-gonic/gin"
)

func DeleteCourse(c *gin.Context) {
	cid, _ := c.Get("cid")

	err := um.RemoveCourseFromUsers(cid)
	if err != nil {
		c.Set("error", err)
		return
	}
	
	course, err := cm.GetByID(cid)
	if err != nil {
		c.Set("error", err)
		return
	}
	
	for _, aid := range course.Assignments {
		assign, err := am.Get(aid)
		if err != nil {
			c.Set("error", err)
			return
		}

		for _, sid := range assign.Submissions {
			err = gfs.Delete(sid)
			if err != nil {
				c.Set("error", err)
				return
			}
		}

		err = sm.DeleteByAssignmentID(aid)
		if err != nil {
			c.Set("error", err)
			return
		}
	}
	
	err = cm.Delete(cid)
	if err != nil {
		c.Set("error", err)
		return
	}

	c.JSON(200, gin.H{
		"message": "Course Deleted.",
	})
}
