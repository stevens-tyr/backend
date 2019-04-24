package cms

import (
	"github.com/gin-gonic/gin"
)

func DeleteAssignment(c *gin.Context) {
	aid, _ := c.Get("aid")
	cid, _ := c.Get("cid")

	err := cm.RemoveAssignment(aid, cid)
	if err != nil {
		c.Set("error", err)
		return
	}

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

	err = am.Delete(aid)
	if err != nil {
		c.Set("error", err)
		return
	}

	err = sm.DeleteByAssignmentID(aid)
	if err != nil {
		c.Set("error", err)
		return
	}

	c.JSON(200, gin.H{
		"message": "Assignment Deleted.",
	})
}
