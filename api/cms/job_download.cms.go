package cms

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"

	"backend/errors"
)

func JobDownloadSubmission(c *gin.Context) {
	key := c.Param("secret")
	if key != os.Getenv("JOB_SECRET") {
		c.Set("error", errors.ErrorInvalidJobSecret)
		return
	}

	sid, _ := c.Get("sid")
	sub, err := sm.Get(sid)
	if err != nil {
		fmt.Println("here")
		c.Set("error", err)
		return
	}

	file, numBytes, err := gfs.Download(sub.FileID)
	if err != nil {
		c.Set("error", err)
		return
	}

	additonalHeaders := map[string]string{
		"Content-Disposition": fmt.Sprintf(`attachment; filename="%s-%s.tar.gz"`, c.Param("sid"), c.Param("num")),
	}

	c.DataFromReader(200, numBytes, "application/tar+gzip", file, additonalHeaders)
}

func JobDownloadSupportingFiles(c *gin.Context) {
	key := c.Param("secret")
	if key != os.Getenv("JOB_SECRET") {
		c.Set("error", errors.ErrorInvalidJobSecret)
		return
	}

	aid, _ := c.Get("aid")
	assign, err := am.Get(aid)
	if err != nil {
		c.Set("error", err)
		return
	}

	file, numBytes, err := gfs.Download(assign.SupportingFiles)
	if err != nil {
		c.Set("error", err)
		return
	}

	additonalHeaders := map[string]string{
		"Content-Disposition": fmt.Sprintf(`attachment; filename="%s-%s.tar.gz"`, c.Param("sid"), c.Param("num")),
	}

	c.DataFromReader(200, numBytes, "application/tar+gzip", file, additonalHeaders)
}
